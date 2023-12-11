/*
Package watch implements a watcher using a pub-sub model. It allows
clients (subscribers) to subscribe to specific database keys and receive
PUT and DEL events on those keys.
*/
package watch

// W represents a watcher with an event queue and a look-up table. PUT
// and DEL operations in the database generate events pushed into the
// watcher's event queue. The look-up table maps subscribers to their
// subscribed keys, and the watcher forwards events to them.
type W struct {
	queue chan Event
	close chan struct{}
	table *lookUp
}

// New creates a watcher with an event queue of the specified size.
func New(queueSize uint) *W {
	return &W{
		queue: make(chan Event, queueSize),
		close: make(chan struct{}),
		table: newLookUp(),
	}
}

// Close shuts down the watcher.
func (w W) Close() { w.close <- struct{}{} }

// Push adds a new event to the watcher's event queue.
func (w W) Push(events ...Event) {
	for _, ev := range events {
		w.queue <- ev
	}
}

// Listen listens for new events in the event queue, publishing them to
// subscribers. It's a blocking function that exits only when the
// Close() method is called on the watcher.
func (w W) Listen() {
	for {
		select {
		case <-w.close:
			return
		case ev := <-w.queue:
			w.table.publish(ev)
		}
	}
}

// Watch subscribes to the specified keys and returns a subscriber
// object.
func (w W) Watch(keys ...string) *Subscriber {
	s := newSubscriber()
	w.table.subscribe(s, keys...)
	return s
}
