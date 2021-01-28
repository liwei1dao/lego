package gate

import (
	"fmt"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/proto"
)

type SessionData struct {
	IP           string
	SessionId    string
	GateServerId string //用户所在网关服务
}

func NewRemoteSession(service base.IClusterService, data map[string]interface{}) (s core.IUserSession, err error) {
	session := &RemoteSession{
		service: service,
		data:    &SessionData{},
	}
	err = session.UpDateMap(data)
	if err == nil {
		s = session
	}
	return s, err
}

type RemoteSession struct {
	service base.IClusterService
	data    *SessionData
}

func (this *RemoteSession) UpDateMap(data map[string]interface{}) (err error) {
	IP := data["IP"]
	if IP != nil {
		this.data.IP = IP.(string)
	} else {
		return fmt.Errorf("用户会话 关键数据 IP 不存在")
	}
	Sessionid := data["SessionId"]
	if Sessionid != nil {
		this.data.SessionId = Sessionid.(string)
	} else {
		return fmt.Errorf("用户会话 关键数据 SessionId 不存在")
	}
	Serverid := data["GateServerId"]
	if Serverid != nil {
		this.data.GateServerId = Serverid.(string)
	} else {
		return fmt.Errorf("用户会话 关键数据 GateServerId 不存在")
	}
	return nil
}
func (this *RemoteSession) GetGateServerId() string {
	return this.data.GateServerId
}
func (this *RemoteSession) GetSessionId() string {
	return this.data.SessionId
}
func (this *RemoteSession) GetIP() string {
	return this.data.IP
}
func (this *RemoteSession) GetGateId() string {
	return this.data.GateServerId
}

func (this *RemoteSession) SendMsg(comdId uint16, msgId uint16, msg interface{}) (err error) {
	m := proto.EncodeToMesage(comdId, msgId, msg)
	_, err = this.service.RpcInvokeById(this.data.GateServerId, RPC_GateSendMsg, false, this.GetSessionId(), m)
	return err
}
func (this *RemoteSession) Close() (err error) {
	_, err = this.service.RpcInvokeById(this.data.GateServerId, RPC_GateAgentClose, false, this.GetSessionId())
	return
}

func NewLocalSession(module IGateModule, data map[string]interface{}) (s core.IUserSession, err error) {
	session := &LocalSession{
		module: module,
		data:   &SessionData{},
	}
	err = session.UpDateMap(data)
	if err == nil {
		s = session
	}
	return s, err
}

type LocalSession struct {
	module IGateModule
	data   *SessionData
}

func (this *LocalSession) UpDateMap(data map[string]interface{}) (err error) {
	IP := data["IP"]
	if IP != nil {
		this.data.IP = IP.(string)
	} else {
		return fmt.Errorf("用户会话 关键数据 IP 不存在")
	}
	Sessionid := data["SessionId"]
	if Sessionid != nil {
		this.data.SessionId = Sessionid.(string)
	} else {
		return fmt.Errorf("用户会话 关键数据 SessionId 不存在")
	}
	Serverid := data["GateServerId"]
	if Serverid != nil {
		this.data.GateServerId = Serverid.(string)
	} else {
		return fmt.Errorf("用户会话 关键数据 GateServerId 不存在")
	}
	return nil
}
func (this *LocalSession) GetGateServerId() string {
	return this.data.GateServerId
}
func (this *LocalSession) GetSessionId() string {
	return this.data.SessionId
}
func (this *LocalSession) GetIP() string {
	return this.data.IP
}
func (this *LocalSession) GetGateId() string {
	return this.data.GateServerId
}
func (this *LocalSession) SendMsg(comdId uint16, msgId uint16, msg interface{}) (err error) {
	m := proto.EncodeToMesage(comdId, msgId, msg)
	if _, e := this.module.SendMsg(this.GetSessionId(), m); e != "" {
		err = fmt.Errorf(e)
	}
	return
}
func (this *LocalSession) Close() (err error) {
	if _, e := this.module.CloseAgent(this.GetSessionId()); e != "" {
		err = fmt.Errorf(e)
	}
	return
}
