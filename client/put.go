package client

import (
	"context"
	"errors"

	pb "github.com/murtaza-u/keye"

	"google.golang.org/grpc/status"
)

// Put adds a key-value pair to the database and returns the modified
// key(s). When used with WithRegex(), Put treats the "key" as a regex,
// updating all matching keys with the specified value.
func (c *C) Put(k string, v []byte, optfns ...OptFunc) ([]string, error) {
	opts := defaultOpts()
	for _, fn := range optfns {
		fn(&opts)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	resp, err := c.api.Put(ctx, &pb.PutParams{
		Key: k,
		Val: v,
		Opts: &pb.PutOpts{
			Regex: opts.regex,
		},
	})
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return nil, errors.New(stat.Message())
		}
		return nil, err
	}

	return resp.GetKeys(), err
}
