package model

type ParamName string
type Params map[ParamName]interface{}

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
