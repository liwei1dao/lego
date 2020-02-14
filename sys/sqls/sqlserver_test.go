package sqls

import (
	// "database/sql"

	"database/sql"
	"fmt"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
)

func TestSqlserverInit(t *testing.T) {
	err := OnInit(nil,
		SetServerAdd("119.23.148.95"),
		SetPort(1433),
		SetUser("sa"),
		SetPassword("Jiangnancyhd2017"),
		SetDatabase("liwei1dao"),
	)
	if err != nil {
		t.Logf("初始化失败=%s", err.Error())
	} else {
		// m1 := make(map[string]string)
		// m1["@Account"] = "liwei1dao"
		// m1["@Password"] = "li13451234"

		// sql := sqlserver.GetExecProcSql("Login", m1)

		row, err := ExecContext("Login", sql.Named("@Account", "liwei1dao"), sql.Named("@Password", "li13451234"))
		// // var rs mssql.ReturnStatus
		// row, err := sqlserver.db.QueryContext(sqlserver.getContext(), sql)
		if err != nil {
			t.Logf("执行存储工程失败:%s", err.Error())
		} else {
			t.Log("执行存储工程成功", row)
			//获取结果集
			for row.Next() {
				var uid int
				var a string
				var p string
				var n string
				row.Scan(&uid, &a, &p, &n)
				fmt.Printf("%d,%s,%s,%s", uid, a, p, n)
			}
		}
		// data, err := QueryContext("SELECT * FROM dbo.[user]")
		// if err != nil {
		// 	t.Logf("执行存储工程失败:%s", err.Error())
		// } else {
		// 	t.Log("执行存储工程失败", data)
		// }
	}
}
