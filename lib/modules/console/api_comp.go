package console

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/http"
	"github.com/liwei1dao/lego/sys/log"
)

type ApiComp struct {
	cbase.ModuleCompBase
	module *Console
}

func (this *ApiComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.module = module.(*Console)
	api := this.module.Group("/lego/api")
	api.POST("/sendemailcaptcha", this.SendEmailCaptchaReq)

	console := this.module.Group("/lego/console")
	console.POST("/getprojectinfo", this.module.CheckToken, this.GetProjectInfo)
	console.POST("/gethostinfo", this.module.CheckToken, this.GetHostInfo)
	console.POST("/getcpuinfo", this.module.CheckToken, this.GetCpuInfo)
	console.POST("/getmemoryinfo", this.module.CheckToken, this.GetMemoryInfo)
	console.POST("/gethostmonitordata", this.module.CheckToken, this.GetHostMonitorData)
	console.POST("/getconsolecluster", this.module.CheckToken, this.GetClusterMonitorData)

	user := this.module.Group("/lego/user")
	user.POST("/loginbycaptcha", this.LoginByCaptchaReq)
	user.POST("/loginbypassword", this.LoginByPasswordReq)
	user.POST("/getuserinfo", this.module.CheckToken, this.GetUserinfoReq)
	user.POST("/loginout", this.module.CheckToken, this.LoginOutReq)

	return
}

/*				 Api相关接口
 * _______________#########_______________________
 * ______________############_____________________
 * ______________#############____________________
 * _____________##__###########___________________
 * ____________###__######_#####__________________
 * ____________###_#######___####_________________
 * ___________###__##########_####________________
 * __________####__###########_####_______________
 * ________#####___###########__#####_____________
 * _______######___###_########___#####___________
 * _______#####___###___########___######_________
 * ______######___###__###########___######_______
 * _____######___####_##############__######______
 * ____#######__#####################_#######_____
 * ____#######__##############################____
 * ___#######__######_#################_#######___
 * ___#######__######_######_#########___######___
 * ___#######____##__######___######_____######___
 * ___#######________######____#####_____#####____
 * ____######________#####_____#####_____####_____
 * _____#####________####______#####_____###______
 * ______#####______;###________###______#________
 * ________##_______####________####______________
 */

//获取验证码
func (this *ApiComp) SendEmailCaptchaReq(c *http.Context) {
	req := &SendEmailCaptchaReq{}
	c.ShouldBindJSON(req)
	defer log.Debugf("SendEmailCaptchaReq:%+v", req)
	if sgin := this.module.ParamSign(map[string]interface{}{"Mailbox": req.Mailbox, "CaptchaType": req.CaptchaType}); sgin != req.Sign {
		log.Errorf("LoginByCaptchaReq SignError sgin:%s", sgin)
		this.module.HttpStatusOK(c, core.ErrorCode_SignError, nil)
		return
	}
	if req.Mailbox == "" {
		this.module.HttpStatusOK(c, ErrorCode_ReqParameterError, nil)
		return
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	captchacode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	if err := this.module.captcha.SendEmailCaptcha(req.Mailbox, captchacode); err == nil {
		this.module.captcha.WriteCaptcha(req.Mailbox, captchacode, req.CaptchaType)
		this.module.HttpStatusOK(c, ErrorCode_Success, nil)
	} else {
		this.module.HttpStatusOK(c, ErrorCode_ReqParameterError, nil)
	}
}

/*				 Console相关接口
 * _______________#########_______________________
 * ______________############_____________________
 * ______________#############____________________
 * _____________##__###########___________________
 * ____________###__######_#####__________________
 * ____________###_#######___####_________________
 * ___________###__##########_####________________
 * __________####__###########_####_______________
 * ________#####___###########__#####_____________
 * _______######___###_########___#####___________
 * _______#####___###___########___######_________
 * ______######___###__###########___######_______
 * _____######___####_##############__######______
 * ____#######__#####################_#######_____
 * ____#######__##############################____
 * ___#######__######_#################_#######___
 * ___#######__######_######_#########___######___
 * ___#######____##__######___######_____######___
 * ___#######________######____#####_____#####____
 * ____######________#####_____#####_____####_____
 * _____#####________####______#####_____###______
 * ______#####______;###________###______#________
 * ________##_______####________####______________
 */

//获取集群信息
func (this *ApiComp) GetProjectInfo(c *http.Context) {
	uId := c.GetUInt32(UserKey)
	defer log.Debugf("GetClusterInfo:%d", uId)
	var (
		err   error
		udata *Cache_UserData
	)
	if udata, err = this.module.cache.QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		this.module.HttpStatusOK(c, ErrorCode_Success, this.module.options)
		return
	} else {
		this.module.HttpStatusOK(c, ErrorCode_AuthorityLow, nil)
		return
	}
}

//获取控制台信息
func (this *ApiComp) GetHostInfo(c *http.Context) {
	uId := c.GetUInt32(UserKey)
	defer log.Debugf("GetHostInfo:%d", uId)
	var (
		err   error
		udata *Cache_UserData
	)
	if udata, err = this.module.cache.QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		hostInfo := this.module.hostmonitorcomp.hostInfo
		this.module.HttpStatusOK(c, ErrorCode_Success, hostInfo)
		return
	} else {
		this.module.HttpStatusOK(c, ErrorCode_AuthorityLow, nil)
		return
	}
}

//获取控制台信息
func (this *ApiComp) GetCpuInfo(c *http.Context) {
	uId := c.GetUInt32(UserKey)
	defer log.Debugf("GetHostInfo:%d", uId)
	var (
		err   error
		udata *Cache_UserData
	)
	if udata, err = this.module.cache.QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		cpuinfo := this.module.hostmonitorcomp.cpuInfo
		this.module.HttpStatusOK(c, ErrorCode_Success, cpuinfo)
		return
	} else {
		this.module.HttpStatusOK(c, ErrorCode_AuthorityLow, nil)
		return
	}
}

//获取控制台信息
func (this *ApiComp) GetMemoryInfo(c *http.Context) {
	uId := c.GetUInt32(UserKey)
	defer log.Debugf("GetconSoleinfo:%d", uId)
	var (
		err   error
		udata *Cache_UserData
	)
	if udata, err = this.module.cache.QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		memoryInfo := this.module.hostmonitorcomp.memoryInfo
		this.module.HttpStatusOK(c, ErrorCode_Success, memoryInfo)
		return
	} else {
		this.module.HttpStatusOK(c, ErrorCode_AuthorityLow, nil)
		return
	}
}

//获取主机监控信息
func (this *ApiComp) GetHostMonitorData(c *http.Context) {
	req := &QueryHostMonitorDataReq{}
	c.ShouldBindJSON(req)
	uId := c.GetUInt32(UserKey)
	defer log.Debugf("GetHostMonitorData:%v", req)
	var (
		err   error
		udata *Cache_UserData
	)
	if udata, err = this.module.cache.QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		resp := this.module.hostmonitorcomp.GetHostMonitorData(req.QueryTime)
		this.module.HttpStatusOK(c, ErrorCode_Success, resp)
		return
	} else {
		this.module.HttpStatusOK(c, ErrorCode_AuthorityLow, nil)
		return
	}
}

//获取控制台信息
func (this *ApiComp) GetClusterMonitorData(c *http.Context) {
	req := &QueryClusterMonitorDataReq{}
	c.ShouldBindJSON(req)
	uId := c.GetUInt32(UserKey)
	defer log.Debugf("GetConsoleCluster:%d", uId)
	var (
		err   error
		udata *Cache_UserData
	)
	if udata, err = this.module.cache.QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		resp := this.module.clustermonitorcomp.GetClusterMonitorDataResp(req.QueryTime)
		this.module.HttpStatusOK(c, ErrorCode_Success, resp)
		return
	} else {
		this.module.HttpStatusOK(c, ErrorCode_AuthorityLow, nil)
		return
	}
}

/*				 Console相关接口
 * _______________#########_______________________
 * ______________############_____________________
 * ______________#############____________________
 * _____________##__###########___________________
 * ____________###__######_#####__________________
 * ____________###_#######___####_________________
 * ___________###__##########_####________________
 * __________####__###########_####_______________
 * ________#####___###########__#####_____________
 * _______######___###_########___#####___________
 * _______#####___###___########___######_________
 * ______######___###__###########___######_______
 * _____######___####_##############__######______
 * ____#######__#####################_#######_____
 * ____#######__##############################____
 * ___#######__######_#################_#######___
 * ___#######__######_######_#########___######___
 * ___#######____##__######___######_____######___
 * ___#######________######____#####_____#####____
 * ____######________#####_____#####_____####_____
 * _____#####________####______#####_____###______
 * ______#####______;###________###______#________
 * ________##_______####________####______________
 */

//登录请求 验证码
func (this *ApiComp) LoginByCaptchaReq(c *http.Context) {
	req := &LoginByCaptchaReq{}
	c.ShouldBindJSON(req)
	defer log.Debugf("LoginByCaptchaReq:%+v", req)
	if sgin := this.module.ParamSign(map[string]interface{}{"PhonOrEmail": req.PhonOrEmail, "Captcha": req.Captcha}); sgin != req.Sign {
		log.Errorf("LoginByCaptchaReq SignError sgin:%s", sgin)
		this.module.HttpStatusOK(c, ErrorCode_SignError, nil)
		return
	}
	var (
		captchecode string
		token       string
		err         error
		udata       *DB_UserData
	)
	if captchecode, err = this.module.captcha.QueryCaptcha(req.PhonOrEmail, CaptchaType_LoginCaptcha); err == nil && captchecode == req.Captcha {
		if udata, err = this.module.db.LoginUserDataByPhonOrEmail(&DB_UserData{
			PhonOrEmail: req.PhonOrEmail,
			Password:    this.module.options.UserInitialPassword,
			NickName:    req.PhonOrEmail,
			UserRole:    UserRole_Guester,
		}); err == nil {
			this.module.cache.CleanToken(udata.Token)
			token = this.module.createToken(udata.Id)
			udata.Token = token
			this.module.db.UpDataUserData(udata)
			this.module.cache.WriteToken(token, udata.Id)
			this.module.cache.WriteUserData(&Cache_UserData{
				Db_UserData: udata,
				IsOnLine:    true,
			})
			this.module.HttpStatusOK(c, ErrorCode_Success, token)
			return
		}
	}
	this.module.HttpStatusOK(c, ErrorCode_ReqParameterError, nil)
	return
}

//登录请求 密码
func (this *ApiComp) LoginByPasswordReq(c *http.Context) {
	req := &LoginByPasswordReq{}
	c.ShouldBindJSON(req)
	defer log.Debugf("LoginByPasswordReq:%+v", req)
	if sgin := this.module.ParamSign(map[string]interface{}{"PhonOrEmail": req.PhonOrEmail, "Password": req.Password}); sgin != req.Sign {
		log.Errorf("LoginByPasswordReq SignError sgin:%s", sgin)
		this.module.HttpStatusOK(c, ErrorCode_SignError, nil)
		return
	}
	var (
		err   error
		token string
		udata *DB_UserData
	)
	if udata, err = this.module.db.QueryUserDataByPhonOrEmail(req.PhonOrEmail); err == nil {
		if udata.Password == req.Password {
			this.module.cache.CleanToken(udata.Token)
			token = this.module.createToken(udata.Id)
			udata.Token = token
			this.module.db.UpDataUserData(udata)
			this.module.cache.WriteToken(token, udata.Id)
			this.module.cache.WriteUserData(&Cache_UserData{
				Db_UserData: udata,
				IsOnLine:    true,
			})
			this.module.HttpStatusOK(c, ErrorCode_Success, token)
			return
		}
	}
	this.module.HttpStatusOK(c, ErrorCode_ReqParameterError, nil)
	return
}

//登录请求 Token
func (this *ApiComp) GetUserinfoReq(c *http.Context) {
	uId := c.GetUInt32(UserKey)
	defer log.Debugf("GetUserinfoReq:%d", uId)
	var (
		err   error
		udata *DB_UserData
	)
	if udata, err = this.module.db.QueryUserDataById(uId); err == nil {
		this.module.cache.WriteUserData(&Cache_UserData{
			Db_UserData: udata,
			IsOnLine:    true,
		})
		this.module.HttpStatusOK(c, ErrorCode_Success, udata)
		return
	}
	this.module.HttpStatusOK(c, ErrorCode_TokenExpired, nil)
	return
}

//登出
func (this *ApiComp) LoginOutReq(c *http.Context) {
	req := &LoginOutReq{}
	c.ShouldBindJSON(req)
	uId := c.GetUInt32(UserKey)
	defer log.Debugf("LoginOutReq:%+v UserId:%d", req, uId)
	var (
		result *DB_UserData
		err    error
	)
	if result, err = this.module.db.QueryUserDataById(uId); err == nil {
		this.module.cache.CleanUserData(uId)
		this.module.cache.CleanToken(result.Token)
		this.module.HttpStatusOK(c, ErrorCode_Success, nil)
		return
	}
	log.Errorf("LoginOutReq err:%v", err)
	this.module.HttpStatusOK(c, ErrorCode_SqlExecutionError, nil)
	return
}
