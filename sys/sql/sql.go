package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	// 更具需要启动的数据库服务 自己在上层业务添加驱动代码
	// _ "github.com/denisenkom/go-mssqldb" //sqlserver 驱动 需要提前导入否赠会连接失败
	// _ "github.com/go-sql-driver/mysql"   //mysql 驱动
	// _ "github.com/godror/godror"		    //oracle 驱动
)

func newSys(options Options) (sys *Sql, err error) {
	sys = &Sql{options: options}
	sys.db, err = sqlx.Connect(string(options.SqlType), options.SqlUrl)
	return
}

type Sql struct {
	options Options
	db      *sqlx.DB
}

func (this *Sql) Close() error {
	return this.db.Close()
}

func (this *Sql) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), time.Second*time.Duration(this.options.TimeOut))
	return
}

func (this *Sql) Query(query string, args ...interface{}) (data *sql.Rows, err error) {
	data, err = this.db.Query(query, args...)
	return
}

func (this *Sql) Exec(query string, args ...interface{}) (data sql.Result, err error) {
	data, err = this.db.Exec(query, args...)
	return
}

func (this *Sql) QueryContext(query string, args ...interface{}) (data *sql.Rows, err error) {
	data, err = this.db.QueryContext(this.getContext(), query, args...)
	return
}

func (this *Sql) ExecContext(query string, args ...interface{}) (data sql.Result, err error) {
	data, err = this.db.ExecContext(this.getContext(), query, args...)
	return
}

func (this *Sql) QueryxContext(query string, args ...interface{}) (data *sqlx.Rows, err error) {
	data, err = this.db.QueryxContext(this.getContext(), query, args...)
	return
}

func (this *Sql) QueryRowxContext(query string, args ...interface{}) (data *sqlx.Row) {
	return this.db.QueryRowxContext(this.getContext(), query, args...)
}

func (this *Sql) GetContext(dest interface{}, query string, args ...interface{}) error {
	return this.db.GetContext(this.getContext(), dest, query, args...)
}

func (this *Sql) SelectContext(dest interface{}, query string, args ...interface{}) error {
	return this.db.SelectContext(this.getContext(), dest, query, args...)
}
