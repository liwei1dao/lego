package email_test

import (
	"fmt"
	"net/smtp"
	"testing"

	eml "github.com/jordan-wright/email"
)

//发送邮件
func Test_SendMail_Gamil(t *testing.T) {
	e := eml.NewEmail()
	e.From = "liwei1dao@gmail.com"
	e.To = []string{"2211068034@qq.com"}
	// e.Bcc = []string{"test_bcc@example.com"}
	// e.Cc = []string{"test_cc@example.com"}
	e.Subject = "Awesome Subject"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("")
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "liwei1dao@gmail.com", "afsqdeweopnpakuw", "smtp.gmail.com"))
	fmt.Printf("Test_SendMail:%v", err)
}

//发送邮件
func Test_SendMail_QQ(t *testing.T) {
	e := eml.NewEmail()
	e.From = "2211068034@qq.com"
	e.To = []string{"2211068034@qq.com"}
	// e.Bcc = []string{"test_bcc@example.com"}
	// e.Cc = []string{"test_cc@example.com"}
	e.Subject = "Awesome Subject"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("")
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2211068034@qq.com", "ffngwepbmewcdifc", "smtp.qq.com"))
	fmt.Printf("Test_SendMail:%v", err)
}
