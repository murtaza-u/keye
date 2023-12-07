package srv

import (
	"context"
	"regexp"

	pb "github.com/murtaza-u/keye"

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

	var keys []string

	err := s.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return status.Errorf(codes.Internal,
				"failed to open database bucket: %s", err.Error())
		}

		if !opts.GetRegex() {
			err := b.Delete([]byte(k))
			if err != nil {
				return status.Errorf(codes.Internal,
					"failed to delete key from the bucket: %s", err.Error())
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
			if err := b.Delete(k); err != nil {
				return status.Errorf(codes.Internal,
					"failed to delete key from the bucket: %s", err.Error())
			}
			keys = append(keys, string(k))
			return nil
		})

		return err
	})
	return &pb.DelResponse{
		Keys: keys,
	}, err
}
