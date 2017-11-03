package algorithm

import "math/rand"

//simple shuffle algorithm
type ShuffleInterface interface {
	Len() int
	Swap(i, j int)
}

type IntSlice []int

func (p IntSlice) Len() int      { return len(p) }
func (p IntSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func FisherYatesShuffle(dat ShuffleInterface, r *rand.Rand) {
	dat_len := dat.Len()
	for i := dat_len - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		dat.Swap(i, j)
	}
}

func FisherYatesShuffleInt(src []int, r *rand.Rand) {
	FisherYatesShuffle(IntSlice(src), r)
}
