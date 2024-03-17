package openai

import (
	"context"

	"github.com/liwei1dao/lego/sys/log"
	openai "github.com/sashabaranov/go-openai"
)

func newSys(options *Options) (sys *OpenAI, err error) {
	sys = &OpenAI{
		options: options,
	}
	sys.client = openai.NewClient(options.Token)
	return
}

type OpenAI struct {
	options *Options
	client  *openai.Client
}

//发送消息 非流式模式
func (this *OpenAI) SendReq(ctx context.Context, content string) (results string, err error) {
	var (
		resp openai.ChatCompletionResponse
	)
	resp, err = this.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: this.options.Model,
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

//发送消息 流式模式
func (this *OpenAI) SendReqByStream(ctx context.Context, content string) (stream *openai.ChatCompletionStream, err error) {
	req := openai.ChatCompletionRequest{
		Model:     this.options.Model,
		MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
		Stream: true,
	}
	stream, err = this.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		this.options.Log.Error("CreateChatCompletionStream", log.Field{Key: "err", Value: err.Error()})
		return
	}
	return
}

//发送音频语音转文本
//filepath = "recording.mp3"
func (this *OpenAI) SendAudioToText(ctx context.Context, filepath string) (result string, err error) {
	var (
		req  openai.AudioRequest
		resp openai.AudioResponse
	)
	req = openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filepath,
	}
	resp, err = this.client.CreateTranscription(ctx, req)
	if err != nil {
		this.options.Log.Error("Transcription", log.Field{Key: "err", Value: err.Error()})
		return
	}
	result = resp.Text
	return
}
