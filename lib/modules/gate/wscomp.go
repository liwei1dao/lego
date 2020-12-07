package gate

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
)

type WsServerComp struct {
	cbase.ModuleCompBase
	module IGateModule
	// CertFile   string
	// KeyFile    string
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

func (this *WsServerComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	this.module = module.(IGateModule)
	if this.NewWsAgent == nil {
		err = fmt.Errorf("启动WsServiceComp 组件失败 代理接口没有实现")
		return
	}
	err = this.ModuleCompBase.Init(service, module, comp, options)
	return
}

func (this *WsServerComp) Start() (err error) {
	err = this.ModuleCompBase.Start()
	ln, err := net.Listen("tcp", this.module.GetOptions().GetWSAddr())
	if err != nil {
		err = fmt.Errorf("WsServerComp Listen Fail %v", err)
		return
	}
	if this.module.GetOptions().GetWSCertFile() != "" || this.module.GetOptions().GetWSKeyFile() != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}

		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(this.module.GetOptions().GetWSCertFile(), this.module.GetOptions().GetWSKeyFile())
		if err != nil {
			log.Errorf("%v", err)
		}
		ln = tls.NewListener(ln, config)
	}
	this.ln = ln
	this.handler = &WSHandler{
		upgrader: websocket.Upgrader{
			HandshakeTimeout: this.module.GetOptions().GetWSOuttime(),
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
		myComp: this,
	}
	httpServer := &http.Server{
		Addr:           this.module.GetOptions().GetWSAddr(),
		Handler:        this.handler,
		ReadTimeout:    this.module.GetOptions().GetWSOuttime(),
		WriteTimeout:   this.module.GetOptions().GetWSOuttime(),
		MaxHeaderBytes: 1024,
	}
	log.Infof("WsServerComp Listen %s Success", this.module.GetOptions().GetWSAddr())
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
	wc := NewWsConn(conn, handler.myComp.module.GetOptions().GetWSHeartbeat())
	agent, err := handler.myComp.NewWsAgent(handler.myComp.module, wc)
	if err == nil {
		agent.OnRun()
		agent.Destory()
	} else {
		log.Errorf("WsServerComp NewWsAgent err %v", err)
		wc.Close()
	}
}

func (this *WsServerComp) Destroy() (err error) {
	err = this.ModuleCompBase.Destroy()
	this.ln.Close()
	this.handler.wg.Wait()
	return
}
