package model

import (
	"math"

	"github.com/liwei1dao/lego/sys/recom/core"
	"github.com/liwei1dao/lego/sys/recom/floats"
)

func NewBPR(params Params) *BPR {
	bpr := new(BPR)
	bpr.SetParams(params)
	return bpr
}

//贝叶斯个性化推荐 模型
type BPR struct {
	ModelBase
	UserFactor  [][]float64 // p_u
	ItemFactor  [][]float64 // q_i
	nFactors    int
	nEpochs     int
	lr          float64
	reg         float64
	initMean    float64
	initStdDev  float64
	UserRatings []*core.MarginalSubSet
	ItemPop     *ItemPop
}

func (bpr *BPR) SetParams(params Params) {
	bpr.ModelBase.SetParams(params)
	bpr.nFactors = bpr.Params.GetInt(NFactors, 10)
	bpr.nEpochs = bpr.Params.GetInt(NEpochs, 100)
	bpr.lr = bpr.Params.GetFloat64(Lr, 0.05)
	bpr.reg = bpr.Params.GetFloat64(Reg, 0.01)
	bpr.initMean = bpr.Params.GetFloat64(InitMean, 0)
	bpr.initStdDev = bpr.Params.GetFloat64(InitStdDev, 0.001)
}

func (bpr *BPR) Predict(userId, itemId uint32) float64 {
	// Convert sparse IDs to dense IDs
	userIndex := bpr.UserIndexer.ToIndex(userId)
	itemIndex := bpr.ItemIndexer.ToIndex(itemId)
	if userIndex == core.NotId || bpr.UserRatings[userIndex].Len() == 0 {
		// If users not exist in dataset, use ItemPop model.
		return bpr.ItemPop.Predict(userId, itemId)
	}
	return bpr.predict(userIndex, itemIndex)
}

func (bpr *BPR) Fit(trainSet core.DataSetInterface) {
	bpr.Init(trainSet)
	bpr.UserFactor = bpr.rng.NewNormalMatrix(trainSet.UserCount(), bpr.nFactors, bpr.initMean, bpr.initStdDev)
	bpr.ItemFactor = bpr.rng.NewNormalMatrix(trainSet.ItemCount(), bpr.nFactors, bpr.initMean, bpr.initStdDev)
	bpr.UserRatings = trainSet.Users()
	bpr.ItemPop = NewItemPop(nil)
	bpr.ItemPop.Fit(trainSet)
	temp := make([]float64, bpr.nFactors)
	userFactor := make([]float64, bpr.nFactors)
	positiveItemFactor := make([]float64, bpr.nFactors)
	negativeItemFactor := make([]float64, bpr.nFactors)
	for epoch := 0; epoch < bpr.nEpochs; epoch++ {
		for i := 0; i < trainSet.Count(); i++ {
			var userIndex, ratingCount int
			for {
				userIndex = bpr.rng.Intn(trainSet.UserCount())
				ratingCount = trainSet.UserByIndex(userIndex).Len()
				if ratingCount > 0 {
					break
				}
			}
			posIndex := trainSet.UserByIndex(userIndex).GetIndex(bpr.rng.Intn(ratingCount))
			// Select a negative sample
			negIndex := -1
			for {
				temp := bpr.rng.Intn(trainSet.ItemCount())
				tempId := bpr.ItemIndexer.ToID(temp)
				if !trainSet.UserByIndex(userIndex).Contain(tempId) {
					negIndex = temp
					break
				}
			}
			diff := bpr.predict(userIndex, posIndex) - bpr.predict(userIndex, negIndex)
			grad := math.Exp(-diff) / (1.0 + math.Exp(-diff))
			copy(userFactor, bpr.UserFactor[userIndex])
			copy(positiveItemFactor, bpr.ItemFactor[posIndex])
			copy(negativeItemFactor, bpr.ItemFactor[negIndex])
			// Update positive item latent factor: +w_u
			floats.MulConstTo(userFactor, grad, temp)
			floats.MulConstAddTo(positiveItemFactor, -bpr.reg, temp)
			floats.MulConstAddTo(temp, bpr.lr, bpr.ItemFactor[posIndex])
			// Update negative item latent factor: -w_u
			floats.MulConstTo(userFactor, -grad, temp)
			floats.MulConstAddTo(negativeItemFactor, -bpr.reg, temp)
			floats.MulConstAddTo(temp, bpr.lr, bpr.ItemFactor[negIndex])
			// Update user latent factor: h_i-h_j
			floats.SubTo(positiveItemFactor, negativeItemFactor, temp)
			floats.MulConst(temp, grad)
			floats.MulConstAddTo(userFactor, -bpr.reg, temp)
			floats.MulConstAddTo(temp, bpr.lr, bpr.UserFactor[userIndex])
		}
	}
}

func (bpr *BPR) predict(userIndex int, itemIndex int) float64 {
	ret := 0.0
	// + q_i^Tp_u
	if itemIndex != core.NotId && userIndex != core.NotId {
		userFactor := bpr.UserFactor[userIndex]
		itemFactor := bpr.ItemFactor[itemIndex]
		ret += floats.Dot(userFactor, itemFactor)
	}
	return ret
}
