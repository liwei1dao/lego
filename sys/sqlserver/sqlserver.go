package sqlserver

import (
	"context"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb" //需要提前导入否赠会连接失败
	"github.com/jmoiron/sqlx"
)

func newSys(options Options) (sys *SqlServer, err error) {
	sys = &SqlServer{options: options}
	sys.db, err = sqlx.Connect("sqlserver", options.SqlUrl)

	return
}

type SqlServer struct {
	options Options
	db      *sqlx.DB
}

func (this *SqlServer) Close() error {
	return this.db.Close()
}

func (this *SqlServer) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.options.TimeOut)
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
