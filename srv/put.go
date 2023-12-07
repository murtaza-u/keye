package srv

import (
	"context"
	"regexp"

	pb "github.com/murtaza-u/keye"

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
	var keys []string

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
			keys = append(keys, k)
			return nil
		}

		reg, err := regexp.Compile(k)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid regex")
		}

		err = b.ForEach(func(k, _ []byte) error {
			if !reg.Match(k) {
				return nil
			}
			if err := b.Put(k, v); err != nil {
				return status.Errorf(codes.Internal,
					"failed to put key into the bucket: %s", err.Error())
			}
			keys = append(keys, string(k))
			return nil
		})

		return err
	})
	return &pb.PutResponse{
		Keys: keys,
	}, err
}
