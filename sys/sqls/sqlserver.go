package sqls

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

func newSqlServer(opt ...Option) (*SqlServer, error) {
	sqlserver := new(SqlServer)
	sqlserver.opts = newOptions(opt...)
	err := sqlserver.init()
	return sqlserver, err
}

type SqlServer struct {
	opts Options
	db   *sql.DB
}

func (this *SqlServer) init() (err error) {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable",
		this.opts.ServerAdd, this.opts.User, this.opts.Password, this.opts.Port, this.opts.Database)
	// Create connection pool
	this.db, err = sql.Open("mssql", connString)
	return
}

func (this *SqlServer) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.opts.TimeOut)
	return
}

func (this *SqlServer) ExecContext(query string, args ...interface{}) (err error) {
	_, err = this.db.ExecContext(this.getContext(), query, args...)
	return
}

func (this *SqlServer) QueryContext(query string, args ...interface{}) (data *sql.Rows, err error) {
	data, err = this.db.QueryContext(this.getContext(), query, args...)
	return
}
