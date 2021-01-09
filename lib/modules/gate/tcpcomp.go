package gate

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
)

type TcpServerComp struct {
	cbase.ModuleCompBase
	module      IGateModule
	listener    net.Listener
	wgLn        sync.WaitGroup
	wgConns     sync.WaitGroup
	NewTcpAgent func(gate IGateModule, coon IConn) (IAgent, error)
}

func (this *TcpServerComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	this.module = module.(IGateModule)
	if this.NewTcpAgent == nil {
		err = fmt.Errorf("启动Tcp 组件失败 代理接口没有实现")
		return
	}
	err = this.ModuleCompBase.Init(service, module, comp, options)
	return
}

func (this *TcpServerComp) Start() (err error) {
	err = this.ModuleCompBase.Start()
	ln, err := net.Listen("tcp", this.module.GetOptions().GetTcpAddr())
	if err != nil {
		err = fmt.Errorf("TcpServerComp Listen:%s 失败 %s", this.module.GetOptions().GetTcpAddr(), err.Error())
		return
	} else {
		log.Infof("TcpServerComp Listen:%s 成功", this.module.GetOptions().GetTcpAddr())
	}
	this.listener = ln
	go this.run()
	return
}

func (this *TcpServerComp) run() {
	var tempDelay time.Duration
locp:
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Infof("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			break locp
		}
		tempDelay = 0
		tc := NewTcpConn(conn)
		agent, err := this.NewTcpAgent(this.module, tc)
		if err == nil {
			this.wgConns.Add(1)
			go func() {
				agent.OnRun()
				agent.Destory()
				this.wgConns.Done()
			}()
		} else {
			tc.Close()
		}
	}
	log.Infof("Tcp 监听退出")
}

func (this *TcpServerComp) Destroy() (err error) {
	err = this.ModuleCompBase.Destroy()
	this.listener.Close()
	this.wgConns.Wait()
	return
}
