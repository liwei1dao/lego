package sqls

import (
	"github.com/liwei1dao/lego/core"
)

var (
	service   core.IService
	sqlserver *SqlServer
)

func OnInit(s core.IService, opt ...Option) (err error) {
	service = s
	sqlserver, err = newSqlServer(opt...)
	return
}
