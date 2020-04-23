package algorithm

import (
	"math/rand"
	"testing"
	"time"
)

func Test_algorithm(t *testing.T) {
	t.Log(DivisionIntByWeights(101, []int{2, 1, 1}))
	t.Log(DivisionIntByWeights(102, []int{1, 10, 2}))

	ba, r := newReservoir()
	t.Log(ba.SampleOne(r))
}

func newReservoir() (*Reservoir, *rand.Rand) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ba := NewReservoir()
	ba.Fill(1, 1.0)
	ba.Fill(2, 2.0)
	ba.Fill(3, 3.0)
	ba.Fill(4, 4.0)
	ba.Fill(5, 5.0)
	ba.Fill(6, 6.0)
	ba.Fill(7, 7.0)
	return ba, r
}

func Benchmark_algorithm(b *testing.B) {
	ba, r := newReservoir()
	for i := 0; i < b.N; i++ {
		ba.SampleOne(r)
	}
}

func Benchmark_Division1(b *testing.B) {
	b.Log(DivisionIntByWeights(100000000, []int{10000000, 20000000, 30000000, 4000000, 50000000}))
	for i := 0; i < b.N; i++ {
		DivisionIntByWeights(100000000, []int{10000000, 20000000, 30000000, 4000000, 50000000})
	}
}
