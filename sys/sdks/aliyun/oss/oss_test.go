package oss

import (
	"fmt"
	"testing"
)

func Test_OSSUploadFile(t *testing.T) {
	if err := OnInit(nil,
		SetEndpoint("http://gohitool.oss-accelerate.aliyuncs.com"),
		SetAccessKeyId("LTAI4G1hvDpFe6gP7QyDeJK7"),
		SetAccessKeySecret("mJfDBKS4GewCPwpPsDWaqYhYGf1qUZ"),
		SetBucketName("gohitool"),
	); err != nil {
		fmt.Printf("初始化OSS 系统失败 err:%v", err)
		t.Logf("初始化OSS 系统失败 err:%s", err.Error())
		return
	} else {
		t.Logf("初始化OSS 成功")
	}
	// if err := CreateBucket("hitoolchat"); err != nil {
	// 	t.Logf("创建 CreateBucket  err:%s", err.Error())
	// } else {
	// 	t.Logf("创建 CreateBucket 成功")
	// }
	if err := UploadFile("test/liwei2dao.jpg", "F:/liwei1dao.jpg"); err != nil {
		t.Logf("上传OSS 系统失败 err:%s", err.Error())
	} else {
		t.Logf("上传OSS 成功")
	}
	// if file, err := GetObject("test/liwei1dao.jpg"); err != nil {
	// 	t.Logf("下载OSS 系统失败 err:%s", err.Error())
	// } else {
	// 	t.Logf("下载OSS 成功 len:%d", len(file))
	// }
}
