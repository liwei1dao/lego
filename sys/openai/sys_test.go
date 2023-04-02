package openai_test

import (
	"context"
	"fmt"
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

//测试OpenAI
func Test_OpenAI(t *testing.T) {
	client := openai.NewClient("you Token")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "你好! 你的名字是？能告诉我 1+2=?",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
