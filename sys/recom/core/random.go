package core

import "math/rand"

func NewRandomGenerator(seed int64) RandomGenerator {
	return RandomGenerator{rand.New(rand.NewSource(int64(seed)))}
}

type RandomGenerator struct {
	*rand.Rand
}

func (rng RandomGenerator) NewNormalVector(size int, mean, stdDev float64) []float64 {
	ret := make([]float64, size)
	for i := 0; i < len(ret); i++ {
		ret[i] = rng.NormFloat64()*stdDev + mean
	}
	return ret
}
func (rng RandomGenerator) NewNormalMatrix(row, col int, mean, stdDev float64) [][]float64 {
	ret := make([][]float64, row)
	for i := range ret {
		ret[i] = rng.NewNormalVector(col, mean, stdDev)
	}
	return ret
}
