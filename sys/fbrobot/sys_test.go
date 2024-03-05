package fbrobot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/event"
	"github.com/liwei1dao/lego/sys/fbrobot"
)

var (
	sys fbrobot.ISys
)

func Test_Sys(t *testing.T) {
	var err error
	if sys, err = fbrobot.NewSys(
		fbrobot.Set_ListenPort(9898),
		fbrobot.Set_VerifyToken(""),
		fbrobot.Set_AppSecret("9898"),
		fbrobot.Set_AccessToken("9898"),
		fbrobot.Set_PageID("9898"),
		fbrobot.Set_MessageReceived(MessageReceived),
	); err == nil {
		event.Register(core.Event_Key("TestEvent"), func() {
			fmt.Printf("TestEvent TriggerEvent")
		})
		event.TriggerEvent(core.Event_Key("TestEvent"))
	}
}

// 接收用户消息
func MessageReceived(ctx context.Context, event fbrobot.Event, opts fbrobot.MessageOpts, msg fbrobot.ReceivedMessage) {
	// log.Println("event:", event, " opt:", opts, " msg:", msg)
	profile, err := sys.GetProfile(ctx, opts.Sender.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := sys.SendSimpleMessage(ctx, opts.Sender.ID, fmt.Sprintf("Hello   , %s %s, %s", profile.FirstName, profile.LastName, msg.Text))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", resp)
}
