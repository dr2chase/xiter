package xiter

// Slice returns a Seq over the elements of s.
func Slice[T any, S ~[]T](s S) Seq[T] {
	return func(yield func(T) bool) bool {
		for _, v := range s {
			if !yield(v) {
				return false
			}
		}
		return false
	}
}

// MapEntry represents a key-value pair.
type MapEntry[K comparable, V any] struct {
	Key K
	Val V
}

func (e MapEntry[K, V]) key() K { return e.Key }
func (e MapEntry[K, V]) val() V { return e.Val }

// MapEntries returns a Seq over the key-value pairs of m.
func MapEntries[K comparable, V any, M ~map[K]V](m M) Seq[MapEntry[K, V]] {
	return func(yield func(MapEntry[K, V]) bool) bool {
		for k, v := range m {
			if !yield(MapEntry[K, V]{k, v}) {
				return false
			}
		}
		return false
	}
}

// MapKeys returns a Seq over the keys of m.
func MapKeys[K comparable, V any, M ~map[K]V](m M) Seq[K] {
	return Map(MapEntries(m), MapEntry[K, V].key)
}

// MapValues returns a Seq over the values of m.
func MapValues[K comparable, V any, M ~map[K]V](m M) Seq[V] {
	return Map(MapEntries(m), MapEntry[K, V].val)
}
