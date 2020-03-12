package sqls

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func newSqlServer(opt ...Option) (*SqlServer, error) {
	sqlserver := new(SqlServer)
	sqlserver.opts = newOptions(opt...)
	err := sqlserver.init()
	return sqlserver, err
}

type SqlServer struct {
	opts Options
	db   *sqlx.DB
}

func (this *SqlServer) init() (err error) {
	this.db, err = sqlx.Connect("sqlserver", this.opts.SqlUrl)
	if err != nil {
		fmt.Printf("err=%s\n", err)
	}
	return
}

func (this *SqlServer) Close() {
	this.db.Close()
}

func (this *SqlServer) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.opts.TimeOut)
	return
}

func (this *SqlServer) QueryContext(query string, args ...interface{}) (data *sql.Rows, err error) {
	data, err = this.db.QueryContext(this.getContext(), query, args...)
	return
}

func (this *SqlServer) ExecContext(query string, args ...interface{}) (data sql.Result, err error) {
	data, err = this.db.ExecContext(this.getContext(), query, args...)
	return
}

func (this *SqlServer) QueryxContext(query string, args ...interface{}) (data *sqlx.Rows, err error) {
	data, err = this.db.QueryxContext(this.getContext(), query, args...)
	return
}

func (this *SqlServer) QueryRowxContext(query string, args ...interface{}) (data *sqlx.Row) {
	return this.db.QueryRowxContext(this.getContext(), query, args...)
}

func (this *SqlServer) GetContext(dest interface{}, query string, args ...interface{}) error {
	return this.db.GetContext(this.getContext(), dest, query, args...)
}

func (this *SqlServer) SelectContext(dest interface{}, query string, args ...interface{}) error {
	return this.db.SelectContext(this.getContext(), dest, query, args...)
}
