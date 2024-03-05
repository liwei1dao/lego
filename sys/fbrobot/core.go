package fbrobot

import "context"

var (
	//GraphAPI specifies host used for API requests
	GraphAPI = "https://graph.facebook.com"
)

/*
系统描述:facebook 聊天机器人
*/
type (
	//sdk相关回处理回调函数
	MessageReceivedHandler  func(context.Context, Event, MessageOpts, ReceivedMessage)
	MessageDeliveredHandler func(context.Context, Event, MessageOpts, Delivery)
	PostbackHandler         func(context.Context, Event, MessageOpts, Postback)
	AuthenticationHandler   func(context.Context, Event, MessageOpts, *Optin)

	//错误对象
	rawError struct {
		Error Error `json:"error"`
	}
	Error struct {
		Message   string `json:"message"`
		Type      string `json:"type"`
		Code      int    `json:"code"`
		ErrorData string `json:"error_data"`
		TraceID   string `json:"fbtrace_id"`
	}

	//个人信息
	Profile struct {
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		ProfilePicture string `json:"profile_pic,omitempty"`
		Locale         string `json:"locale,omitempty"`
		Timezone       int    `json:"timezone,omitempty"`
		Gender         string `json:"gender,omitempty"`
	}
	MessageResponse struct {
		RecipientID string `json:"recipient_id"`
		MessageID   string `json:"message_id"`
	}
	//系统接口
	ISys interface {
		//获取用户信息
		GetProfile(ctx context.Context, userID string) (*Profile, error)
		//发送消息
		SendSimpleMessage(ctx context.Context, recipient string, message string) (*MessageResponse, error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}
