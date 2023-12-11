package watch

import "sync"

// Subscriber provides methods to receive subscribed events and cancel
// the subscription if needed.
type Subscriber struct {
	mu     sync.Mutex
	ev     chan Event
	closed bool
}

// newSubscriber creates a new subscriber.
func newSubscriber() *Subscriber {
	return &Subscriber{
		ev: make(chan Event),
	}
}

// NextEvent returns a channel of events, blocking when no events are
// available. Example usage:
//
//	for {
//	    ev := <-s.NextEvent()
//	}
func (s *Subscriber) NextEvent() <-chan Event { return s.ev }

// Close cancels the subsciption.
func (s *Subscriber) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.ev)
	s.closed = true
}

// push adds a new event to the subscriber queue. Clients can use the
// NextEvent() method to receive events.
func (s *Subscriber) push(ev Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return true
	}
	s.ev <- ev
	return false
}
