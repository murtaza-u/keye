package srv

import (
	"log"
	"regexp"
	"time"

	pb "github.com/murtaza-u/keye"

	"go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Watch implements the gRPC API service Watch method, enabling clients
// to watch over keys in the database. Operations on watched keys
// trigger events to the client. When the regex option is enabled, Watch
// treats "key" as a regex.
//
// Note: Watch throws an error if the regex doesn't match any existing
// key in the database.
func (s *Srv) Watch(in *pb.WatchParams, stream pb.Api_WatchServer) error {
	k := in.GetKey()
	if k == "" {
		return status.Errorf(codes.InvalidArgument, "missing key")
	}

	opts := in.GetOpts()
	if opts == nil {
		opts = &pb.WatchOpts{
			Regex: false,
		}
	}

	// retrieve matching keys
	var keys []string
	err := s.db.View(func(tx *bbolt.Tx) error {
		if !opts.GetRegex() {
			keys = append(keys, k)
			return nil
		}

		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return status.Errorf(codes.NotFound, ErrKeyNotFound.Error())
		}

		reg, err := regexp.Compile(k)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid regex")
		}

		b.ForEach(func(k, v []byte) error {
			if !reg.Match(k) {
				return nil
			}
			keys = append(keys, string(k))
			return nil
		})

		if len(keys) == 0 {
			return status.Errorf(codes.NotFound, ErrNoKeysMatchRegex.Error())
		}

		return nil
	})
	if err != nil {
		return err
	}

	sub := s.watcher.Watch(keys...)
	defer sub.Close()

	t := time.NewTicker(s.wpi)
	defer t.Stop()

	for {
		select {
		// keepalive messages are sent to the client at a configured
		// interval to maintain a live connection.
		case <-t.C:
			err = stream.Send(&pb.WatchResponse{
				Event: pb.Event_EVENT_KEEPALIVE,
				Kv:    nil,
			})
		case <-stream.Context().Done():
			return nil
		case ev := <-sub.NextEvent():
			err = stream.Send(&pb.WatchResponse{
				Event: pb.Event(ev.Type),
				Kv: &pb.KV{
					Key: ev.KV.K,
					Val: ev.KV.V,
				},
			})

			// Reset the ticker if the occurred event is not a keepalive
			// event.
			t.Reset(s.wpi)
		}

		if err != nil {
			if stat, ok := status.FromError(err); ok {
				log.Println(stat.Err())
				return nil
			}
			log.Println(err)
			return nil
		}
	}
}
