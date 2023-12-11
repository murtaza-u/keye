package watch

import (
	"container/list"
	"sync"
)

// lookUp is a table that maps subscribers to their subscribed keys.
type lookUp struct {
	sync.RWMutex
	table map[string]*list.List
}

// newLookUp creates a new empty lookup table.
func newLookUp() *lookUp {
	return &lookUp{
		table: make(map[string]*list.List),
	}
}

// subscribe subscribes a subscriber to the specified keys.
func (l *lookUp) subscribe(s *Subscriber, keys ...string) {
	l.Lock()
	defer l.Unlock()

	for _, k := range keys {
		if l.table[k] == nil {
			l.table[k] = list.New()
		}
		l.table[k].PushBack(s)
	}
}

// publish sends an event to all its subscribers and removes any
// canceled subscriptions.
func (l *lookUp) publish(ev Event) {
	k := ev.KV.K

	l.RLock()
	lst := l.table[k]
	l.RUnlock()

	if lst == nil {
		return
	}

	for e := lst.Front(); e != nil; e = e.Next() {
		s, ok := e.Value.(*Subscriber)
		if !ok {
			lst.Remove(e)
			continue
		}
		closed := s.push(Event{
			Type: ev.Type,
			KV:   ev.KV,
		})
		if closed {
			lst.Remove(e)
		}
	}

	l.Lock()
	l.table[k] = lst
	if lst.Len() == 0 {
		delete(l.table, k)
	}
	l.Unlock()
}
