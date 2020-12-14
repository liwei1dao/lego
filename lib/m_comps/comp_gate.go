package m_comps

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

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
	MaxGoroutine int
}

func (this *MComp_GateCompOptions) LoadConfig(settings map[string]interface{}) (err error) {
	this.MaxGoroutine = 100
	if settings != nil {
		err = mapstructure.Decode(settings, this)
	}
	return
}

func (this *MComp_GateCompOptions) GetGateMaxGoroutine() int {
	return this.MaxGoroutine
}

/*
模块 网关组件
接待网关 代理 的业务需求
*/
type MComp_GateComp struct {
	cbase.ModuleCompBase
	service    core.IService
	comp       IMComp_GateComp
	ComId      uint16 //协议分类Id
	IsLog      bool   //是否输出消息日志
	Msghandles map[uint16]*msgRecep
	Mrlock     sync.RWMutex
	Workerpool workerpools.IWorkerPool
}

type msgRecep struct {
	msgId   uint16
	MsgType reflect.Type
	F       func(session core.IUserSession, msg interface{})
}

func (this *MComp_GateComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	this.ModuleCompBase.Init(service, module, comp, options)
	this.service = service
	this.comp = comp.(IMComp_GateComp)
	this.Msghandles = make(map[uint16]*msgRecep)
	this.Workerpool, err = workerpools.NewSys(workerpools.SetMaxWorkers(options.(IMComp_GateCompOptions).GetGateMaxGoroutine()), workerpools.SetTaskTimeOut(time.Second*2))
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
	this.Workerpool.Submit(func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}) {
		defer cancel()        //任务结束通知上层
		defer cbase.Recover() //打印消息处理异常信息

		this.Mrlock.RLock()
		msghandles, ok := this.Msghandles[msg.GetMsgId()]
		this.Mrlock.RUnlock()
		if !ok {
			log.Errorf("模块网关路由【%d】没有注册消息【%d】接口", this.ComId, msg.GetMsgId())
			return
		}
		msgdata, e := proto.ByteDecodeToStruct(msghandles.MsgType, msg.GetBuffer())
		if e != nil {
			log.Errorf("收到异常消息【%d:%d】来自【%s】的消息err:%v", this.ComId, msg.GetMsgId(), session.GetSessionId(), e)
			session.Close()
			return
		}
		if this.IsLog {
			log.Infof("收到【%d:%d】来自【%s】的消息:%v", this.ComId, msg.GetMsgId(), session.GetSessionId(), msgdata)
		}
		msghandles.F(session, msgdata)
	})
	return core.ErrorCode_Success, ""
}

func (this *MComp_GateComp) RegisterHandle(mId uint16, msg interface{}, f func(session core.IUserSession, msg interface{})) {
	if _, ok := this.Msghandles[mId]; ok {
		log.Errorf("重复 注册网关【%d】消息【%d】", this.ComId, mId)
		return
	}
	this.Mrlock.Lock()
	this.Msghandles[mId] = &msgRecep{
		msgId:   mId,
		MsgType: reflect.TypeOf(msg),
		F:       f,
	}
	this.Mrlock.Unlock()
}
