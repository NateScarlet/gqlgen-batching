package iterator

func Keys[K, V any](
	seq func(func(K, V) bool),
	yieldValue func(V) bool,
) func(func(K) bool) {
	return func(yield func(K) bool) {
		seq(func(k K, v V) bool {
			if yieldValue != nil && !yieldValue(v) {
				return false
			}
			return yield(k)
		})
	}
}

func Values[T1, T2 any](
	seq func(func(T1, T2) bool),
	yieldKey func(T2) bool,
) func(func(T2) bool) {
	return func(yield func(T2) bool) {
		seq(func(k T1, v T2) bool {
			if yieldKey != nil && !yieldKey(v) {
				return false
			}
			return yield(v)
		})
	}
}
