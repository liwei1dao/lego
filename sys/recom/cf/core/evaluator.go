package core

import (
	"github.com/liwei1dao/lego/sys/recom/cf/base"
	"github.com/liwei1dao/lego/sys/recom/cf/floats"
)

type RankMetric func(targetSet *base.MarginalSubSet, rankList []string) float64

// EvaluateRank evaluates a model in top-n tasks.
func EvaluateRank(estimator ModelInterface, testSet DataSetInterface, excludeSet DataSetInterface, n int, metrics ...RankMetric) []float64 {
	sum := make([]float64, len(metrics))
	count := 0.0
	items := Items(testSet, excludeSet)
	// For all users
	for userIndex := 0; userIndex < testSet.UserCount(); userIndex++ {
		userId := testSet.UserIndexer().ToID(userIndex)
		// Find top-n items in test set
		targetSet := testSet.UserByIndex(userIndex)
		if targetSet.Len() > 0 {
			// Find top-n items in predictions
			rankList, _ := Top(items, userId, n, excludeSet.User(userId), estimator)
			count++
			for i, metric := range metrics {
				sum[i] += metric(targetSet, rankList)
			}
		}
	}
	floats.MulConst(sum, 1/count)
	return sum
}

func Precision(targetSet *base.MarginalSubSet, rankList []string) float64 {
	hit := 0.0
	for _, itemId := range rankList {
		if targetSet.Contain(itemId) {
			hit++
		}
	}
	return float64(hit) / float64(len(rankList))
}
