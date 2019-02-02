package algorithm

import (
	"sort"
)

/*
	divison by weights
*/
func DivisionIntByWeights(sum int, weights []int) ([]int, []int) {
	if len(weights) == 0 {
		return nil, nil
	}

	allWeight := 0
	for _, v := range weights {
		allWeight += v
	}

	if allWeight <= 0 {
		return nil, nil
	}

	sort.Ints(weights)

	rets := make([]int, len(weights), len(weights))
	aliquant := sum
	//according to their weights divide approximately
	for k := range weights {
		n := sum * weights[k] / allWeight
		aliquant -= n
		rets[k] += n
	}

	//divide the rest part by sort index
	for aliquant > 0 {
		for k := range weights {
			if weights[k] > 0 {
				rets[k] += 1
				aliquant -= 1
			}
			if aliquant == 0 {
				break
			}
		}
	}

	return weights, rets
}
