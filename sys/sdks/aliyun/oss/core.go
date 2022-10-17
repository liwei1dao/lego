package oss

import (
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type (
	ISys interface {
		///创建数据桶
		CreateBucket(bucketName string) (err error)
		///上传对象
		UploadObject(objectKey string, reader io.Reader, options ...oss.Option) (err error)
		///上传本地文件
		UploadFile(objectName string, localFileName string) (err error)
		///获取对象
		GetObject(objectName string, options ...oss.Option) ([]byte, error)
		///下载文件
		DownloadFile(objectName string, downloadedFileName string) (err error)
		///删除文件
		DeleteFile(objectName string) (err error)
		//获取临时访问地址
		GetURL(objectName string, expired int64, options ...oss.Option) (url string, err error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func CreateBucket(bucketName string) (err error) {
	return defsys.CreateBucket(bucketName)
}

func UploadObject(objectKey string, reader io.Reader, options ...oss.Option) (err error) {
	return defsys.UploadObject(objectKey, reader, options...)
}

func UploadFile(localFileName string, objectName string) (err error) {
	return defsys.UploadFile(localFileName, objectName)
}

func GetObject(objectName string, options ...oss.Option) ([]byte, error) {
	return defsys.GetObject(objectName, options...)
}

func DownloadFile(objectName string, downloadedFileName string) (err error) {
	return defsys.DownloadFile(objectName, downloadedFileName)
}

func DeleteFile(objectName string) (err error) {
	return defsys.DeleteFile(objectName)
}

func GetURL(objectName string, expired int64, options ...oss.Option) (url string, err error) {
	return defsys.GetURL(objectName, expired, options...)
}
