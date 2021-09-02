package sql

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
)

///测试SqlServer 存储过程
func Test_SqlServer(t *testing.T) {
	err := OnInit(nil,
		SetSqlUrl("sqlserver://sa:Jiangnancyhd2017@119.23.148.95:1433?database=THAccountsDB&connection+timeout=30"),
	)
	if err != nil {
		t.Logf("初始化失败=%s", err.Error())
	} else {
		row, err := QueryContext("Login", sql.Named("@Account", "liwei1dao"), sql.Named("@Password", "li13451234"))
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
	}
}

func Test_MySql(t *testing.T) {
	err := OnInit(nil,
		SetSqlType(MySql),
		SetSqlUrl("root:Idss@sjzt2021@tcp(172.20.27.125:3306)/mysql"),
	)
	if err != nil {
		t.Logf("初始化失败=%s", err.Error())
	} else {
		t.Logf("初始化成功")
	}
}
