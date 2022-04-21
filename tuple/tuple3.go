package tuple

type Value3[T1 any, T2 any, T3 any] struct {
	First  T1
	Second T2
	Third  T3
}

func New3[T1 any, T2 any, T3 any](first T1, second T2, third T3) Value3[T1, T2, T3] {
	return Value3[T1, T2, T3]{first, second, third}
}

func (v Value3[T1, T2, T3]) Fields() (T1, T2, T3) {
	return v.First, v.Second, v.Third
}
