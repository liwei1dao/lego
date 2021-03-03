package sqls

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type (
	ISys interface {
		Close() error
		ExecContext(query string, args ...interface{}) (data sql.Result, err error)
		QueryContext(query string, args ...interface{}) (data *sql.Rows, err error)
		QueryxContext(query string, args ...interface{}) (data *sqlx.Rows, err error)
		QueryRowxContext(query string, args ...interface{}) (data *sqlx.Row)
		GetContext(dest interface{}, query string, args ...interface{}) error
		SelectContext(dest interface{}, query string, args ...interface{}) error
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func Close() error {
	return defsys.Close()
}

func ExecContext(query string, args ...interface{}) (data sql.Result, err error) {
	return defsys.ExecContext(query, args...)
}

func QueryContext(query string, args ...interface{}) (data *sql.Rows, err error) {
	return defsys.QueryContext(query, args...)
}

func QueryxContext(query string, args ...interface{}) (data *sqlx.Rows, err error) {
	data, err = defsys.QueryxContext(query, args...)
	return
}

func QueryRowxContext(query string, args ...interface{}) (data *sqlx.Row) {
	return defsys.QueryRowxContext(query, args...)
}

func GetContext(dest interface{}, query string, args ...interface{}) error {
	return defsys.GetContext(dest, query, args...)
}

func SelectContext(dest interface{}, query string, args ...interface{}) error {
	return defsys.SelectContext(dest, query, args...)
}
