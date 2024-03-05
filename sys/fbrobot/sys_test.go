package fbrobot_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/go-ego/gse"
	"github.com/liwei1dao/lego/sys/fbrobot"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/nlp"
)

var (
	sys fbrobot.ISys
)

func Test_Sys(t *testing.T) {
	var err error

	if err = log.OnInit(nil, log.SetFileName("./kafka.log")); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}

	if err = nlp.OnInit(nil); err != nil {
		fmt.Printf("nlp init err:%v", err)
		return
	}

	if sys, err = fbrobot.NewSys(
		fbrobot.Set_ListenPort(9898),
		fbrobot.Set_VerifyToken("your facebook VerifyToken"),
		fbrobot.Set_AppSecret("your facebook AppSecret"),
		fbrobot.Set_AccessToken("your facebook AccessToken"),
		fbrobot.Set_PageID("your page"),
		fbrobot.Set_MessageReceived(MessageReceived),
	); err == nil {
		log.Debugf("fbrobot new fail!")
		return
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Debugf("test end signal")
	}

}

// 接收用户消息
func MessageReceived(ctx context.Context, event fbrobot.Event, opts fbrobot.MessageOpts, msg fbrobot.ReceivedMessage) (err error) {
	var (
		pos     []gse.SegPos
		profile *fbrobot.Profile
	)
	pos = nlp.Pos(msg.Text) //cixici
	for _, v := range pos {
		if v.Text == "你好" && v.Pos == "l" {
			profile, err = sys.GetProfile(ctx, opts.Sender.ID)
			if err != nil {
				log.Debugln(err)
				return
			}
			_, err = sys.SendSimpleMessage(ctx, opts.Sender.ID, fmt.Sprintf("Hello   , %s %s, %s", profile.FirstName, profile.LastName, msg.Text))
			if err != nil {
				log.Debugln(err)
			}
		}
	}
	return

}
