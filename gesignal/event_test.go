package gesignal

import (
	"testing"

	"github.com/quasilyte/ge/tuple"
)

type eventConnection struct {
	disposed bool
}

func (c *eventConnection) IsDisposed() bool { return c.disposed }

func TestEvent(t *testing.T) {
	var e Event[int]
	counter := 0
	e.Connect(nil, func(arg int) {
		counter += arg
	})
	e.Connect(nil, func(arg int) {
		counter++
	})
	e.Emit(10)
	if counter != 11 {
		t.Fatal("unexpected counter value")
	}
	e.Emit(1)
	if counter != 13 {
		t.Fatal("unexpected counter value")
	}
	e.Emit(0)
	if counter != 14 {
		t.Fatal("unexpected counter value")
	}
}

func TestEventDisposed(t *testing.T) {
	var e Event[int]
	counter := 0
	conn := eventConnection{disposed: true}
	e.Connect(&conn, func(arg int) {
		counter += arg
	})
	if len(e.handlers) != 1 {
		t.Fatal("unexpected number of handlers")
	}
	e.Emit(1)
	e.Emit(1)
	if counter != 0 {
		t.Fatal("unexpected counter value")
	}
	if len(e.handlers) != 0 {
		t.Fatal("unexpected number of handlers")
	}
}

func TestEventDisposed2(t *testing.T) {
	var e Event[Void]
	var connections [3]eventConnection
	var counters [3]int
	for i := range connections {
		id := i
		e.Connect(&connections[id], func(Void) {
			counters[id]++
		})
	}
	if len(e.handlers) != 3 {
		t.Fatal("unexpected number of handlers")
	}
	e.Emit(Void{})
	for _, v := range counters {
		if v != 1 {
			t.Fatal("unexpected counter value")
		}
	}
	connections[2].disposed = true
	e.Emit(Void{})
	if len(e.handlers) != 2 {
		t.Fatal("unexpected number of handlers")
	}
	if counters[2] != 1 {
		t.Fatal("unexpected counter value")
	}
	if counters[0] != 2 || counters[1] != 2 {
		t.Fatal("unexpected counter value")
	}
	connections[0].disposed = true
	e.Emit(Void{})
	if len(e.handlers) != 1 {
		t.Fatal("unexpected number of handlers")
	}
	if counters[2] != 1 || counters[0] != 2 {
		t.Fatal("unexpected counter value")
	}
	if counters[1] != 3 {
		t.Fatal("unexpected counter value")
	}
	e.Emit(Void{})
	if counters[1] != 4 {
		t.Fatal("unexpected counter value")
	}
	connections[1].disposed = true
	e.Emit(Void{})
	if len(e.handlers) != 0 {
		t.Fatal("unexpected number of handlers")
	}
	if counters[1] != 4 {
		t.Fatal("unexpected counter value")
	}
}

func TestEventDisposed3(t *testing.T) {
	var e Event[Void]
	var connections [3]eventConnection
	var counters [3]int
	for i := range connections {
		id := i
		e.Connect(&connections[id], func(Void) {
			counters[id]++
		})
	}
	connections[0].disposed = true
	connections[2].disposed = true
	e.Emit(Void{})
	e.Emit(Void{})
	e.Emit(Void{})
	if len(e.handlers) != 1 {
		t.Fatal("unexpected number of handlers")
	}
	if counters[0] != 0 || counters[2] != 0 {
		t.Fatal("unexpected counter value")
	}
	if counters[1] != 3 {
		t.Fatal("unexpected counter value")
	}
	connections[1].disposed = true
	for i := 0; i < 2; i++ {
		e.Emit(Void{})
		if len(e.handlers) != 0 {
			t.Fatal("unexpected number of handlers")
		}
		if counters[0] != 0 || counters[1] != 3 || counters[2] != 0 {
			t.Fatal("unexpected counter value")
		}
	}
}

func TestDisconnect(t *testing.T) {
	var e Event[int]
	var conn1 eventConnection
	var conn2 eventConnection
	counter := 0
	e.Connect(&conn1, func(arg int) {
		counter += arg
	})
	e.Connect(&conn2, func(arg int) {
		counter -= arg
	})
	e.Emit(10)
	e.Emit(5)
	if counter != 0 {
		t.Fatal("unexpected counter value")
	}
	if len(e.handlers) != 2 {
		t.Fatal("unexpected number of handlers")
	}
	e.Disconnect(&conn1)
	// The disconnection is delayed until the emit, so the number
	// of handlers will remain the same.
	if len(e.handlers) != 2 {
		t.Fatal("unexpected number of handlers")
	}
	e.Emit(5)
	if counter != -5 {
		t.Fatal("unexpected counter value")
	}
	if len(e.handlers) != 1 {
		t.Fatal("unexpected number of handlers")
	}
	e.Emit(1)
	if counter != -6 {
		t.Fatal("unexpected counter value")
	}
	conn2.disposed = true
	e.Emit(1)
	if counter != -6 {
		t.Fatal("unexpected counter value")
	}
	if len(e.handlers) != 0 {
		t.Fatal("unexpected number of handlers")
	}
}

func BenchmarkEmitVoid(b *testing.B) {
	var e0 Event[Void]
	var e1 Event[Void]
	e1.Connect(nil, func(Void) {})
	var e2 Event[Void]
	e2.Connect(nil, func(Void) {})
	e2.Connect(nil, func(Void) {})

	b.Run("Count0", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e0.Emit(Void{})
		}
	})
	b.Run("Count1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e1.Emit(Void{})
		}
	})
	b.Run("Count2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e2.Emit(Void{})
		}
	})
}

func BenchmarkEmitTuple(b *testing.B) {
	var e0 Event[tuple.Value2[int, int]]
	var e1 Event[tuple.Value2[int, int]]
	e1.Connect(nil, func(tuple.Value2[int, int]) {})
	var e2 Event[tuple.Value2[int, int]]
	e2.Connect(nil, func(tuple.Value2[int, int]) {})
	e2.Connect(nil, func(tuple.Value2[int, int]) {})

	b.Run("Count0", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e0.Emit(tuple.New2(1, 2))
		}
	})
	b.Run("Count1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e1.Emit(tuple.New2(1, 2))
		}
	})
	b.Run("Count2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e2.Emit(tuple.New2(1, 2))
		}
	})
}

func BenchmarkConnectEmitDisconnectEmit(b *testing.B) {
	var e Event[Void]
	conn1 := &eventConnection{}
	conn2 := &eventConnection{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn1.disposed = false
		conn2.disposed = false
		e.Connect(conn1, func(Void) {})
		e.Connect(conn2, func(Void) {})
		e.Emit(Void{})
		e.Disconnect(conn1)
		e.Disconnect(conn2)
		e.Emit(Void{})
	}
}
