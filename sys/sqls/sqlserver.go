package sqls

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (this *SqlServer) QueryContext(query string, args ...interface{}) (data *sql.Rows, err error) {
	data, err = this.db.QueryContext(this.getContext(), query, args...)
	return
}

func (this *SqlServer) ExecContext(procName string, agrs ...sql.NamedArg) (data *sql.Rows, err error) {
	sql := this.getExecProcSql(procName, agrs...)
	return this.QueryContext(sql)
}

func (this *SqlServer) getExecProcSql(procName string, agrs ...sql.NamedArg) string {
	var sql string
	var paras string
	for _, v := range agrs {
		paras += fmt.Sprintf("%s=%v", v.Name, v.Value)
		paras += ","
	}
	// for key, value := range agr {
	// 	paras += key + " = " + value
	// 	paras += ","
	// }
	paras = strings.TrimRight(paras, ",")
	paras += ";"
	sql = fmt.Sprintf("DECLARE    @return_value int; EXEC    @return_value = [dbo].[%s] %s SELECT 'Return Value' = @return_value", procName, paras)
	return sql
}
