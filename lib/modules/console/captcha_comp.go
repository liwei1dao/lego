package console

import (
	"fmt"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/sdks/email"
)

type CaptchaComp struct {
	cbase.ModuleCompBase
	module *Console
	email  email.IEmail
}

func (this *CaptchaComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.module = module.(*Console)
	this.email, err = email.NewSys(
		email.SetServerhost(this.module.options.MailServerhost),
		email.SetFromemail(this.module.options.MailFromemail),
		email.SetFompasswd(this.module.options.MailFompasswd),
		email.SetServerport(this.module.options.MailServerport))
	return
}

//发送邮箱验证码
func (this *CaptchaComp) SendEmailCaptcha(email, captcha string) (err error) {
	return this.email.SendMail(email, "Verification Code", fmt.Sprintf("Your %s console verification code:<b>%s</b>, please do not forward others", this.module.options.ProjectName, captcha))
}

//获取验证码
func (this *CaptchaComp) QueryCaptcha(cId string, ctype CaptchaType) (captcha string, err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleCaptcha), cId, ctype)
	pool := this.module.cache.GetPool()
	err = pool.GetKeyForValue(Id, &captcha)
	return
}

//写入验证码
func (this *CaptchaComp) WriteCaptcha(cId, captcha string, ctype CaptchaType) {
	Id := fmt.Sprintf(string(Cache_ConsoleCaptcha), cId, ctype)
	pool := this.module.cache.GetPool()
	pool.SetExKeyForValue(Id, captcha, this.module.options.CaptchaExpirationdate)
}
