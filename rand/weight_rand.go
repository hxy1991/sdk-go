package rand

import (
	"github.com/mroth/weightedrand"
	"math/rand"
)

type WeightedRandChoice interface {
	GetWeight() int
}

func WeightRand(choices []WeightedRandChoice) (WeightedRandChoice, error) {
	r := RandPool.pool.Get().(*rand.Rand)
	defer RandPool.pool.Put(r)

	var weightList []weightedrand.Choice
	for _, choice := range choices {
		weightList = append(weightList, weightedrand.Choice{
			Item:   choice,
			Weight: uint(choice.GetWeight()),
		})
	}

	chooser, err := weightedrand.NewChooser(weightList...)
	if err != nil {
		return nil, err
	}

	return chooser.PickSource(r).(WeightedRandChoice), nil
}
