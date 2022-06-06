package gesignal

type oneshotConnector[T any] struct {
	conn  connection
	fired bool
}

func (c *oneshotConnector[T]) IsDisposed() bool {
	if c.fired {
		return true
	}
	return c.conn != nil && c.conn.IsDisposed()
}

func ConnectOneShot[T any](event *Event[T], conn connection, slot func(T)) {
	oneshot := &oneshotConnector[T]{conn: conn}
	event.Connect(oneshot, func(arg T) {
		oneshot.fired = true
		slot(arg)
	})
}
