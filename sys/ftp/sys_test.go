package ftp

import (
	"fmt"
	"testing"
)

func Test_sys_ftp(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"SType":    FTP,
		"IP":       "172.20.27.145",
		"Port":     21,
		"User":     "ftpuser",
		"Password": "123456",
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		fmt.Printf("start sys succ")
		entries, err := ReadDir("./")
		fmt.Printf("start sys entriesL%v err:%v", entries, err)
	}
}

func Test_sys_sftp(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"SType":    SFTP,
		"IP":       "172.20.27.145",
		"Port":     22,
		"User":     "root",
		"Password": "idss@wuhan",
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		fmt.Printf("start sys succ")
		entries, err := ReadDir("/opt/idss/gm/")
		fmt.Printf("start sys entriesL%v err:%v", entries, err)
	}
}
