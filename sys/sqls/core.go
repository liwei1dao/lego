package sqls

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
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

func Close() {
	sqlserver.Close()
}

func ExecContext(query string, args ...interface{}) (data sql.Result, err error) {
	return sqlserver.ExecContext(query, args...)
}

func QueryContext(query string, args ...interface{}) (data *sql.Rows, err error) {
	return sqlserver.QueryContext(query, args...)
}

func QueryxContext(query string, args ...interface{}) (data *sqlx.Rows, err error) {
	data, err = sqlserver.QueryxContext(query, args...)
	return
}

func QueryRowxContext(query string, args ...interface{}) (data *sqlx.Row) {
	return sqlserver.QueryRowxContext(query, args...)
}

func GetContext(dest interface{}, query string, args ...interface{}) error {
	return sqlserver.GetContext(dest, query, args...)
}

func SelectContext(dest interface{}, query string, args ...interface{}) error {
	return sqlserver.SelectContext(dest, query, args...)
}
