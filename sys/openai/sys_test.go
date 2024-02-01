package openai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	lgopenai "github.com/liwei1dao/lego/sys/openai"
	openai "github.com/sashabaranov/go-openai"
)

//测试OpenAI
func Test_SendReq(t *testing.T) {
	var (
		sys  lgopenai.ISys
		resp string
		err  error
	)
	if sys, err = lgopenai.NewSys(lgopenai.SetToken("your-token")); err != nil {
		fmt.Println(err.Error())
		return
	}
	if resp, err = sys.SendReq(context.Background(), "我想测试下你的api接口"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

//测试流式方式
func Test_SendReqByStream(t *testing.T) {
	var (
		sys  lgopenai.ISys
		resp *openai.ChatCompletionStream
		err  error
	)
	if sys, err = lgopenai.NewSys(lgopenai.SetToken("your-token")); err != nil {
		fmt.Println(err.Error())
		return
	}
	if resp, err = sys.SendReqByStream(context.Background(), "我想测试下你的api接口"); err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Close()
	fmt.Printf("Stream response: ")
	for {
		response, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}
		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}
		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
