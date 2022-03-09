package ftp

import (
	"fmt"
	"testing"
)

func Test_sys(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"IP":       "172.20.27.145",
		"Port":     21,
		"User":     "zmlftp",
		"Password": "123456",
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		fmt.Printf("start sys succ")
		err = MakeDir("/ftptest")
		fmt.Printf("start sys MakeDir err:%v", err)

	}
}
