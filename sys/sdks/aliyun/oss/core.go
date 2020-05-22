package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/liwei1dao/lego/core"
)

var (
	service core.IService
	defoss  IOSS
)

type (
	IOSS interface {
		CreateBucket(bucketName string) (err error)
		UploadFile(objectName string, localFileName string) (err error)
		GetObject(objectName string, options ...oss.Option) ([]byte, error)
		DownloadFile(objectName string, downloadedFileName string) (err error)
		DeleteFile(objectName string) (err error)
	}
)

func OnInit(s core.IService, opt ...Option) (err error) {
	defoss, err = newOSS(opt...)
	return
}

func CreateBucket(bucketName string) (err error) {
	return defoss.CreateBucket(bucketName)
}

func UploadFile(localFileName string, objectName string) (err error) {
	return defoss.UploadFile(localFileName, objectName)
}

func GetObject(objectName string, options ...oss.Option) ([]byte, error) {
	return defoss.GetObject(objectName, options...)
}

func DownloadFile(objectName string, downloadedFileName string) (err error) {
	return defoss.DownloadFile(objectName, downloadedFileName)
}

func DeleteFile(objectName string) (err error) {
	return defoss.DeleteFile(objectName)
}
