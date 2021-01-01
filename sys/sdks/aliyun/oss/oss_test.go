package oss

import (
	"testing"
)

func Test_OSSUploadFile(t *testing.T) {
	if err := OnInit(nil,
		SetEndpoint("http://oss-cn-hongkong.aliyuncs.com"),
		SetAccessKeyId("xxxxxxxx"),
		SetAccessKeySecret("xxxxxxxxx"),
		SetBucketName("xxxxxxxxx"),
	); err != nil {
		t.Logf("初始化OSS 系统失败 err:%s", err.Error())
	} else {
		t.Logf("初始化OSS 成功")
	}
	// if err := CreateBucket("hitoolchat"); err != nil {
	// 	t.Logf("创建 CreateBucket  err:%s", err.Error())
	// } else {
	// 	t.Logf("创建 CreateBucket 成功")
	// }
	// if err := UploadFile("test/liwei1dao.jpg", "F:/liwei1dao.jpg"); err != nil {
	// 	t.Logf("上传OSS 系统失败 err:%s", err.Error())
	// } else {
	// 	t.Logf("上传OSS 成功")
	// }
	// if file, err := GetObject("test/liwei1dao.jpg"); err != nil {
	// 	t.Logf("下载OSS 系统失败 err:%s", err.Error())
	// } else {
	// 	t.Logf("下载OSS 成功 len:%d", len(file))
	// }
}
