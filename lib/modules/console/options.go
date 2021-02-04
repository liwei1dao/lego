package console

import (
	"github.com/liwei1dao/lego/lib/modules/http"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Options struct {
	http.Options
	RedisUrl                  string                 `json:"-"` //缓存地址 不显示控制台配置信息中
	UserCacheExpirationDate   int                    `json:"-"` //用户缓存过期时间 单位秒
	TokenCacheExpirationDate  int                    `json:"-"` //Token缓存过期时间 单位秒
	MonitorTotalTime          int                    `json:"-"` //监控总时长 小时为单位
	MongodbUrl                string                 `json:"-"` //数据库
	MongodbDatabase           string                 `json:"-"` //数据集
	MailServerhost            string                 `json:"-"` //邮件系统配置
	MailFromemail             string                 `json:"-"` //邮件系统配置
	MailFompasswd             string                 `json:"-"` //邮件系统配置
	MailServerport            int                    `json:"-"` //邮件系统配置
	MailCaptchaExpirationdate string                 `json:"-"` //邮件系统配置
	CaptchaExpirationdate     int                    `json:"-"` //Captcha缓存过期时间 单位秒
	SignKey                   string                 `json:"-"` //签名密钥
	UserInitialPassword       string                 `json:"-"` //用户初始密码
	ProjectName               string                 //项目名称
	ProjectDes                string                 //项目描述
	ProjectVersion            float64                //版本
	ProjectTime               string                 //立项时间
	ProjectMember             map[string]interface{} //项目经理 成员-职务
}

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
	if err = this.Options.LoadConfig(settings); err == nil {
		if settings != nil {
			err = mapstructure.Decode(settings, this)
		}
	}
	return
}
