package xiter

import (
	"bytes"
	"cmp"
	"slices"
	"testing"
)

func TestMap(t *testing.T) {
	s := OfSlice([]int{1, 2, 3})
	n := Collect(Map(s, func(v int) float64 { return float64(v * 2) }))
	if [3]float64(n) != [...]float64{2, 4, 6} {
		t.Fatal(n)
	}
}

func TestFilter(t *testing.T) {
	s := OfSlice([]int{1, 2, 3})
	n := Collect(Filter(s, func(v int) bool { return v%2 != 0 }))
	if [2]int(n) != [...]int{1, 3} {
		t.Fatal(n)
	}
}

func TestSkip(t *testing.T) {
	s := Collect(Skip(Limit(Generate(
		0, 1),
		3),
		2),
	)
	if !Equal(OfSlice(s), Of(2)) {
		t.Fatal(s)
	}
}

func TestLimit(t *testing.T) {
	s := Collect(Limit(Generate(
		0, 2),
		3),
	)
	if [3]int(s) != [...]int{0, 2, 4} {
		t.Fatal(s)
	}
}

func TestConcat(t *testing.T) {
	s := Collect(Concat(OfSlice([]int{1, 2, 3}), OfSlice([]int{3, 2, 5})))
	if [6]int(s) != [...]int{1, 2, 3, 3, 2, 5} {
		t.Fatal(s)
	}
}

func TestZip(t *testing.T) {
	s1 := OfSlice([]int{1, 2, 3, 4, 5})
	s2 := OfSlice([]int{2, 3, 4, 5, 6})
	seq := Zip(s1, s2)
	seq(func(v Zipped[int, int]) bool {
		if v.V2-v.V1 != 1 {
			t.Fatalf("unexpected values: %+v", v)
		}
		return true
	})
}

func TestZipShort1(t *testing.T) {
	s1 := OfSlice([]int{1, 2, 3, 4})
	s2 := OfSlice([]int{2, 3, 4, 5, 1})
	seq := Zip(s1, s2)
	seq(func(v Zipped[int, int]) bool {
		if v.V2-v.V1 != 1 {
			t.Fatalf("unexpected values: %+v", v)
		}
		return true
	})
	t.Logf("Greetings from TestZipShort1")
}

func TestZipShort2(t *testing.T) {
	s1 := OfSlice([]int{1, 2, 3, 4, -1, -2})
	s2 := OfSlice([]int{2, 3, 4, 5})
	seq := Zip(s1, s2)
	seq(func(v Zipped[int, int]) bool {
		if v.V1 == -2 {
			return false
		}
		if v.V2-v.V1 != 1 {
			t.Fatalf("unexpected values: %+v", v)
		}
		return true
	})
	t.Logf("Greetings from TestZipShort2")
}

func BenchmarkZip(b *testing.B) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{2, 3, 4, 5, 6}

	for i := 0; i < b.N; i++ {
		s1 := OfSlice(slice1)
		s2 := OfSlice(slice2)
		seq := Zip(s1, s2)
		seq(func(v Zipped[int, int]) bool {
			return true
		})
	}
}

func TestIsSorted(t *testing.T) {
	if IsSorted(OfSlice([]int{1, 2, 3, 2})) {
		t.Fatal("is not sorted")
	}
	if !IsSorted(OfSlice([]int{1, 2, 3, 4, 5})) {
		t.Fatal("is sorted")
	}
	if !IsSorted(OfSlice([]int{48, 48})) {
		t.Fatal("is sorted")
	}
}

func TestMerge(t *testing.T) {
	s1 := OfSlice([]int{2, 3, 5})
	s2 := OfSlice([]int{1, 2, 3, 4, 5})
	r := Collect(Merge(s1, s2))
	if [8]int(r) != [...]int{1, 2, 2, 3, 3, 4, 5, 5} {
		t.Fatal(r)
	}
}

func TestMergeFuncA(t *testing.T) {
	s1 := OfSlice([]int{0, 2, 4, 6, 8})
	s2 := OfSlice([]int{1, 3, 5, 7, 9})
	seq := MergeFunc(s1, s2, func(x, y int) int { return (x >> 1) - (y >> 1) })
	r := Collect(seq)
	if [10]int(r) != [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		t.Fatal(r)
	}
}

func TestMergeFuncB(t *testing.T) {
	s1 := OfSlice([]int{1, 3, 5, 7, 9})
	s2 := OfSlice([]int{0, 2, 4, 6, 8})
	seq := MergeFunc(s1, s2, func(x, y int) int { return (x >> 1) - (y >> 1) })
	r := Collect(seq)
	if [10]int(r) != [...]int{1, 0, 3, 2, 5, 4, 7, 6, 9, 8} {
		t.Fatal(r)
	}
}

func TestMergeFuncC(t *testing.T) {
	s1 := OfSlice([]int{0, 1, 2, 3, 5})
	s2 := OfSlice([]int{4, 6, 7, 8, 9})
	seq := MergeFunc(s1, s2, func(x, y int) int { return (x >> 1) - (y >> 1) })
	r := Collect(seq)
	// NB 5 precedes 4 because 5 >>1 == 4 >> 1 and 5 is in s1
	if [10]int(r) != [...]int{0, 1, 2, 3, 5, 4, 6, 7, 8, 9} {
		t.Fatal(r)
	}
}

func TestMergeFuncD(t *testing.T) {
	s1 := OfSlice([]int{4, 6, 7, 8, 9})
	s2 := OfSlice([]int{0, 1, 2, 3, 5})
	seq := MergeFunc(s1, s2, func(x, y int) int { return (x >> 1) - (y >> 1) })
	r := Collect(seq)
	if [10]int(r) != [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		t.Fatal(r)
	}
}

func TestMergeFuncEarlyA(t *testing.T) {
	s1 := OfSlice([]int{4, 6, 7, 8, -1, 9})
	s2 := OfSlice([]int{0, 1, 2, 3, 5})
	seq := MergeFunc(s1, s2, func(x, y int) int { return (x >> 1) - (y >> 1) })
	var r []int
	seq(func(v int) bool {
		if v < 0 {
			return false
		}
		r = append(r, v)
		return true
	})
	if [9]int(r) != [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8} {
		t.Fatal(r)
	}
}

func TestMergeFuncEarlyB(t *testing.T) {
	s1 := OfSlice([]int{4, 6, 7, 8, 9})
	s2 := OfSlice([]int{0, 1, 2, 3, -1, 5})
	seq := MergeFunc(s1, s2, func(x, y int) int { return (x >> 1) - (y >> 1) })
	var r []int
	seq(func(v int) bool {
		if v < 0 {
			return false
		}
		r = append(r, v)
		return true
	})
	if [4]int(r) != [...]int{0, 1, 2, 3} {
		t.Fatal(r)
	}
}

func splitmerge[T cmp.Ordered](s []T) Seq[T] {
	if len(s) <= 1 {
		return OfSlice(s)
	}

	left := splitmerge(s[:len(s)/2])
	right := splitmerge(s[len(s)/2:])
	return Merge(left, right)
}

func mergesort[T cmp.Ordered](s []T) {
	AppendTo(splitmerge(s), s[:0])
}

func splitmerge2P[T cmp.Ordered](s []T) Seq[T] {
	if len(s) <= 1 {
		return OfSlice(s)
	}

	left := splitmerge(s[:len(s)/2])
	right := splitmerge(s[len(s)/2:])
	return MergeFunc2Pull(left, right, cmp.Compare)
}

func mergesort2P[T cmp.Ordered](s []T) {
	AppendTo(splitmerge2P(s), s[:0])
}

func TestMergeSort(t *testing.T) {
	s := []int{3, 2, 5, 1, 6, 2}
	mergesort(s)
	if [6]int(s) != [...]int{1, 2, 2, 3, 5, 6} {
		t.Fatal(s)
	}
}

func BenchmarkMerge1Pull(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := []int{4, 61, 28, 19, 57, 72, 40, 90, 8, 87, 39, 25, 60, 79, 53, 51, 47, 94, 36, 34, 22, 50, 10, 2, 58, 73, 83, 31, 91, 64, 17, 86, 70, 3, 14, 5, 48, 24, 54, 69, 1, 92, 99, 33, 89, 7, 45, 11, 74, 84, 55, 97, 26, 71, 65, 88, 27, 49, 23, 18, 93, 78, 21, 59, 82, 66, 95, 13, 42, 32, 56, 68, 80, 96, 46, 77, 12, 38, 16, 63, 43, 85, 29, 35, 52, 15, 62, 30, 76, 98, 67, 44, 20, 37, 81, 75, 5, 6, 4, 8, 9, 7, 5, 3, 2, 1}
		mergesort(s)
	}
}

func BenchmarkMerge2Pull(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := []int{4, 61, 28, 19, 57, 72, 40, 90, 8, 87, 39, 25, 60, 79, 53, 51, 47, 94, 36, 34, 22, 50, 10, 2, 58, 73, 83, 31, 91, 64, 17, 86, 70, 3, 14, 5, 48, 24, 54, 69, 1, 92, 99, 33, 89, 7, 45, 11, 74, 84, 55, 97, 26, 71, 65, 88, 27, 49, 23, 18, 93, 78, 21, 59, 82, 66, 95, 13, 42, 32, 56, 68, 80, 96, 46, 77, 12, 38, 16, 63, 43, 85, 29, 35, 52, 15, 62, 30, 76, 98, 67, 44, 20, 37, 81, 75, 5, 6, 4, 8, 9, 7, 5, 3, 2, 1}
		mergesort2P(s)
	}
}

func FuzzMergeSort(f *testing.F) {
	f.Add([]byte("The quick brown fox jumped over the lazy dog."))
	f.Fuzz(func(t *testing.T, s []byte) {
		check := bytes.Clone(s)
		slices.Sort(check)

		mergesort(s)
		if !Equal(OfSlice(s), OfSlice(check)) {
			t.Fatal(s)
		}
	})
}

func TestChunks(t *testing.T) {
	s := Collect(Map(Chunks(OfSlice([]int{1, 2, 3, 4, 5}),
		2),
		slices.Clone),
	)
	if !slices.EqualFunc(s, [][]int{{1, 2}, {3, 4}, {5}}, slices.Equal) {
		t.Fatal(s)
	}
}

func TestSplit2(t *testing.T) {
	s1, s2 := CollectSplit(Split2(FromPair(OfSlice([]Pair[int32, int64]{{1, 2}, {3, 4}, {5, 6}}))))
	if !slices.Equal(s1, []int32{1, 3, 5}) {
		t.Fatal(s1)
	}
	if !slices.Equal(s2, []int64{2, 4, 6}) {
		t.Fatal(s2)
	}
}

func TestCache(t *testing.T) {
	var i int
	f := func(yield func(int) bool) {
		yield(i)
		i++
		return
	}
	seq := Cache(f)
	if s := Collect(seq); !slices.Equal(s, []int{0}) {
		t.Fatal(s)
	}
	if s := Collect(seq); !slices.Equal(s, []int{0}) {
		t.Fatal(s)
	}
}

func TestEnumerate(t *testing.T) {
	s := Collect(ToPair(Enumerate(Limit(Generate(0, 2), 3))))
	if !slices.Equal(s, []Pair[int, int]{{0, 0}, {1, 2}, {2, 4}}) {
		t.Fatal(s)
	}
}

func TestOr(t *testing.T) {
	s := Collect(Or(Of[int](), nil, Of(1, 2, 3), Of(4, 5, 6)))
	if !slices.Equal(s, []int{1, 2, 3}) {
		t.Fatal(s)
	}
}
