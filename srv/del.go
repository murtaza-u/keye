package srv

import (
	"context"
	"log/slog"
	"regexp"

	"github.com/murtaza-u/keye/internal/pb"
	"github.com/murtaza-u/keye/watch"

	"go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Del implements the gRPC API service Del method, deleting the
// specified key and its corresponding value from the database and
// returning the deleted key(s). When the regex option is enabled, Del
// treats the "key" as a regex, deleting all matching keys.
func (s *Srv) Del(ctx context.Context, in *pb.DelParams) (*pb.DelResponse, error) {
	k := in.GetKey()
	if k == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing key")
	}

	opts := in.GetOpts()
	if opts == nil {
		opts = &pb.DelOpts{
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
			v := b.Get([]byte(k))
			if v == nil {
				return status.Errorf(codes.NotFound, ErrKeyNotFound.Error())
			}
			err := b.Delete([]byte(k))
			if err != nil {
				return status.Errorf(codes.Internal,
					"failed to delete key from the bucket: %s", err.Error())
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
			if err := b.Delete(k); err != nil {
				return status.Errorf(codes.Internal,
					"failed to delete key from the bucket: %s", err.Error())
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

	keys := make([]string, len(kvs))
	for i, kv := range kvs {
		keys[i] = kv.K
	}

	slog.Debug("matched keys",
		slog.String("method", "Del"), slog.Int("count", len(keys)),
		slog.Any("keys", keys),
	)

	s.watcher.Push(watch.NewDelEvents(kvs...)...)

	return &pb.DelResponse{
		Keys: keys,
	}, err
}
