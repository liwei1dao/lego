package core

import "math/rand"

func Split(data DataSetInterface, testRatio float64) (train, test DataSetInterface) {
	testSize := int(float64(data.Count()) * testRatio)
	perm := rand.Perm(data.Count())
	// Test Data
	testIndex := perm[:testSize]
	test = data.SubSet(testIndex)
	// Train Data
	trainIndex := perm[testSize:]
	train = data.SubSet(trainIndex)
	return
}
