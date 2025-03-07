package funcs

import (
	"math/rand"
	"time"
)

// WeightedRandomSelectionI  加权随机选择 (可用于业务流量分发等场景)
type WeightedRandomSelectionI interface {
	GetWeight() (weight int)
}

type WeightedRandomSelections[T WeightedRandomSelectionI] []T

func NewWeightedRandomSelections[T WeightedRandomSelectionI](ts ...T) WeightedRandomSelections[T] {
	return WeightedRandomSelections[T](ts)
}

func (ts WeightedRandomSelections[T]) Select() (out T) {
	if len(ts) == 0 {
		return out
	}

	var totalWeight int
	weights := make([]int, len(ts))

	for i, t := range ts {
		weights[i] = t.GetWeight()
		totalWeight += t.GetWeight()
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := rng.Intn(totalWeight)
	cumulativeWeight := 0

	for i, t := range ts {
		cumulativeWeight += weights[i]
		if rnd < cumulativeWeight {
			return t
		}
	}
	return out
}
