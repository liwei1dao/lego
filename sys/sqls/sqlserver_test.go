package sqls

import (
	// "database/sql"

	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func TestSqlserverInit(t *testing.T) {
	err := OnInit(nil,
		SetSqlUrl("sqlserver://sa:Jiangnancyhd2017@119.23.148.95:1433?database=THAccountsDB&connection+timeout=30"),
	)
	if err != nil {
		t.Logf("初始化失败=%s", err.Error())
	} else {
		// m1 := make(map[string]string)
		// m1["@Account"] = "liwei1dao"
		// m1["@Password"] = "li13451234"

		// sql := sqlserver.GetExecProcSql("Login", m1)

		row, err := QueryContext("Login", sql.Named("@Account", "liwei1dao"), sql.Named("@Password", "li13451234"))
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
	}
}

func TestLoginAccount(t *testing.T) {
	err := OnInit(nil,
		SetSqlUrl("sqlserver://sa:Jiangnancyhd2017@119.23.148.95:1433?database=THAccountsDB&connection+timeout=30"),
	)
	if err != nil {
		t.Logf("初始化失败=%s", err.Error())
	} else {
		var strErrorDescribe = ""
		if data, err := QueryContext("GSP_MB_EfficacyAccounts",
			sql.Named("LogonType", 1),
			sql.Named("cbDeviceType", 1),
			sql.Named("strAccounts", ""),
			sql.Named("strPassword", ""),
			sql.Named("strMobilePhone", "18888755166"),
			sql.Named("strIPhonePassword", "260095"),
			sql.Named("strMachineID", "ios"),
			sql.Named("strErrorDescribe", sql.Out{Dest: &strErrorDescribe, In: false}),
		); err != nil {
			t.Logf("执行存储工程失败:%s", err.Error())
		} else {
			t.Log("执行存储工程成功", data)
			//获取结果集
			for data.Next() {
				var UserID, GameID, Accounts, NickName, DynamicPass, UnderWrite, FaceID, CustomID, Gender, Ingot, Experience, Score, Insure, Beans, LoveLiness, MemberOrder, MemberOverDate, MoorMachine, InsureEnabled, LogonMode, HeadUrl, BingID, UserUin, xinyong = 100040, 0, "", "", "", "", 0, 0, 1, 0, 9412, 0, 0, 1000.00, 0, 0, time.Now(), 0, 1, "", 0, 0, "", 5
				data.Scan(&UserID, &GameID, &Accounts, &NickName, &DynamicPass, &UnderWrite, &FaceID, &CustomID, &Gender, &Ingot, &Experience, &Score, &Insure, &Beans, &LoveLiness, &MemberOrder, &MemberOverDate, &MoorMachine, &InsureEnabled, &LogonMode, &HeadUrl, &BingID, &UserUin, &xinyong)
				fmt.Printf("%d,%d,%s,%s", UserID, GameID, Accounts, NickName)
			}
		}
	}
}
