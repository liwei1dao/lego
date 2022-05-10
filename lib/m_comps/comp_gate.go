package m_comps

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib"
	"github.com/liwei1dao/lego/lib/modules/gate"
	"github.com/liwei1dao/lego/lib/s_comps"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/proto"
	"github.com/liwei1dao/lego/sys/workerpools"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type MComp_GateCompOptions struct {
	DefGoroutine  int //网关工作池默认协程数
	MaxGoroutine  int //网关工作池最大协程数
	HandleTimeOut int //工作池处理超时时间
}

func (this *MComp_GateCompOptions) LoadConfig(settings map[string]interface{}) (err error) {
	this.DefGoroutine = 100
	this.MaxGoroutine = 1000
	this.HandleTimeOut = 5
	if settings != nil {
		err = mapstructure.Decode(settings, this)
	}
	return
}

func (this *MComp_GateCompOptions) GetDefGoroutine() int {
	return this.DefGoroutine
}
func (this *MComp_GateCompOptions) GetGateMaxGoroutine() int {
	return this.MaxGoroutine
}
func (this *MComp_GateCompOptions) GetHandleTimeOut() int {
	return this.HandleTimeOut
}

/*
模块 网关组件
接待网关 代理 的业务需求
*/
type MComp_GateComp struct {
	cbase.ModuleCompBase
	service    core.IService
	comp       IMComp_GateComp
	options    IMComp_GateCompOptions
	ComId      uint16 //协议分类Id
	IsLog      bool   //是否输出消息日志
	Msghandles map[uint16]*msgRecep
	Workerpool workerpools.ISys
}

type msgRecep struct {
	msgId   uint16
	MsgType reflect.Type
	F       func(ctx context.Context, session core.IUserSession, msg interface{})
}

func (this *MComp_GateComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	this.ModuleCompBase.Init(service, module, comp, options)
	this.service = service
	this.comp = comp.(IMComp_GateComp)
	this.options = options.(IMComp_GateCompOptions)
	this.Msghandles = make(map[uint16]*msgRecep)
	this.Workerpool, err = workerpools.NewSys(workerpools.SetDefWorkers(this.options.GetDefGoroutine()), workerpools.SetMaxWorkers(this.options.GetGateMaxGoroutine()), workerpools.SetTaskTimeOut(time.Second*time.Duration(this.options.GetHandleTimeOut())))
	return
}

func (this *MComp_GateComp) Start() (err error) {
	if err = this.ModuleCompBase.Start(); err != nil {
		return
	}
	isRegisterLocalRoute := false
	//注册本地路由
	m, e := this.service.GetModule(lib.SM_GateModule)
	if e == nil {
		m.(gate.IGateModule).RegisterLocalRoute(this.ComId, this.comp.ReceiveMsg)
		isRegisterLocalRoute = true
	}
	if !isRegisterLocalRoute {
		//注册远程路由
		cc, e := this.service.GetComp(lib.SC_ServiceGateRouteComp)
		if e == nil {
			cc.(s_comps.ISC_GateRouteComp).RegisterRoute(this.ComId, this.comp.ReceiveMsg)
			isRegisterLocalRoute = true
		}
	}

	if !isRegisterLocalRoute {
		return fmt.Errorf("MC_GateComp 未成功注册路由!")
	}
	return
}

func (this *MComp_GateComp) ReceiveMsg(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string) {
	this.Workerpool.Submit(this.handlemsg, session, msg)
	return core.ErrorCode_Success, ""
}

func (this *MComp_GateComp) RegisterHandle(mId uint16, msg interface{}, f func(ctx context.Context, session core.IUserSession, msg interface{})) {
	if _, ok := this.Msghandles[mId]; ok {
		log.Errorf("重复 注册网关【%d】消息【%d】", this.ComId, mId)
		return
	}
	this.Msghandles[mId] = &msgRecep{
		msgId:   mId,
		MsgType: reflect.TypeOf(msg),
		F:       f,
	}
}

func (this *MComp_GateComp) handlemsg(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}) {
	defer cancel()
	session, msg := agrs[0].(core.IUserSession), agrs[1].(proto.IMessage)    //任务结束通知上层
	defer lego.Recover(fmt.Sprintf("MComp_GateComp ReceiveMsg msg:%v", msg)) //打印消息处理异常信息
	msghandles, ok := this.Msghandles[msg.GetMsgId()]
	if !ok {
		log.Errorf("模块网关路由【%d】没有注册消息【%d】接口", this.ComId, msg.GetMsgId())
		return
	}
	msgdata, e := proto.ByteDecodeToStruct(msghandles.MsgType, msg.GetBuffer())
	if e != nil {
		log.Errorf("收到异常消息【%d:%d】来自【%s】的消息:%v err:%v", this.ComId, msg.GetMsgId(), session.GetSessionId(), msg.GetBuffer(), e)
		session.Close()
		return
	}
	if this.IsLog {
		log.Infof("收到【%d:%d】来自【%s】的消息:%v", this.ComId, msg.GetMsgId(), session.GetSessionId(), msgdata)
	}
	msghandles.F(ctx, session, msgdata)
}
