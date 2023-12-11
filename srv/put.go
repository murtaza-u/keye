package srv

import (
	"context"
	"regexp"

	pb "github.com/murtaza-u/keye"
	"github.com/murtaza-u/keye/watch"

	"go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Put implements the gRPC API service Put method, adding a key-value
// pair to the database and returning the modified key(s). When the
// regex option is enabled, Put treats the "key" as a regex, updating
// all matching keys with the specified value.
func (s *Srv) Put(ctx context.Context, in *pb.PutParams) (*pb.PutResponse, error) {
	k := in.GetKey()
	if k == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing key")
	}
	v := in.GetVal()
	if v == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing value")
	}

	opts := in.GetOpts()
	if opts == nil {
		opts = &pb.PutOpts{
			Regex: false,
		}
	}
	var kvs []watch.KV

	err := s.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return status.Errorf(codes.Internal,
				"failed to open database bucket: %s", err.Error())
		}

		if !opts.GetRegex() {
			err := b.Put([]byte(k), v)
			if err != nil {
				return status.Errorf(codes.Internal,
					"failed to put key into the bucket: %s", err.Error())
			}
			kvs = append(kvs, watch.KV{K: k, V: v})
			return nil
		}

		reg, err := regexp.Compile(k)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid regex")
		}

		err = b.ForEach(func(k, v []byte) error {
			if !reg.Match(k) {
				return nil
			}
			if err := b.Put(k, v); err != nil {
				return status.Errorf(codes.Internal,
					"failed to put key into the bucket: %s", err.Error())
			}
			kvs = append(kvs, watch.KV{K: string(k), V: v})
			return nil
		})
		if err != nil {
			return err
		}

		if len(kvs) == 0 {
			return status.Errorf(codes.NotFound, ErrNoKeysMatchRegex.Error())
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.watcher.Push(watch.NewPutEvents(kvs...)...)

	keys := make([]string, len(kvs))
	for i, kv := range kvs {
		keys[i] = kv.K
	}

	return &pb.PutResponse{
		Keys: keys,
	}, err
}
