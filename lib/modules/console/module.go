package console

import (
	"fmt"
	nethttp "net/http"
	reflect "reflect"
	"sort"
	"strings"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/lib"
	"github.com/liwei1dao/lego/lib/modules/http"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/crypto"
)

type Console struct {
	http.Http
	options            *Options
	db                 *DBComp
	cache              *CacheComp
	captcha            *CaptchaComp
	hostmonitorcomp    *HostMonitorComp
	clustermonitorcomp *ClusterMonitorComp
}

func (this *Console) GetType() core.M_Modules {
	return lib.SM_ConsoleModule
}

func (this *Console) NewOptions() (options core.IModuleOptions) {
	return new(Options)
}

func (this *Console) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	this.options = options.(*Options)
	err = this.Http.Init(service, module, options)
	return
}

func (this *Console) OnInstallComp() {
	this.ModuleBase.OnInstallComp()
	this.db = this.RegisterComp(new(DBComp)).(*DBComp)
	this.cache = this.RegisterComp(new(CacheComp)).(*CacheComp)
	this.captcha = this.RegisterComp(new(CaptchaComp)).(*CaptchaComp)
	this.hostmonitorcomp = this.RegisterComp(new(HostMonitorComp)).(*HostMonitorComp)
	this.clustermonitorcomp = this.RegisterComp(new(ClusterMonitorComp)).(*ClusterMonitorComp)
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

//签名接口
func (this *Console) ParamSign(param map[string]interface{}) (sign string) {
	var keys []string
	for k, _ := range param {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	builder := strings.Builder{}
	for _, v := range keys {
		builder.WriteString(v)
		builder.WriteString("=")
		switch reflect.TypeOf(param[v]).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
			builder.WriteString(fmt.Sprintf("%d", param[v]))
			break
		default:
			builder.WriteString(fmt.Sprintf("%s", param[v]))
			break
		}
		builder.WriteString("&")
	}
	builder.WriteString("key=" + this.options.SignKey)
	log.Infof("orsign:%s", builder.String())
	sign = crypto.MD5EncToLower(builder.String())
	return
}

//token 验证中间键
func (this *Console) CheckToken(c *http.Context) {
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
func (this *Console) HttpStatusOK(c *http.Context, code core.ErrorCode, data interface{}) {
	defer log.Debugf("%s resp: code:%d data:%+v", c.Request.RequestURI, code, data)
	c.JSON(nethttp.StatusOK, http.OutJson{ErrorCode: code, Message: GetErrorCodeMsg(code), Data: data})
}

//创建token
func (this *Console) CreateToken(uId uint32) (token string) {
	orsign := fmt.Sprintf("%d|%d", uId, time.Now().Unix())
	token = crypto.AesEncryptCBC(orsign, this.options.SignKey)
	return
}

//校验token
func (this *Console) checkToken(token string, uId uint32) (check bool) {
	orsign := crypto.AesDecryptCBC(token, this.options.SignKey)
	tempstr := strings.Split(orsign, "|")
	if len(tempstr) == 2 && tempstr[0] == fmt.Sprintf("%d", uId) {
		return true
	}
	return false
}
