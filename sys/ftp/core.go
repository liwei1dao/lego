package ftp

type (
	IFtp interface {
		MakeDir(path string) error
		RemoveDir(path string) error
		Close() error
	}
)

var (
	defsys IFtp
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newSys(newOptions(config, option...)); err == nil {
	}
	return
}

func NewSys(option ...Option) (sys IFtp, err error) {
	if sys, err = newSys(newOptionsByOption(option...)); err == nil {
	}
	return
}

func MakeDir(path string) error {
	return defsys.MakeDir(path)
}
func RemoveDir(path string) error {
	return defsys.RemoveDir(path)
}
func Close() error {
	return defsys.Close()
}
