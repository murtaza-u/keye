package client

import (
	"context"
	"errors"

	pb "github.com/murtaza-u/keye"
	"github.com/murtaza-u/keye/watch"

	"google.golang.org/grpc/status"
)

// Watch watches specified key(s) in the database. If used with
// WithRegex(), it treats the "key" as a regex, watching all matched
// keys.
func (c *C) Watch(ctx context.Context, k string, optfns ...OptFunc) (*Watcher, error) {
	opts := defaultOpts()
	for _, fn := range optfns {
		fn(&opts)
	}

	stream, err := c.api.Watch(ctx, &pb.WatchParams{
		Key: k,
		Opts: &pb.WatchOpts{
			Regex: opts.regex,
		},
	})
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return nil, errors.New(stat.Message())
		}
		return nil, err
	}

	w := newWatcher()
	go w.listen(stream)

	return w, nil
}

// Watcher notifies clients of subscribed key events.
type Watcher struct {
	ev  chan watch.Event
	err chan error
}

// NextEvent returns the next event channel.
func (w Watcher) NextEvent() <-chan watch.Event { return w.ev }

// Error returns the error channel.
func (w Watcher) Error() <-chan error { return w.err }

// listen is a blocking method that streams events from the database
// server to the client via the watcher's event channel. It terminates
// in case of an error, returning the error through the watcher's error
// channel.
func (w Watcher) listen(stream pb.Api_WatchClient) {
	defer stream.CloseSend()

	for {
		resp, err := stream.Recv()
		if err != nil {
			if stat, ok := status.FromError(err); ok {
				w.err <- errors.New(stat.Message())
				return
			}
			w.err <- err
			return
		}

		var typ int
		switch resp.GetEvent() {
		case pb.Event_EVENT_KEEPALIVE:
			continue
		case pb.Event_EVENT_PUT:
			typ = watch.EventPut
		case pb.Event_EVENT_DEL:
			typ = watch.EventDel
		}

		kv := resp.GetKv()
		w.ev <- watch.Event{
			Type: typ,
			KV: watch.KV{
				K: kv.GetKey(),
				V: kv.GetVal(),
			},
		}
	}
}

func newWatcher() *Watcher {
	return &Watcher{
		ev:  make(chan watch.Event, 1),
		err: make(chan error),
	}
}
