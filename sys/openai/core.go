package openai

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type (
	ISys interface {
		//非流式
		SendReq(ctx context.Context, content string) (results string, err error)
		//流式数据回应
		SendReqByStream(ctx context.Context, content string) (stream *openai.ChatCompletionStream, err error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func SendReq(ctx context.Context, content string) (results string, err error) {
	return defsys.SendReq(ctx, content)
}
