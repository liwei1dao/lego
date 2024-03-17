package storage

import "io"

type (
	ISys interface {
		UploadFile(r io.Reader, object string) (err error)
		DownloadFile(w io.Writer, destFileName string) (err error)
	}
)

var defsys ISys

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func UploadFile(r io.Reader, object string) (err error) {
	return defsys.UploadFile(r, object)
}
func DownloadFile(w io.Writer, object string) (err error) {
	return defsys.DownloadFile(w, object)
}
