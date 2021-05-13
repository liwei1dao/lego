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
	module IConsole
}

func (this *ApiComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.module = module.(IConsole)
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
	user.POST("/registerbycaptcha", this.RegisterByCaptchaReq)
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
	if sgin := this.module.ParamSign(map[string]interface{}{"PhonOrEmail": req.PhonOrEmail, "CaptchaType": req.CaptchaType}); sgin != req.Sign {
		log.Errorf("LoginByCaptchaReq SignError sgin:%s", sgin)
		this.module.HttpStatusOK(c, core.ErrorCode_SignError, nil)
		return
	}
	if req.PhonOrEmail == "" {
		this.module.HttpStatusOK(c, ErrorCode_ReqParameterError, nil)
		return
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	captchacode := fmt.Sprintf("%04v", rnd.Int31n(10000))
	if err := this.module.Captcha().SendEmailCaptcha(req.PhonOrEmail, captchacode); err == nil {
		this.module.Captcha().WriteCaptcha(req.PhonOrEmail, captchacode, req.CaptchaType)
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
	if udata, err = this.module.Cache().QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		this.module.HttpStatusOK(c, ErrorCode_Success, this.module.Options())
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
	if udata, err = this.module.Cache().QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		hostInfo := this.module.Hostmonitorcomp().HostInfo()
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
	if udata, err = this.module.Cache().QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		cpuinfo := this.module.Hostmonitorcomp().CpuInfo()
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
	if udata, err = this.module.Cache().QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		memoryInfo := this.module.Hostmonitorcomp().MemoryInfo()
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
	if udata, err = this.module.Cache().QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		resp := this.module.Hostmonitorcomp().GetHostMonitorData(req.QueryTime)
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
	if udata, err = this.module.Cache().QueryUserData(uId); err == nil && udata.Db_UserData.UserRole >= UserRole_Developer {
		resp := this.module.Clustermonitorcomp().GetClusterMonitorDataResp(req.QueryTime)
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

//注册请求 验证码
func (this *ApiComp) RegisterByCaptchaReq(c *http.Context) {
	req := &RegisterByCaptchaReq{}
	c.ShouldBindJSON(req)
	defer log.Debugf("RegisterByCaptchaReq:%+v", req)
	if sgin := this.module.ParamSign(map[string]interface{}{"PhonOrEmail": req.PhonOrEmail, "Password": req.Password, "Captcha": req.Captcha}); sgin != req.Sign {
		log.Errorf("RegisterByCaptchaReq SignError sgin:%s", sgin)
		this.module.HttpStatusOK(c, ErrorCode_SignError, nil)
		return
	}
	var (
		captchecode string
		token       string
		err         error
		udata       *DB_UserData
	)
	if captchecode, err = this.module.Captcha().QueryCaptcha(req.PhonOrEmail, CaptchaType_RegisterCaptcha); err == nil && captchecode == req.Captcha {
		if udata, err = this.module.DB().LoginUserDataByPhonOrEmail(&DB_UserData{
			PhonOrEmail: req.PhonOrEmail,
			Password:    req.Password,
			NickName:    req.PhonOrEmail,
			UserRole:    UserRole_Guester,
		}); err == nil {
			this.module.Cache().CleanToken(udata.Token)
			token = this.module.CreateToken(udata.Id)
			udata.Token = token
			this.module.DB().UpDataUserData(udata)
			this.module.Cache().WriteToken(token, udata.Id)
			this.module.Cache().WriteUserData(&Cache_UserData{
				Db_UserData: udata,
				IsOnLine:    true,
			})
			this.module.HttpStatusOK(c, ErrorCode_Success, token)
			return
		}
	}
	this.module.HttpStatusOK(c, ErrorCode_ReqParameterError, nil)
}

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
	if captchecode, err = this.module.Captcha().QueryCaptcha(req.PhonOrEmail, CaptchaType_LoginCaptcha); err == nil && captchecode == req.Captcha {
		if udata, err = this.module.DB().LoginUserDataByPhonOrEmail(&DB_UserData{
			PhonOrEmail: req.PhonOrEmail,
			Password:    this.module.Options().GetUserInitialPassword(),
			NickName:    req.PhonOrEmail,
			UserRole:    UserRole_Guester,
		}); err == nil {
			this.module.Cache().CleanToken(udata.Token)
			token = this.module.CreateToken(udata.Id)
			udata.Token = token
			this.module.DB().UpDataUserData(udata)
			this.module.Cache().WriteToken(token, udata.Id)
			this.module.Cache().WriteUserData(&Cache_UserData{
				Db_UserData: udata,
				IsOnLine:    true,
			})
			this.module.HttpStatusOK(c, ErrorCode_Success, token)
			return
		}
	}
	this.module.HttpStatusOK(c, ErrorCode_ReqParameterError, nil)
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
	if udata, err = this.module.DB().QueryUserDataByPhonOrEmail(req.PhonOrEmail); err == nil {
		if udata.Password == req.Password {
			this.module.Cache().CleanToken(udata.Token)
			token = this.module.CreateToken(udata.Id)
			udata.Token = token
			this.module.DB().UpDataUserData(udata)
			this.module.Cache().WriteToken(token, udata.Id)
			this.module.Cache().WriteUserData(&Cache_UserData{
				Db_UserData: udata,
				IsOnLine:    true,
			})
			this.module.HttpStatusOK(c, ErrorCode_Success, udata)
			return
		}
	}
	this.module.HttpStatusOK(c, ErrorCode_ReqParameterError, nil)
}

//登录请求 Token
func (this *ApiComp) GetUserinfoReq(c *http.Context) {
	uId := c.GetUInt32(UserKey)
	defer log.Debugf("GetUserinfoReq:%d", uId)
	var (
		err   error
		udata *DB_UserData
	)
	if udata, err = this.module.DB().QueryUserDataById(uId); err == nil {
		this.module.Cache().WriteUserData(&Cache_UserData{
			Db_UserData: udata,
			IsOnLine:    true,
		})
		this.module.HttpStatusOK(c, ErrorCode_Success, udata)
		return
	}
	this.module.HttpStatusOK(c, ErrorCode_TokenExpired, nil)
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
	if result, err = this.module.DB().QueryUserDataById(uId); err == nil {
		this.module.Cache().CleanUserData(uId)
		this.module.Cache().CleanToken(result.Token)
		this.module.HttpStatusOK(c, ErrorCode_Success, nil)
		return
	}
	log.Errorf("LoginOutReq err:%v", err)
	this.module.HttpStatusOK(c, ErrorCode_SqlExecutionError, nil)
}
