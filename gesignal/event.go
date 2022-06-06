package gesignal

// Void is a helper type that is recommended to be used instead of `struct{}`.
// You may want to use Void as the event type parameter when there is no
// useful data to be transmitted.
type Void struct{}

// Event is a slot-signal container.
// It holds all currect event listeners and invokes their callbacks
// when the associated event is triggered.
// An event is triggered when Emit() method is called.
//
// If you need 0 arguments callback, use Void type for the argument.
// If you need more than 1 argument in your callback, use tuple helper package.
// For example, a tuple.Value3[int, float, string] can be used to pass
// three arguments to your callback.
type Event[T any] struct {
	handlers []eventHandler[T]
}

// Connect adds an event listener that will be called for every Emit called for this event.
// When connection is disposed, an associated callback will be unregistered.
// If this connection should be persistent, pass a nil value as conn.
// For a non-nil conn, it's possible to disconnect from event by using Disconnect method.
func (e *Event[T]) Connect(conn connection, slot func(arg T)) {
	e.handlers = append(e.handlers, eventHandler[T]{
		c: conn,
		f: slot,
	})
}

// Disconnect removes an event listener identified by this connection.
// Note that you can't disconnect a listener that was connected with nil connection object.
func (e *Event[T]) Disconnect(conn connection) {
	for i, h := range e.handlers {
		if h.c == conn {
			e.handlers[i].c = theRemovedConnection
			break
		}
	}
}

// Emit triggers the associated event and calls all active callbacks with provided argument.
func (e *Event[T]) Emit(arg T) {
	// This method is slightly faster than the self-append alternative.
	length := 0
	for _, h := range e.handlers {
		if h.c != nil && h.c.IsDisposed() {
			continue
		}
		h.f(arg)
		e.handlers[length] = h
		length++
	}
	e.handlers = e.handlers[:length]
}

func (e *Event[T]) IsEmpty() bool {
	return len(e.handlers) == 0
}

type eventHandler[T any] struct {
	c connection
	f func(T)
}

type connection interface {
	IsDisposed() bool
}

type removedConnection struct{}

func (r *removedConnection) IsDisposed() bool { return true }

var theRemovedConnection = &removedConnection{}
