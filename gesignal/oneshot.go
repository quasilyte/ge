package gesignal

type oneshotConnector[T any] struct {
	d     disposable
	fired bool
}

func (c *oneshotConnector[T]) IsDisposed() bool {
	if c.fired {
		return true
	}
	return c.d != nil && c.d.IsDisposed()
}

func ConnectOneShot[T any](event *Event[T], d disposable, slot func(T)) {
	oneshot := &oneshotConnector[T]{d: d}
	event.Connect(oneshot, func(arg T) {
		oneshot.fired = true
		slot(arg)
	})
}
