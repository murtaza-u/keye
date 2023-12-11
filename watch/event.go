package watch

const (
	// EventPut is fired when a key is modified or added to the
	// database.
	EventPut = iota + 1
	// EventDel is fired when a key is deleted from the database.
	EventDel
)

// Event is generated when a PUT or DEL operation occurs in the
// database.
type Event struct {
	// Type of the event. Possible values: EventPut, EventDel.
	Type int
	// KV is the affected key-value pair.
	KV KV
}

// KV represents a database key and its corresponding value.
type KV struct {
	K string
	V []byte
}

// NewPutEvents is a helper function that constructs new PUT events with
// the given key-value pairs.
func NewPutEvents(kvs ...KV) []Event {
	if kvs == nil {
		return nil
	}

	events := make([]Event, len(kvs))
	for i, kv := range kvs {
		events[i] = Event{
			Type: EventPut,
			KV:   kv,
		}
	}

	return events
}

// NewDelEvents is a helper function that constructs new DEL events with
// the given key-value pairs.
func NewDelEvents(kvs ...KV) []Event {
	if kvs == nil {
		return nil
	}

	events := make([]Event, len(kvs))
	for i, kv := range kvs {
		events[i] = Event{
			Type: EventDel,
			KV:   kv,
		}
	}

	return events
}
