package fbrobot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

func newSys(options *Options) (sys *FBRobot, err error) {
	sys = &FBRobot{
		option: options,
		fbapi: &messenger.Messenger{
			VerifyToken: options.VerifyToken,
			AppSecret:   options.AppSecret,
			AccessToken: options.AccessToken,
			PageID:      options.PageID,
		},
	}
	go sys.run()
	return
}

type FBRobot struct {
	option *Options
	fbapi  *messenger.Messenger
}

func (this *FBRobot) run() {
	http.HandleFunc("/webhook", this.Handler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", this.option.ListenPort), nil); err != nil {
		this.option.Log.Errorln(err)
	}
}

// 接收facebook messager的消息
func (this *FBRobot) Handler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		query := req.URL.Query()
		verifyToken := query.Get("hub.verify_token")
		this.option.Log.Debugf("Handle", req.Method, " token=", verifyToken, " verify_token=", this.option.VerifyToken)
		if verifyToken != this.option.VerifyToken {
			rw.WriteHeader(http.StatusUnauthorized)
			this.option.Log.Debugf("StatusUnauthorized")
			return
		}
		rw.WriteHeader(http.StatusOK)
		this.option.Log.Debugf("RET:", query.Get("hub.challenge"))
		rw.Write([]byte(query.Get("hub.challenge")))
	} else if req.Method == "POST" {
		this.handlePOST(rw, req)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func (this *FBRobot) handlePOST(rw http.ResponseWriter, req *http.Request) {
	read, err := ioutil.ReadAll(req.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if this.option.AppSecret != "" {
		if req.Header.Get("x-hub-signature") == "" || !checkIntegrity(this.option.AppSecret, read, req.Header.Get("x-hub-signature")[5:]) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	event := &upstreamEvent{}
	err = json.Unmarshal(read, event)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, entry := range event.Entries {
		for _, message := range entry.Messaging {
			if message.Delivery != nil {
				if this.option.MessageDelivered != nil {
					go this.option.MessageDelivered(req.Context(), entry.Event, message.MessageOpts, *message.Delivery)
				}
			} else if message.Message != nil {
				if this.option.MessageReceived != nil {
					go this.option.MessageReceived(req.Context(), entry.Event, message.MessageOpts, *message.Message)
				}
			} else if message.Postback != nil {
				if this.option.Postback != nil {
					go this.option.Postback(req.Context(), entry.Event, message.MessageOpts, *message.Postback)
				}
			} else if this.option.Authentication != nil {
				go this.option.Authentication(req.Context(), entry.Event, message.MessageOpts, message.Optin)
			}
		}
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))
}

// 获取个人信息
func (this *FBRobot) GetProfile(ctx context.Context, userID string) (*Profile, error) {
	resp, err := this.doRequest(ctx, "GET", fmt.Sprintf(GraphAPI+"/v2.6/%s?fields=first_name,last_name,profile_pic,locale,timezone,gender", userID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	read, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		er := new(rawError)
		json.Unmarshal(read, er)
		return nil, errors.New("Error occured: " + er.Error.Message)
	}
	profile := new(Profile)
	return profile, json.Unmarshal(read, profile)
}

// 请求
func (this *FBRobot) doRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	query := req.URL.Query()
	query.Set("access_token", this.option.AccessToken)
	req.URL.RawQuery = query.Encode()
	return http.DefaultClient.Do(req)
}

// 发送简单消息
func (this *FBRobot) SendSimpleMessage(ctx context.Context, recipient string, message string) (*MessageResponse, error) {
	return this.SendMessage(ctx, MessageQuery{
		Recipient: Recipient{
			ID: recipient,
		},
		Message: SendMessage{
			Text: message,
		},
	})
}

func (this *FBRobot) SendMessage(ctx context.Context, mq MessageQuery) (*MessageResponse, error) {
	byt, err := json.Marshal(mq)
	if err != nil {
		return nil, err
	}
	resp, err := this.doRequest(ctx, "POST", GraphAPI+"/v2.6/me/messages", bytes.NewReader(byt))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	read, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		er := new(rawError)
		json.Unmarshal(read, er)
		return nil, errors.New("Error occured: " + er.Error.Message)
	}
	response := &MessageResponse{}
	err = json.Unmarshal(read, response)
	return response, err
}
