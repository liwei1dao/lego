package console

import (
	"encoding/json"

	"github.com/liwei1dao/lego/lib/modules/http"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type IOptions interface {
	http.IOptions
	GetRedisUrl() string
	GetUserCacheExpirationDate() int
	GetTokenCacheExpirationDate() int
	GetMonitorTotalTime() int
	GetMongodbUrl() string
	GetMongodbDatabase() string
	GetMailServerhost() string
	GetMailFromemail() string
	GetMailFompasswd() string
	GetMailServerport() int
	GetCaptchaExpirationdate() int
	GetSignKey() string
	GetUserInitialPassword() string
	GetProjectName() string
	ToString() string
}

type Options struct {
	http.Options
	RedisUrl                 string                 `json:"-"` //缓存地址 不显示控制台配置信息中
	UserCacheExpirationDate  int                    `json:"-"` //用户缓存过期时间 单位秒
	TokenCacheExpirationDate int                    `json:"-"` //Token缓存过期时间 单位秒
	MonitorTotalTime         int                    `json:"-"` //监控总时长 小时为单位
	MongodbUrl               string                 `json:"-"` //数据库
	MongodbDatabase          string                 `json:"-"` //数据集
	MailServerhost           string                 `json:"-"` //邮件系统配置
	MailFromemail            string                 `json:"-"` //邮件系统配置
	MailFompasswd            string                 `json:"-"` //邮件系统配置
	MailServerport           int                    `json:"-"` //邮件系统配置
	CaptchaExpirationdate    int                    `json:"-"` //Captcha缓存过期时间 单位秒
	SignKey                  string                 `json:"-"` //签名密钥
	UserInitialPassword      string                 `json:"-"` //用户初始密码
	ProjectName              string                 //项目名称
	ProjectDes               string                 //项目描述
	ProjectVersion           float64                //版本
	ProjectTime              string                 //立项时间
	ProjectMember            map[string]interface{} //项目经理 成员-职务
}

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
	this.UserCacheExpirationDate = 60
	this.TokenCacheExpirationDate = 3600
	this.CaptchaExpirationdate = 60
	if err = this.Options.LoadConfig(settings); err == nil {
		if settings != nil {
			err = mapstructure.Decode(settings, this)
		}
	}
	return
}

func (this *Options) GetRedisUrl() string {
	return this.RedisUrl
}

func (this *Options) GetUserCacheExpirationDate() int {
	return this.UserCacheExpirationDate
}
func (this *Options) GetTokenCacheExpirationDate() int {
	return this.TokenCacheExpirationDate
}
func (this *Options) GetMonitorTotalTime() int {
	return this.MonitorTotalTime
}
func (this *Options) GetMongodbUrl() string {
	return this.MongodbUrl
}
func (this *Options) GetMongodbDatabase() string {
	return this.MongodbDatabase
}
func (this *Options) GetMailServerhost() string {
	return this.MailServerhost
}
func (this *Options) GetMailFromemail() string {
	return this.MailFromemail
}
func (this *Options) GetMailFompasswd() string {
	return this.MailFompasswd
}
func (this *Options) GetMailServerport() int {
	return this.MailServerport
}

func (this *Options) GetCaptchaExpirationdate() int {
	return this.CaptchaExpirationdate
}
func (this *Options) GetSignKey() string {
	return this.SignKey
}
func (this *Options) GetUserInitialPassword() string {
	return this.UserInitialPassword
}

func (this *Options) GetProjectName() string {
	return this.ProjectName
}

func (this *Options) ToString() string {
	str, _ := json.Marshal(this)
	return string(str)
}
