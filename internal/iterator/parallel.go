package iterator

import (
	"context"
	"iter"
	"slices"
	"sync"
)

// Parallel2 returns order preserved result.
//
// limit is the max goroutine for parallel execution, non-positive value means no limit.
func Parallel2[T any, K any, V any](ctx context.Context, limit int, seq iter.Seq[T], f func(T) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		type output struct {
			index int
			key   K
			value V
		}
		var out = make(chan output)
		go func() {
			defer close(out)
			var wg = new(sync.WaitGroup)
			func() {
				type input struct {
					index int
					value T
				}
				var in = make(chan input)
				defer close(in)
				var nextIndex int
				var workerCount int
				for i := range seq {
					select {
					case <-ctx.Done():
						return
					case in <- input{nextIndex, i}:
						nextIndex++
					default:
						if limit <= 0 || workerCount < limit {
							workerCount++
							wg.Add(1)
							go func() {
								defer wg.Done()
								for i := range in {
									k, v := f(i.value)
									select {
									case <-ctx.Done():
										return
									case out <- output{i.index, k, v}:
									}
								}
							}()
						}
						select {
						case <-ctx.Done():
							return
						case in <- input{nextIndex, i}:
							nextIndex++
						}
					}
				}
			}()
			wg.Wait()
		}()

		var nextIndex int
		var b []output
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-out:
				if !ok {
					if len(b) > 0 {
						panic("parallel execution result incomplete")
					}
					return
				}
				if v.index == nextIndex {
					// returned in order
					nextIndex++
					if !yield(v.key, v.value) {
						return
					}
					for len(b) > 0 && b[len(b)-1].index == nextIndex {
						// continue with buffer
						nextIndex++
						v, b = b[len(b)-1], b[:len(b)-1]
						if !yield(v.key, v.value) {
							return
						}
					}
				} else {
					// not in order, save to buffer
					var i, _ = slices.BinarySearchFunc(b, v, func(el, target output) int {
						return target.index - el.index
					})
					b = slices.Insert(b, i, v)
				}
			}
		}
	}
}

func Parallel[T any, K any](ctx context.Context, limit int, seq func(func(T) bool), f func(T) K) func(func(K) bool) {
	return Keys(Parallel2(ctx, limit, seq, func(i T) (K, struct{}) {
		return f(i), struct{}{}
	}), nil)
}
