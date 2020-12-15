package cf

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/recom/cf/base"
	"github.com/liwei1dao/lego/sys/recom/cf/core"
	"github.com/liwei1dao/lego/sys/recom/cf/model"
)

func Test_Main(t *testing.T) {
	// Load dataset
	data := core.NewDataSet([]string{"1001", "1002", "1003", "1001"}, []string{"2", "1", "3", "3"}, []float64{5, 5, 5, 5})
	// Split dataset
	// train, test := core.Split(data, 0.2)
	// Create model
	bpr := model.NewBPR(base.Params{
		base.NFactors:   10,
		base.Reg:        0.01,
		base.Lr:         0.05,
		base.NEpochs:    100,
		base.InitMean:   0,
		base.InitStdDev: 0.001,
	})
	// Fit model
	bpr.Fit(data, nil)
	// Evaluate model
	// scores := core.EvaluateRank(bpr, test, train, 10, core.Precision, core.Recall, core.NDCG)
	// fmt.Printf("Precision@10 = %.5f\n", scores[0])
	// fmt.Printf("Recall@10 = %.5f\n", scores[1])
	// fmt.Printf("NDCG@10 = %.5f\n", scores[2])
	// Generate recommendations for user(4):
	// Get all items in the full dataset
	items := core.Items(data)
	// Get user(4)'s ratings in the training dataset
	excludeItems := data.User("1003")
	// Get top 10 recommended items (excluding rated items) for user(4) using BPR
	recommendItems, _ := core.Top(items, "1003", 10, excludeItems, bpr)
	fmt.Printf("Recommend for user(1003) = %v\n", recommendItems)
}
