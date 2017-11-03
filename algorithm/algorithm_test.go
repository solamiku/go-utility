package algorithm

import (
	"math/rand"
	"testing"
	"time"
)

func Test_algorithm(t *testing.T) {

}

func Benchmark_algorithm(b *testing.B) {
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
