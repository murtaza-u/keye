package client

import (
	"context"
	"errors"

	"github.com/murtaza-u/keye/internal/pb"

	"google.golang.org/grpc/status"
)

// Get retrieves keys. By default, Get returns the value for "key", if
// any. When WithRegex() is passed, Get treats the "key" as a regex.
// When WithKeysOnly() is passed, Get only returns the key(s) without
// the value(s).
func (c *C) Get(k string, optfns ...OptFunc) ([]*pb.KV, error) {
	opts := defaultOpts()
	for _, fn := range optfns {
		fn(&opts)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	resp, err := c.api.Get(ctx, &pb.GetParams{
		Key: k,
		Opts: &pb.GetOpts{
			Regex:    opts.regex,
			KeysOnly: opts.keysOnly,
		},
	})
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return nil, errors.New(stat.Message())
		}
		return nil, err
	}

	return resp.GetKvs(), err
}
