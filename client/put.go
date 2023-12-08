package client

import (
	"context"

	pb "github.com/murtaza-u/keye"
)

// Put adds a key-value pair to the database and returns the modified
// key(s). When used with WithRegex(), Put treats the "key" as a regex,
// updating all matching keys with the specified value.
func (c *C) Put(k string, v []byte, optfuncs ...optFunc) ([]string, error) {
	opts := defaultOpts()
	for _, fn := range optfuncs {
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
		return nil, err
	}

	return resp.GetKeys(), err
}
