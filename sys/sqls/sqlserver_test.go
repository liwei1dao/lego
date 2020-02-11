package sqls

import "testing"

func TestSqlserverInit(t *testing.T) {
	err := OnInit(nil,
		SetServerAdd("119.23.148.95"),
		SetPort(1433),
		SetUser("sa"),
		SetPassword("Jiangnancyhd2017"),
		SetDatabase("ReportServer"),
	)
	if err != nil {
		t.Logf("初始化失败=%s", err.Error())
	} else {
		t.Logf("初始化成功")
	}
}
