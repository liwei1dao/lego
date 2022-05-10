package console

import (
	"encoding/json"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type IOptions interface {
	GetListenPort() int
	GetCertPath() string
	GetKeyPath() string
	GetCors() bool
	GetSignKey() string
	GetRedisUrl() string
	GetRedisDB() int
	GetRedisPassword() string
	GetInitUserIdNum() int
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
	GetUserInitialPassword() string
	GetProjectName() string
	ToString() string
}

type Options struct {
	ListenPort               int                    `json:"-"`
	CertPath                 string                 `json:"-"`
	KeyPath                  string                 `json:"-"`
	Cors                     bool                   `json:"-"` //是否跨域
	SignKey                  string                 `json:"-"` //签名密钥
	RedisUrl                 string                 `json:"-"` //缓存地址 不显示控制台配置信息中
	RedisDB                  int                    `json:"-"` //缓存DB
	RedisPassword            string                 `json:"-"` //缓存密码
	InitUserIdNum            int                    `json:"-"` //初始化用户Id数
	UserCacheExpirationDate  int                    `json:"-"` //用户缓存过期时间 单位秒
	TokenCacheExpirationDate int                    `json:"-"` //Token缓存过期时间 单位秒
	MonitorTotalTime         int                    `json:"-"` //监控总时长 小时为单位
	MongodbUrl               string                 `json:"-"` //数据库
	MongodbDatabase          string                 `json:"-"` //数据集
	MailFromemail            string                 `json:"-"` //邮件系统配置
	MailFompasswd            string                 `json:"-"` //邮件系统配置
	CaptchaExpirationdate    int                    `json:"-"` //Captcha缓存过期时间 单位秒
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
	this.InitUserIdNum = 100000
	if settings != nil {
		err = mapstructure.Decode(settings, this)
	}
	return
}
func (this *Options) GetListenPort() int {
	return this.ListenPort
}
func (this *Options) GetCertPath() string {
	return this.CertPath
}
func (this *Options) GetKeyPath() string {
	return this.KeyPath
}
func (this *Options) GetCors() bool {
	return this.Cors
}
func (this *Options) GetSignKey() string {
	return this.SignKey
}
func (this *Options) GetRedisUrl() string {
	return this.RedisUrl
}
func (this *Options) GetRedisDB() int {
	return this.RedisDB
}

func (this *Options) GetRedisPassword() string {
	return this.RedisPassword
}
func (this *Options) GetInitUserIdNum() int {
	return this.InitUserIdNum
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

func (this *Options) GetMailFromemail() string {
	return this.MailFromemail
}
func (this *Options) GetMailFompasswd() string {
	return this.MailFompasswd
}

func (this *Options) GetCaptchaExpirationdate() int {
	return this.CaptchaExpirationdate
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
