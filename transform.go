package xiter

import "cmp"

// Map returns a Seq that yields the values of seq transformed via f.
func Map[T1, T2 any](seq Seq[T1], f func(T1) T2) Seq[T2] {
	return func(yield func(T2) bool) bool {
		return seq(func(v T1) bool {
			return yield(f(v))
		})
	}
}

// Filter returns a Seq that yields only the values of seq for which
// f(value) returns true.
func Filter[T any](seq Seq[T], f func(T) bool) Seq[T] {
	return func(yield func(T) bool) bool {
		return seq(func(v T) bool {
			if !f(v) {
				return true
			}
			return yield(v)
		})
	}
}

// Limit returns a Seq that yields at most n values from seq.
func Limit[T any](seq Seq[T], n int) Seq[T] {
	return func(yield func(T) bool) bool {
		return seq(func(v T) bool {
			if !yield(v) {
				return false
			}
			n--
			return n > 0
		})
	}
}

// Concat creates a new Seq that yields the values of each of the
// provided Seqs in turn.
func Concat[T any](seqs ...Seq[T]) Seq[T] {
	return func(yield func(T) bool) bool {
		for _, seq := range seqs {
			seq(yield)
		}
		return false
	}
}

// Zipped holds values from an iteration of a Seq returned by [Zip].
type Zipped[T1, T2 any] struct {
	V1  T1
	OK1 bool

	V2  T2
	OK2 bool
}

// Zip returns a new Seq that yields the values of seq1 and seq2
// simultaneously.
func Zip[T1, T2 any](seq1 Seq[T1], seq2 Seq[T2]) Seq[Zipped[T1, T2]] {
	return func(yield func(Zipped[T1, T2]) bool) bool {
		p1, stop := Pull(seq1)
		defer stop()
		p2, stop := Pull(seq2)
		defer stop()

		for {
			var val Zipped[T1, T2]
			val.V1, val.OK1 = p1()
			val.V2, val.OK2 = p2()
			if (!val.OK1 && !val.OK2) || !yield(val) {
				return false
			}
		}
	}
}

// Merge returns a sequence that yields values from the ordered
// sequences seq1 and seq2 one at a time to produce a new ordered
// sequence made up of all of the elements of both seq1 and seq2.
func Merge[T cmp.Ordered](seq1, seq2 Seq[T]) Seq[T] {
	return MergeFunc(seq1, seq2, cmp.Compare)
}

// MergeFunc is like [Merge], but uses a custom comparison function
// for determining the order of values.
func MergeFunc[T any](seq1, seq2 Seq[T], compare func(T, T) int) Seq[T] {
	return func(yield func(T) bool) bool {
		p1, stop := Pull(seq1)
		defer stop()
		p2, stop := Pull(seq2)
		defer stop()

		v1, ok1 := p1()
		v2, ok2 := p2()
		for ok1 || ok2 {
			var c int
			if ok1 && ok2 {
				c = compare(v1, v2)
			}

			switch {
			case !ok2 || c<0:
				if !yield(v1) {
					return false
				}
				v1, ok1 = p1()
			case !ok1 || c>0:
				if !yield(v2) {
					return false
				}
				v2, ok2 = p2()
			default:
				if !yield(v1) || !yield(v2) {
					return false
				}
				v1, ok1 = p1()
				v2, ok2 = p2()
			}
		}

		return false
	}
}

// Windows returns a slice over successive overlapping portions of
// size n of the values yielded by seq. In other words,
//
//	Windows(Generate(0, 1), 3)
//
// will yield
//
//	[0, 1, 2]
//	[1, 2, 3]
//	[2, 3, 4]
//
// and so on. The slice yielded is reused from one iteration to the
// next, so it should not be held onto after each iteration has ended.
// [Map] and [slices.Clone] may come in handy for dealing with
// situations where this is necessary.
func Windows[T any](seq Seq[T], n int) Seq[[]T] {
	return func(yield func([]T) bool) bool {
		win := make([]T, 0, n)

		seq(func(v T) bool {
			if len(win) < n-1 {
				win = append(win, v)
				return true
			}
			if len(win) < n {
				win = append(win, v)
				return yield(win)
			}

			copy(win, win[1:])
			win[len(win)-1] = v
			return yield(win)
		})
		if len(win) < n {
			yield(win)
		}
		return false
	}
}

// Chunks works just like [Windows] except that the yielded slices of
// elements do not overlap. In other words,
//
//	Chunks(Generate(0, 1), 3)
//
// will yield
//
//	[0, 1, 2]
//	[3, 4, 5]
//	[6, 7, 8]
//
// Like with Windows, the slice is reused between iterations.
func Chunks[T any](seq Seq[T], n int) Seq[[]T] {
	return func(yield func([]T) bool) bool {
		win := make([]T, 0, n)

		seq(func(v T) bool {
			if len(win) == n {
				clear(win)
				win = win[:0]
			}

			if len(win) < n-1 {
				win = append(win, v)
				return true
			}
			if len(win) < n {
				win = append(win, v)
				return yield(win)
			}
			return true
		})
		if len(win) < n {
			yield(win)
		}
		return false
	}
}
