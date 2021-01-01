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
)

/*
模块 网关组件
接待网关 代理 的业务需求
*/
type MComp_GateComp struct {
	cbase.ModuleCompBase
	service      core.IService
	comp         IMComp_GateComp
	ComId        uint16 //协议分类Id
	IsLog        bool   //是否输出消息日志
	Msghandles   map[uint16]*msgRecep
	Mrlock       sync.RWMutex
	MaxGoroutine int //最大并发数据
	Workerpool   workerpools.IWorkerPool
}

type msgRecep struct {
	msgId   uint16
	MsgType reflect.Type
	F       func(session core.IUserSession, msg interface{})
}

func (this *MComp_GateComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, setting map[string]interface{}) (err error) {
	this.ModuleCompBase.Init(service, module, comp, setting)
	this.service = service
	this.comp = comp.(IMComp_GateComp)
	if v, ok := setting["GateMaxGoroutine"]; ok {
		this.MaxGoroutine = v.(int)
	} else {
		log.Warnf("Module:%s Lack Config:GateMaxGoroutine", module.GetType())
		this.MaxGoroutine = 1000
	}
	this.Msghandles = make(map[uint16]*msgRecep)
	this.Workerpool, err = workerpools.NewTaskPools(workerpools.SetMaxWorkers(this.MaxGoroutine), workerpools.SetTaskTimeOut(time.Second*2))
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
		_gatecomp, _session, _msg := agrs[0].(*MComp_GateComp), agrs[1].(core.IUserSession), agrs[2].(proto.IMessage)

		_gatecomp.Mrlock.RLock()
		msghandles, ok := _gatecomp.Msghandles[msg.GetMsgId()]
		_gatecomp.Mrlock.RUnlock()
		if !ok {
			log.Errorf("模块网关路由【%d】没有注册消息【%d】接口", _gatecomp.ComId, _msg.GetMsgId())
			return
		}
		msgdata, e := proto.MsgUnMarshal(msghandles.MsgType, _msg.GetMsg())
		if e != nil {
			log.Errorf("收到异常消息【%d:%d】来自【%s】的消息:%v err:%v", _gatecomp.ComId, _msg.GetMsgId(), _session.GetSessionId(), _msg.GetMsg(), e)
			_session.Close()
			return
		}
		if _gatecomp.IsLog {
			log.Infof("收到【%d:%d】来自【%s】的消息:%v", _gatecomp.ComId, _msg.GetMsgId(), _session.GetSessionId(), msgdata)
		}
		msghandles.F(_session, msgdata)
	}, this, session, msg)
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
