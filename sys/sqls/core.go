package sqls

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"

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

func ExecContext(query string, args ...interface{}) (err error) {
	return sqlserver.ExecContext(query, args...)
}

func QueryContext(query string, args ...interface{}) (data *sql.Rows, err error) {
	return sqlserver.QueryContext(query, args...)
}
