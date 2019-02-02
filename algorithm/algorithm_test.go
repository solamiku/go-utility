package algorithm

import (
	"math/rand"
	"testing"
	"time"
)

func Test_algorithm(t *testing.T) {
	t.Log(DivisionIntByWeights(101, []int{2, 1, 1}))
	t.Log(DivisionIntByWeights(102, []int{1, 10, 2}))
}

func HIDE_Benchmark_algorithm(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ba := NewBasicReservoirSampling([]BasicWeightObj{
		BasicWeightObj{Id: 1, Weight: 1.0},
		BasicWeightObj{Id: 10, Weight: 1.0},
		BasicWeightObj{Id: 20, Weight: 1.0},
		BasicWeightObj{Id: 30, Weight: 1.0},
	})
	for i := 0; i < b.N; i++ {
		ba.GetWeightResultId(r)
	}
}

func Benchmark_Division1(b *testing.B) {
	b.Log(DivisionIntByWeights(100000000, []int{10000000, 20000000, 30000000, 4000000, 50000000}))
	for i := 0; i < b.N; i++ {
		DivisionIntByWeights(100000000, []int{10000000, 20000000, 30000000, 4000000, 50000000})
	}
}
