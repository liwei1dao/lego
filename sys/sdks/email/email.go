package email

import (
	"gopkg.in/gomail.v2"
)

func newSys(options Options) (sys *Email, err error) {
	sys = &Email{options.Serverhost, options.Serverport, options.Fromemail, options.Fompasswd}
	return
}

type Email struct {
	serverhost string
	serverport int
	fromemail  string
	fompasswd  string
}

//发送邮件
func (this *Email) SendMail(temail string, title string, content string) (err error) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", this.fromemail, "liwei1dao")
	m.SetHeader("To", temail)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)
	d := gomail.NewPlainDialer(this.serverhost, this.serverport, this.fromemail, this.fompasswd)
	err = d.DialAndSend(m)
	return
}
