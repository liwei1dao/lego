package hls

import (
	"net"
	"sync"

	"github.com/liwei1dao/lego/sys/livego/core"
)

type Server struct {
	sys      core.ISys
	listener net.Listener
	conns    *sync.Map
}
