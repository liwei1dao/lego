package console

import (
	"fmt"
	nethttp "net/http"
	"strings"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib"
	"github.com/liwei1dao/lego/sys/gin/engine"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/crypto/aes"
)

type Console struct {
	cbase.ModuleBase
	options            IOptions
	db                 *DBComp
	cache              *CacheComp
	captcha            *CaptchaComp
	hostmonitorcomp    *HostMonitorComp
	clustermonitorcomp *ClusterMonitorComp
	apiComp            *ApiComp
}

func (this *Console) GetType() core.M_Modules {
	return lib.SM_ConsoleModule
}

func (this *Console) NewOptions() (options core.IModuleOptions) {
	return new(Options)
}

func (this *Console) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	this.options = options.(IOptions)
	err = this.ModuleBase.Init(service, module, options)
	return
}

func (this *Console) OnInstallComp() {
	this.ModuleBase.OnInstallComp()
	this.db = this.RegisterComp(new(DBComp)).(*DBComp)
	this.cache = this.RegisterComp(new(CacheComp)).(*CacheComp)
	this.captcha = this.RegisterComp(new(CaptchaComp)).(*CaptchaComp)
	this.hostmonitorcomp = this.RegisterComp(new(HostMonitorComp)).(*HostMonitorComp)
	this.clustermonitorcomp = this.RegisterComp(new(ClusterMonitorComp)).(*ClusterMonitorComp)
	this.apiComp = this.RegisterComp(new(ApiComp)).(*ApiComp)
}

func (this *Console) Options() IOptions {
	return this.options
}

func (this *Console) Cache() ICache {
	return this.cache
}

func (this *Console) DB() IDB {
	return this.db
}
func (this *Console) Captcha() ICaptcha {
	return this.captcha
}

func (this *Console) Hostmonitorcomp() IHostmonitorcomp {
	return this.hostmonitorcomp
}

func (this *Console) Clustermonitorcomp() IClustermonitorcomp {
	return this.clustermonitorcomp
}

//token 验证中间键
func (this *Console) CheckToken(c *engine.Context) {
	token := c.Request.Header.Get("X-Token")
	if token == "" {
		this.HttpStatusOK(c, ErrorCode_NoLOgin, nil)
		c.Abort()
		return
	}
	uId, err := this.cache.QueryToken(token)
	if err != nil {
		this.HttpStatusOK(c, ErrorCode_NoLOgin, nil)
		c.Abort()
		return
	}
	c.Set(UserKey, uId)
	c.Next()
}

//统一输出接口
func (this *Console) HttpStatusOK(c *engine.Context, code core.ErrorCode, data interface{}) {
	defer log.Debugf("%s resp: code:%d data:%+v", c.Request.RequestURI, code, data)
	c.JSON(nethttp.StatusOK, OutJson{ErrorCode: code, Message: GetErrorCodeMsg(code), Data: data})
}

//创建token
func (this *Console) CreateToken(uId uint32) (token string) {
	orsign := fmt.Sprintf("%d|%d", uId, time.Now().Unix())
	token = aes.AesEncryptCBC(orsign, this.options.GetSignKey())
	return
}

//校验token
func (this *Console) checkToken(token string, uId uint32) (check bool) {
	orsign := aes.AesDecryptCBC(token, this.options.GetSignKey())
	tempstr := strings.Split(orsign, "|")
	if len(tempstr) == 2 && tempstr[0] == fmt.Sprintf("%d", uId) {
		return true
	}
	return false
}
