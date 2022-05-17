package gesignal

type Void struct{}

type Event[T any] struct {
	handlers []eventHandler[T]
}

type eventHandler[T any] struct {
	d disposable
	f func(T)
}

func (e *Event[T]) Connect(d disposable, slot func(arg T)) {
	e.handlers = append(e.handlers, eventHandler[T]{
		d: d,
		f: slot,
	})
}

func (e *Event[T]) Disconnect(d disposable) {
	for i, h := range e.handlers {
		if h.d == d {
			e.handlers[i].d = theRemovedListener
			break
		}
	}
}

func (e *Event[T]) Emit(arg T) {
	// This method is slightly faster than the self-append alternative.
	length := 0
	for _, h := range e.handlers {
		if h.d != nil && h.d.IsDisposed() {
			continue
		}
		h.f(arg)
		e.handlers[length] = h
		length++
	}
	e.handlers = e.handlers[:length]
}

type disposable interface {
	IsDisposed() bool
}

type removedListener struct{}

func (r *removedListener) IsDisposed() bool { return true }

var theRemovedListener = &removedListener{}
