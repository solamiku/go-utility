package algorithm

import (
	"container/heap"
	"math"
	"math/rand"
)

type rsHeapItem struct {
	rand_w float64
	idx    int
}

type rsHeap []*rsHeapItem

func (rsh rsHeap) Len() int { return len(rsh) }

func (rsh rsHeap) Less(i, j int) bool { return rsh[i].rand_w < rsh[j].rand_w }

func (rsh rsHeap) Swap(i, j int) { rsh[i], rsh[j] = rsh[j], rsh[i] }

func (rsh *rsHeap) Push(x interface{}) {
	*rsh = append(*rsh, x.(*rsHeapItem))
}

func (rsh *rsHeap) Pop() interface{} {
	old := *rsh
	n := len(old)
	rshi := old[n-1]
	*rsh = old[0 : n-1]
	return rshi
}

type ReservoirSamplingInterface interface {
	Len() int
	GetWeight(idx int) float64
}

func ReservoirSampleWeight(sample ReservoirSamplingInterface, res_len int, r *rand.Rand) []int {
	result := make([]int, res_len)
	sample_len := sample.Len()

	h := &rsHeap{}
	for i := 0; i < sample_len; i++ {
		rshi := &rsHeapItem{
			idx:    i,
			rand_w: math.Pow(r.Float64(), (1.0 / sample.GetWeight(i))),
		}
		if h.Len() < res_len {
			heap.Push(h, rshi)
		} else {
			if (*h)[0].rand_w < rshi.rand_w {
				heap.Pop(h)
				heap.Push(h, rshi)
			}
		}
	}

	res_idx := 0
	for h.Len() > 0 {
		rshi := heap.Pop(h).(*rsHeapItem)
		result[res_idx] = rshi.idx
		res_idx++
	}
	return result

}

//id-weight basic
type BasicWeightObj struct {
	Id     int
	Weight float64
}
type basicReservoirSampling struct {
	objs []BasicWeightObj
}

func (brs basicReservoirSampling) Len() int {
	return len(brs.objs)
}

func (brs basicReservoirSampling) GetWeight(idx int) float64 {
	if idx < 0 || idx > len(brs.objs) {
		return 0
	}
	return brs.objs[idx].Weight
}

func (brs basicReservoirSampling) GetWeightResultId(r *rand.Rand) int {
	rs := ReservoirSampleWeight(brs, 1, r)
	if len(rs) > 0 {
		return brs.objs[rs[0]].Id
	}
	return -1
}

func (brs basicReservoirSampling) GetWeightResultIds(r *rand.Rand, num int) []int {
	rs := ReservoirSampleWeight(brs, num, r)
	ids := make([]int, 0, num)
	for _, v := range rs {
		ids = append(ids, brs.objs[v].Id)
	}
	return ids
}

func NewBasicReservoirSampling(objs []BasicWeightObj) basicReservoirSampling {
	return basicReservoirSampling{
		objs: objs,
	}
}
