package tuple

type Value2[T1 any, T2 any] struct {
	First  T1
	Second T2
}

func New2[T1 any, T2 any](first T1, second T2) Value2[T1, T2] {
	return Value2[T1, T2]{first, second}
}

func (v Value2[T1, T2]) Fields() (T1, T2) {
	return v.First, v.Second
}
