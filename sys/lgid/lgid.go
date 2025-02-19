package lgid

import (
	"github.com/bwmarrin/snowflake"
	"github.com/rs/xid"
)

func newSys(options Options) (sys *LgID, err error) {
	sys = &LgID{options: options}
	sys.number, err = snowflake.NewNode(1)
	return
}

type LgID struct {
	options Options
	number  *snowflake.Node
}

// 获取一个新的数值唯一ID
func (this *LgID) NewNumberID() int64 {
	id := this.number.Generate()
	return id.Int64()
}

// 获取一个新的数值唯一ID
func (this *LgID) NewStringID() string {
	id := xid.New()
	return id.String()
}

func (this *LgID) NewStringIDForNumber() string {
	id := this.number.Generate()
	return id.String()
}
