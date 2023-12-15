package srv

import (
	"context"
	"log/slog"
	"regexp"

	"github.com/murtaza-u/keye/internal/pb"

	"go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get implements the gRPC API service Get method, retrieving keys. By
// default, it returns the value for "key". When the regex option is
// enabled, Get treats the "key" as a regex. When the keysOnly option is
// enabled, Get only returns the key(s) without the value(s).
func (s *Srv) Get(ctx context.Context, in *pb.GetParams) (*pb.GetResponse, error) {
	k := in.GetKey()
	if k == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing key")
	}

	opts := in.GetOpts()
	if opts == nil {
		opts = &pb.GetOpts{
			Regex:    false,
			KeysOnly: false,
		}
	}

	var kvs []*pb.KV
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return status.Errorf(codes.NotFound, ErrKeyNotFound.Error())
		}

		if !opts.GetRegex() {
			v := b.Get([]byte(k))
			if v == nil {
				return status.Errorf(codes.NotFound, ErrKeyNotFound.Error())
			}
			kv := new(pb.KV)
			kv.Key = string(k)
			if !opts.GetKeysOnly() {
				kv.Val = v
			}
			kvs = append(kvs, kv)
			return nil
		}

		reg, err := regexp.Compile(k)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid regex")
		}

		b.ForEach(func(k, v []byte) error {
			if !reg.Match(k) {
				return nil
			}
			kv := new(pb.KV)
			kv.Key = string(k)
			if !opts.GetKeysOnly() {
				kv.Val = v
			}
			kvs = append(kvs, kv)
			return nil
		})

		if len(kvs) == 0 {
			return status.Errorf(codes.NotFound, ErrKeyNotFound.Error())
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(kvs))
	for i, kv := range kvs {
		keys[i] = kv.GetKey()
	}
	slog.Debug("matched keys",
		slog.String("method", "Get"), slog.Int("count", len(keys)),
		slog.Any("keys", keys),
	)

	return &pb.GetResponse{
		Kvs: kvs,
	}, nil
}
