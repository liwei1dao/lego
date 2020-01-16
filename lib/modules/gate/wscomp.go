package gate

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
	"net"
	"net/http"
	"sync"
	"time"
)

type WsServerComp struct {
	cbase.ModuleCompBase
	module     IGateModule
	wsaddr     string
	CertFile   string
	KeyFile    string
	outTime    time.Duration
	ln         net.Listener
	handler    *WSHandler
	NewWsAgent func(gate IGateModule, coon IConn) (IAgent, error)
}

type WSHandler struct {
	maxConnNum      int
	pendingWriteNum int
	maxMsgLen       uint32
	upgrader        websocket.Upgrader
	mutexConns      sync.Mutex
	wg              sync.WaitGroup
	myComp          *WsServerComp
}

func (this *WsServerComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, settings map[string]interface{}) (err error) {
	if m, ok := module.(IGateModule); !ok {
		return fmt.Errorf("TcpServerComp Init module is no IGateModule")
	} else {
		this.module = m
	}
	if WSAddr, ok := settings["WsAddr"]; ok {
		this.wsaddr = WSAddr.(string)
	} else {
		err = fmt.Errorf("启动WsServiceComp 组件失败 配置错误 WsAddr")
		return
	}
	if HTTPTimeout, ok := settings["HttpTimeout"]; ok {
		this.outTime = time.Second * time.Duration(HTTPTimeout.(int64))
	} else {
		err = fmt.Errorf("启动WsServiceComp 组件失败 HttpTimeout 配置错误")
		return
	}
	if this.NewWsAgent == nil {
		err = fmt.Errorf("启动WsServiceComp 组件失败 代理接口没有实现")
		return
	}
	err = this.ModuleCompBase.Init(service, module, comp, settings)
	return
}

func (this *WsServerComp) Start() (err error) {
	err = this.ModuleCompBase.Start()
	ln, err := net.Listen("tcp", this.wsaddr)
	if err != nil {
		err = fmt.Errorf("WsServerComp Listen 失败 %v", err)
		return
	}
	if this.CertFile != "" || this.KeyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}

		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(this.CertFile, this.KeyFile)
		if err != nil {
			log.Errorf("%v", err)
		}
		ln = tls.NewListener(ln, config)
	}
	this.ln = ln
	this.handler = &WSHandler{
		upgrader: websocket.Upgrader{
			HandshakeTimeout: this.outTime,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
		myComp: this,
	}
	httpServer := &http.Server{
		Addr:           this.wsaddr,
		Handler:        this.handler,
		ReadTimeout:    this.outTime,
		WriteTimeout:   this.outTime,
		MaxHeaderBytes: 1024,
	}
	log.Infof("WsServerComp Listen %s 成功", this.wsaddr)
	go httpServer.Serve(ln)
	return
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade error: %v", err)
		return
	}
	handler.wg.Add(1)
	defer handler.wg.Done()
	wc := NewWsConn(conn)
	agent, err := handler.myComp.NewWsAgent(handler.myComp.module, wc)
	if err == nil {
		agent.OnRun()
		agent.Destory()
	} else {
		log.Errorf("WsServerComp NewWsAgent err %s", err.Error())
		wc.Close()
	}
}

func (this *WsServerComp) Destroy() (err error) {
	err = this.ModuleCompBase.Destroy()
	this.ln.Close()
	this.handler.wg.Wait()
	return
}
