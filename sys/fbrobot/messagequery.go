package fbrobot

import (
	"errors"
	"fmt"
)

type SendMessage struct {
	Text       string      `json:"text,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

// Recipient describes the person who will receive the message
// Either ID or PhoneNumber has to be set
type Recipient struct {
	ID          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// notificationType描述了手机收到消息后将执行的行为
type NotificationType string

const (
	//通知类型常规将发出声音/振动和电话通知
	NotificationTypeRegular NotificationType = "REGULAR"
	//NotificationTypeSilentPush 只会发出电话通知
	NotificationTypeSilentPush NotificationType = "SILENT_PUSH"
	//NotificationTypeNoPush 不会发出声音/振动，也不会发出电话通知
	NotificationTypeNoPush NotificationType = "NO_PUSH"
)

type MessageQuery struct {
	Recipient        Recipient        `json:"recipient"`
	Message          SendMessage      `json:"message"`
	NotificationType NotificationType `json:"notification_type,omitempty"`
}

func (this *MessageQuery) RecipientID(recipientID string) error {
	if this.Recipient.PhoneNumber != "" {
		return errors.New("Only one user identification (phone or id) can be specified.")
	}
	this.Recipient.ID = recipientID
	return nil
}

func (this *MessageQuery) RecipientPhoneNumber(phoneNumber string) error {
	if this.Recipient.ID != "" {
		return errors.New("Only one user identification (phone or id) can be specified.")
	}
	this.Recipient.PhoneNumber = phoneNumber
	return nil
}

func (this *MessageQuery) Notification(notification NotificationType) *MessageQuery {
	this.NotificationType = notification
	return this
}

func (this *MessageQuery) Text(text string) error {
	if this.Message.Attachment == nil {
		this.Message.Attachment = &Attachment{}
	}
	if this.Message.Attachment != nil && this.Message.Attachment.Type == AttachmentTypeTemplate {
		return errors.New("Can't set both text and template.")
	}
	this.Message.Text = text
	return nil
}

func (mq *MessageQuery) resource(typ AttachmentType, url string) (err error) {
	if mq.Message.Attachment == nil {
		mq.Message.Attachment = &Attachment{}
	}
	if mq.Message.Attachment.Payload != nil {
		err = fmt.Errorf("Attachment already specified.")
		return
	}
	mq.Message.Attachment.Type = typ
	mq.Message.Attachment.Payload = &Resource{URL: url}
	return
}

func (mq *MessageQuery) Audio(url string) error {
	return mq.resource(AttachmentTypeAudio, url)
}

func (mq *MessageQuery) Video(url string) error {
	return mq.resource(AttachmentTypeVideo, url)
}

func (mq *MessageQuery) Image(url string) error {
	return mq.resource(AttachmentTypeImage, url)
}
