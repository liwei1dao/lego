package ftp

type (
	FileEntry struct {
		FileName string
		FileSize uint64
		ModTime  int64
		IsDir    bool
	}
	ISys interface {
		ReadDir(path string) (entries []*FileEntry, err error)
		MakeDir(path string) error
		RemoveDir(path string) error
		Close() error
	}
)

var (
	defsys ISys
)

func newSys(options Options) (sys ISys, err error) {
	if options.SType == FTP {
		sys, err = newFtp(options)
	} else if options.SType == SFTP {
		sys, err = newSFtp(options)
	}
	return
}

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newSys(newOptions(config, option...)); err == nil {
	}
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	if sys, err = newSys(newOptionsByOption(option...)); err == nil {
	}
	return
}

func ReadDir(path string) (entries []*FileEntry, err error) {
	return defsys.ReadDir(path)
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
