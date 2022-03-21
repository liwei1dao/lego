package storage_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/sdks/google/storage"
)

func Test_Sys(t *testing.T) {
	if sys, err := storage.NewSys(storage.SetBucketName("resstation.appspot.com")); err != nil {
		fmt.Printf("Sys Init err:%v", err)
	} else {
		buf := bytes.NewBuffer([]byte("liwei1dao"))
		if err = sys.UploadFile(buf, "test.txt"); err != nil {
			fmt.Printf("Sys UploadFile err:%v", err)
		} else {
			fmt.Printf("Sys UploadFile Succ!")
		}
	}
}
