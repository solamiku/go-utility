package algorithm

import (
	"container/heap"
	"math"
	"math/rand"
	"reflect"
)

type sampleItem struct {
	val    interface{}
	weight float64
}

type rsHeapItem struct {
	randWeight float64
	idx        int
}

type rsHeap []*rsHeapItem

func (rsh rsHeap) Len() int { return len(rsh) }

func (rsh rsHeap) Less(i, j int) bool { return rsh[i].randWeight < rsh[j].randWeight }

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

type reservoirSamplingInterface interface {
	Len() int
	GetWeight(idx int) float64
}

func reservoirSampleWeight(sample reservoirSamplingInterface,
	resultsLen int, randGen *rand.Rand) []int {
	minHeap := &rsHeap{}
	for idx, sampleLen := 0, sample.Len(); idx < sampleLen; idx++ {
		rshi := &rsHeapItem{
			idx:        idx,
			randWeight: math.Pow(randGen.Float64(), (1.0 / sample.GetWeight(idx))),
		}
		if minHeap.Len() < resultsLen {
			heap.Push(minHeap, rshi)
			continue
		}
		if (*minHeap)[0].randWeight < rshi.randWeight {
			heap.Pop(minHeap)
			heap.Push(minHeap, rshi)
		}
	}

	indice := make([]int, minHeap.Len())
	idx := 0
	for minHeap.Len() > 0 {
		rshi := heap.Pop(minHeap).(*rsHeapItem)
		indice[idx] = rshi.idx
		idx++
	}
	return indice
}

func NewReservoir() *Reservoir {
	return &Reservoir{items: make([]*sampleItem, 0)}
}

type Reservoir struct {
	items      []*sampleItem
	resultType reflect.Type
}

func (reservoir *Reservoir) Fill(val interface{}, weight float64) error {
	reservoir.items = append(reservoir.items,
		&sampleItem{val: val, weight: weight})
	if reservoir.resultType == nil {
		reservoir.resultType = reflect.SliceOf(reflect.TypeOf(val))
	}
	return nil
}

func (reservoir *Reservoir) Len() int                     { return len(reservoir.items) }
func (reservoir *Reservoir) GetWeight(idx int) float64    { return reservoir.items[idx].weight }
func (reservoir *Reservoir) GetValue(idx int) interface{} { return reservoir.items[idx].val }

func (reservoir *Reservoir) Sample(resultsLen int, randGen *rand.Rand) interface{} {
	if reservoir.Len() == 0 {
		return nil
	}
	if resultsLen > reservoir.Len() {
		resultsLen = reservoir.Len()
	}
	resultIndice := reservoirSampleWeight(reservoir, resultsLen, randGen)
	resultSlice := reflect.MakeSlice(reservoir.resultType, resultsLen, resultsLen)
	for i, idx := range resultIndice {
		resultSlice.Index(i).Set(reflect.ValueOf(reservoir.GetValue(idx)))
	}
	return resultSlice.Interface()
}

func (reservoir *Reservoir) SampleOne(randGen *rand.Rand) interface{} {
	if reservoir.Len() == 0 {
		return nil
	}
	resultIndice := reservoirSampleWeight(reservoir, 1, randGen)
	idx := resultIndice[0]
	return reservoir.GetValue(idx)
}
