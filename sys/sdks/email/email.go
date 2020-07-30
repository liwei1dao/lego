package email

import (
	"fmt"

	"github.com/go-gomail/gomail"
)

func NewEmail(serverhost, fromemail, fompasswd string, serverport int) *sdk_email {
	return &sdk_email{serverhost, serverport, fromemail, fompasswd}
}

type sdk_email struct {
	serverhost string
	serverport int
	fromemail  string
	fompasswd  string
}

// SendEmail body支持html格式字符串
func (this *sdk_email) SendCaptcha(email string, captcha string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", this.fromemail, "")
	m.SetHeader("To", email)
	// 主题
	m.SetHeader("Subject", "一刀工作室")
	// 正文
	m.SetBody("text/html", fmt.Sprintf(`验证码<br>
	<h3>您的验证码是:%s</h3><br>`, captcha))
	d := gomail.NewPlainDialer(this.serverhost, this.serverport, this.fromemail, this.fompasswd)
	// 发送
	err := d.DialAndSend(m)
	return err
}
