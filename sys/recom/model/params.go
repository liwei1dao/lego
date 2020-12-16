package model

import (
	"reflect"

	"github.com/liwei1dao/lego/sys/log"
)

type ParamName string

const (
	Lr            ParamName = "Lr"            // learning rate	学习率
	Reg           ParamName = "Reg"           // regularization strength 正则化强度
	NEpochs       ParamName = "NEpochs"       // number of epochs	纪元数
	NFactors      ParamName = "NFactors"      // number of factors	因素数量
	RandomState   ParamName = "RandomState"   // random state (seed)	随机状态（种子）
	UseBias       ParamName = "UseBias"       // use bias				使用偏见
	InitMean      ParamName = "InitMean"      // mean of gaussian initial parameter	高斯初始参数的平均值
	InitStdDev    ParamName = "InitStdDev"    // standard deviation of gaussian initial parameter	高斯初始参数的标准偏差
	InitLow       ParamName = "InitLow"       // lower bound of uniform initial parameter			统一初始参数的下界
	InitHigh      ParamName = "InitHigh"      // upper bound of uniform initial parameter			统一初始参数的上限
	NUserClusters ParamName = "NUserClusters" // number of user cluster								用户集群数
	NItemClusters ParamName = "NItemClusters" // number of item cluster								项目群数
	Type          ParamName = "Type"          // type for KNN										KNN类型
	UserBased     ParamName = "UserBased"     // user based if true. otherwise item based.			基于用户（如果为true）。 否则基于项目。
	Similarity    ParamName = "Similarity"    // similarity metrics									相似性指标
	K             ParamName = "K"             // number of neighbors								邻居数
	MinK          ParamName = "MinK"          // least number of neighbors							邻居最少
	Optimizer     ParamName = "Optimizer"     // optimizer for optimization (SGD/ALS/BPR)			用于优化的优化器（SGD / ALS / BPR）
	Shrinkage     ParamName = "Shrinkage"     // shrinkage strength of similarity					相似收缩强度
	Alpha         ParamName = "Alpha"         // alpha value, depend on context						alpha值，取决于上下文
)

type Params map[ParamName]interface{}

// Copy hyper-parameters.
func (parameters Params) Copy() Params {
	newParams := make(Params)
	for k, v := range parameters {
		newParams[k] = v
	}
	return newParams
}

// GetInt gets a integer parameter by name. Returns _default if not exists or type doesn't match.
func (parameters Params) GetInt(name ParamName, _default int) int {
	if val, exist := parameters[name]; exist {
		switch val.(type) {
		case int:
			return val.(int)
		default:
			log.Infof("Expect %v to be int, but get %v", name, reflect.TypeOf(name))
		}
	}
	return _default
}

// GetInt64 gets a int64 parameter by name. Returns _default if not exists or type doesn't match. The
// type will be converted if given int.
func (parameters Params) GetInt64(name ParamName, _default int64) int64 {
	if val, exist := parameters[name]; exist {
		switch val.(type) {
		case int64:
			return val.(int64)
		case int:
			return int64(val.(int))
		default:
			log.Infof("Expect %v to be int, but get %v", name, reflect.TypeOf(name))
		}
	}
	return _default
}

// GetBool gets a bool parameter by name. Returns _default if not exists or type doesn't match.
func (parameters Params) GetBool(name ParamName, _default bool) bool {
	if val, exist := parameters[name]; exist {
		switch val.(type) {
		case bool:
			return val.(bool)
		default:
			log.Infof("Expect %v to be int, but get %v", name, reflect.TypeOf(name))
		}
	}
	return _default
}

// GetFloat64 gets a float parameter by name. Returns _default if not exists or type doesn't match. The
// type will be converted if given int.
func (parameters Params) GetFloat64(name ParamName, _default float64) float64 {
	if val, exist := parameters[name]; exist {
		switch val.(type) {
		case float64:
			return val.(float64)
		case int:
			return float64(val.(int))
		default:
			log.Infof("Expect %v to be int, but get %v", name, reflect.TypeOf(name))
		}
	}
	return _default
}

// GetString gets a string parameter. Returns _default if not exists or type doesn't match.
func (parameters Params) GetString(name ParamName, _default string) string {
	if val, exist := parameters[name]; exist {
		switch val.(type) {
		case string:
			return val.(string)
		default:
			log.Infof("Expect %v to be string, but get %v", name, reflect.TypeOf(name))
		}
	}
	return _default
}

// Merge another group of hyper-parameters to current group of hyper-parameters.
func (parameters Params) Merge(params Params) Params {
	merged := make(Params)
	for k, v := range parameters {
		merged[k] = v
	}
	for k, v := range params {
		merged[k] = v
	}
	return merged
}
