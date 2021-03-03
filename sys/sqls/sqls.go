package sqls

import (
	"context"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb" //需要提前导入否赠会连接失败
	"github.com/jmoiron/sqlx"
)

func newSys(options Options) (sys *Sqls, err error) {
	sys = &Sqls{options: options}
	sys.db, err = sqlx.Connect("sqlserver", options.SqlUrl)
	return
}

type Sqls struct {
	options Options
	db      *sqlx.DB
}

func (this *Sqls) Close() error {
	return this.db.Close()
}

func (this *Sqls) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.options.TimeOut)
	return
}

func (this *Sqls) QueryContext(query string, args ...interface{}) (data *sql.Rows, err error) {
	data, err = this.db.QueryContext(this.getContext(), query, args...)
	return
}

func (this *Sqls) ExecContext(query string, args ...interface{}) (data sql.Result, err error) {
	data, err = this.db.ExecContext(this.getContext(), query, args...)
	return
}

func (this *Sqls) QueryxContext(query string, args ...interface{}) (data *sqlx.Rows, err error) {
	data, err = this.db.QueryxContext(this.getContext(), query, args...)
	return
}

func (this *Sqls) QueryRowxContext(query string, args ...interface{}) (data *sqlx.Row) {
	return this.db.QueryRowxContext(this.getContext(), query, args...)
}

func (this *Sqls) GetContext(dest interface{}, query string, args ...interface{}) error {
	return this.db.GetContext(this.getContext(), dest, query, args...)
}

func (this *Sqls) SelectContext(dest interface{}, query string, args ...interface{}) error {
	return this.db.SelectContext(this.getContext(), dest, query, args...)
}
