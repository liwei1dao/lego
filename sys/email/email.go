package email

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

func newSys(options Options) (sys ISys, err error) {
	var (
		host string
		addr string
	)
	switch options.EmailType {
	case GmailEmail:
		host = "smtp.gmail.com"
		addr = "smtp.gmail.com:587"
		break
	case QQEmail:
		host = "smtp.qq.com"
		addr = "smtp.qq.com:25"
		break
	}
	sys = &Email{
		From: options.FromEmail,
		Addr: addr,
		Auth: smtp.PlainAuth("", options.FromEmail, options.Password, host),
	}
	return
}

type Email struct {
	From string
	Addr string
	Auth smtp.Auth
}

//发送邮件
func (this *Email) SendMail(title string, content string, to ...string) (err error) {
	e := email.NewEmail()
	e.From = this.From
	e.To = to
	// e.Bcc = []string{"test_bcc@example.com"}
	// e.Cc = []string{"test_cc@example.com"}
	e.Subject = title
	e.Text = []byte(content)
	e.HTML = []byte("")
	err = e.Send(this.Addr, this.Auth)
	return
}
