package core

import (
	"github.com/liwei1dao/lego/sys/log"
)

type Option func(*Options)
type Options struct {
	IndentionStep                 int    //缩进步骤
	ObjectFieldMustBeSimpleString bool   //对象字段必须是简单字符串
	OnlyTaggedField               bool   //仅仅处理标签字段
	DisallowUnknownFields         bool   //禁止未知字段
	CaseSensitive                 bool   //是否区分大小写
	TagKey                        string //标签
	Debug                         bool   //日志是否开启
	Log                           log.ILog
}

//执行选项
type ExecuteOption func(*ExecuteOptions)
type ExecuteOptions struct {
}

func NewExecuteOptions(opts ...ExecuteOption) ExecuteOptions {
	options := ExecuteOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}
