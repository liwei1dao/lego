package openai

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

func newSys(options Options) (sys *OpenAI, err error) {
	sys = &OpenAI{}
	sys.client = openai.NewClient(options.Token)
	return
}

type OpenAI struct {
	client *openai.Client
}

//发送邮件
func (this *OpenAI) SendReq(content string) (results string, err error) {
	resp, err := this.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return
	}
	results = resp.Choices[0].Message.Content
	return
}
