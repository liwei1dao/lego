package model

import (
	"github.com/liwei1dao/lego/sys/recom/core"
	"github.com/liwei1dao/lego/sys/recom/floats"
)

func NewSVD(params Params) *SVD {
	svd := new(SVD)
	svd.SetParams(params)
	return svd
}

type SVD struct {
	ModelBase
	UserFactor [][]float64 // p_u
	ItemFactor [][]float64 // q_i
	UserBias   []float64   // b_u
	ItemBias   []float64   // b_i
	GlobalMean float64     // mu
	// Hyper parameters
	useBias    bool
	nFactors   int
	nEpochs    int
	lr         float64
	reg        float64
	initMean   float64
	initStdDev float64
	optimizer  string
	// Fallback model
	UserRatings []*core.MarginalSubSet
	ItemPop     *ItemPop
}

func (svd *SVD) SetParams(params Params) {
	svd.ModelBase.SetParams(params)
	// Setup hyper-parameters
	svd.useBias = svd.Params.GetBool(UseBias, true)
	svd.nFactors = svd.Params.GetInt(NFactors, 100)
	svd.nEpochs = svd.Params.GetInt(NEpochs, 20)
	svd.lr = svd.Params.GetFloat64(Lr, 0.005)
	svd.reg = svd.Params.GetFloat64(Reg, 0.02)
	svd.initMean = svd.Params.GetFloat64(InitMean, 0)
	svd.initStdDev = svd.Params.GetFloat64(InitStdDev, 0.1)
}

// Predict by the SVD model.
func (svd *SVD) Predict(userId, itemId uint32) float64 {
	// Convert sparse IDs to dense IDs
	userIndex := svd.UserIndexer.ToIndex(userId)
	itemIndex := svd.ItemIndexer.ToIndex(itemId)
	return svd.predict(userIndex, itemIndex)
}

func (svd *SVD) predict(userIndex int, itemIndex int) float64 {
	ret := svd.GlobalMean
	// + b_u
	if userIndex != core.NotId {
		ret += svd.UserBias[userIndex]
	}
	// + b_i
	if itemIndex != core.NotId {
		ret += svd.ItemBias[itemIndex]
	}
	// + q_i^Tp_u
	if itemIndex != core.NotId && userIndex != core.NotId {
		userFactor := svd.UserFactor[userIndex]
		itemFactor := svd.ItemFactor[itemIndex]
		ret += floats.Dot(userFactor, itemFactor)
	}
	return ret
}

func (svd *SVD) Fit(trainSet core.DataSetInterface) {
	svd.Init(trainSet)
	// Initialize parameters
	svd.GlobalMean = 0
	svd.UserBias = make([]float64, trainSet.UserCount())
	svd.ItemBias = make([]float64, trainSet.ItemCount())
	svd.UserFactor = svd.rng.NewNormalMatrix(trainSet.UserCount(), svd.nFactors, svd.initMean, svd.initStdDev)
	svd.ItemFactor = svd.rng.NewNormalMatrix(trainSet.ItemCount(), svd.nFactors, svd.initMean, svd.initStdDev)
	svd.GlobalMean = trainSet.GlobalMean()
	// Create buffers
	temp := make([]float64, svd.nFactors)
	userFactor := make([]float64, svd.nFactors)
	itemFactor := make([]float64, svd.nFactors)
	// Optimize
	for epoch := 0; epoch < svd.nEpochs; epoch++ {
		loss := 0.0
		for i := 0; i < trainSet.Count(); i++ {
			userIndex, itemIndex, rating := trainSet.GetWithIndex(i)
			// Compute error: e_{ui} = r - \hat r
			upGrad := rating - svd.predict(userIndex, itemIndex)
			loss += upGrad * upGrad
			if svd.useBias {
				userBias := svd.UserBias[userIndex]
				itemBias := svd.ItemBias[itemIndex]
				// Update user Bias: b_u <- b_u + \gamma (e_{ui} - \lambda b_u)
				gradUserBias := upGrad - svd.reg*userBias
				svd.UserBias[userIndex] += svd.lr * gradUserBias
				// Update item Bias: p_i <- p_i + \gamma (e_{ui} - \lambda b_i)
				gradItemBias := upGrad - svd.reg*itemBias
				svd.ItemBias[itemIndex] += svd.lr * gradItemBias
			}
			copy(userFactor, svd.UserFactor[userIndex])
			copy(itemFactor, svd.ItemFactor[itemIndex])
			// Update user latent factor
			floats.MulConstTo(itemFactor, upGrad, temp)
			floats.MulConstAddTo(userFactor, -svd.reg, temp)
			floats.MulConstAddTo(temp, svd.lr, svd.UserFactor[userIndex])
			// Update item latent factor
			floats.MulConstTo(userFactor, upGrad, temp)
			floats.MulConstAddTo(itemFactor, -svd.reg, temp)
			floats.MulConstAddTo(temp, svd.lr, svd.ItemFactor[itemIndex])
		}
	}
}
