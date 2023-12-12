package client

import (
	"context"

	pb "github.com/murtaza-u/keye"
)

// Del deletes the specified key and its corresponding value from the
// database, returning the deleted key(s). When used with WithRegex(),
// Del treats the "key" as a regex, deleting all matching keys.
func (c *C) Del(k string, optfns ...OptFunc) ([]string, error) {
	opts := defaultOpts()
	for _, fn := range optfns {
		fn(&opts)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	resp, err := c.api.Del(ctx, &pb.DelParams{
		Key: k,
		Opts: &pb.DelOpts{
			Regex: opts.regex,
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.GetKeys(), err
}
