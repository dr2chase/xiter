package xiter

import (
	"context"
	"slices"
	"testing"
)

func TestFind(t *testing.T) {
	s, _ := Find(Windows(Generate(
		0, 1),
		3),
		func(win []int) bool { return Sum(OfSlice(win)) >= 100 })
	if [3]int(s) != [...]int{33, 34, 35} {
		t.Fatal(s)
	}
}

func TestContains(t *testing.T) {
	c := Contains(Of(1, 2, 3), 2)
	if !c {
		t.Fatal(c)
	}
}

func TestSum(t *testing.T) {
	s := Sum(OfSlice([]string{"a", " ", "test"}))
	if s != "a test" {
		t.Fatal(s)
	}
}

func TestProduct(t *testing.T) {
	p := Product(Of(3, 2, -5))
	if p != -30 {
		t.Fatal(p)
	}
}

func TestEqual(t *testing.T) {
	a, b := Of(1, 2, 3), Of(1, 2)
	e := Of[int]()
	if Equal(a, b) {
		t.Fatal("a,b were equal")
	}
	if Equal(b, a) {
		t.Fatal("b,a were equal")
	}
	if Equal(a, e) {
		t.Fatal("a,e were equal")
	}
	if Equal(e, a) {
		t.Fatal("e,a were equal")
	}
	if !Equal(e, e) {
		t.Fatal("e,e were not equal")
	}
	if !Equal(a, a) {
		t.Fatal("a,a were not equal")
	}
}

func TestPartition(t *testing.T) {
	s1, s2 := Partition(Of(1, 2, 3, 4, 5), func(v int) bool { return v%2 == 0 })
	if !Equal(OfSlice(s1), Of(2, 4)) {
		t.Fatal(s1)
	}
	if !Equal(OfSlice(s2), Of(1, 3, 5)) {
		t.Fatal(s2)
	}
}

func TestExtent(t *testing.T) {
	s := Of(3, 2, 5, 1, 6, -2, 10)
	min := Min(s)
	max := Max(s)
	if min != -2 {
		t.Fatal(min)
	}
	if max != 10 {
		t.Fatal(max)
	}
}

func TestAny(t *testing.T) {
	r := Any(Of(2, 4, 6, 7), func(v int) bool { return v%2 != 0 })
	if !r {
		t.Fatal(r)
	}
}

func TestAll(t *testing.T) {
	r := All(Of(2, 4, 6, 7), func(v int) bool { return v%2 == 0 })
	if r {
		t.Fatal(r)
	}
}

func TestSendContext(t *testing.T) {
	c := make(chan int, 3)
	SendContext(Of(3, 2, 5), context.Background(), c)
	s := []int{<-c, <-c, <-c}
	select {
	case v := <-c:
		t.Fatal(v)
	default:
	}
	if !slices.Equal(s, []int{3, 2, 5}) {
		t.Fatal(s)
	}
}

func FuzzSendRecvContext(f *testing.F) {
	f.Add([]byte("The quick brown fox jumps over the lazy dog."))

	f.Fuzz(func(t *testing.T, data []byte) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c := make(chan byte, len(data))
		SendContext(OfSlice(data), ctx, c)
		close(c)
		s := Collect(RecvContext(ctx, c))
		if !slices.Equal(data, s) {
			t.Fatal(s)
		}
	})
}
