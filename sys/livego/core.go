package livego

import "github.com/liwei1dao/lego/sys/livego/core"

var (
	EmptyID = ""
)

type (
	GetInFo interface {
		GetInfo() (string, string, string)
	}
	StreamReadWriteCloser interface {
		GetInFo
		Close(error)
		Write(core.ChunkStream) error
		Read(c *core.ChunkStream) error
	}
)
