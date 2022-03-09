package ftp

import (
	"fmt"
	"time"

	"github.com/jlaffaye/ftp"
)

func newSys(options Options) (sys *Ftp, err error) {
	sys = &Ftp{options: options}
	err = sys.init()
	return
}

type Ftp struct {
	options Options
	conn    *ftp.ServerConn
}

func (this *Ftp) init() (err error) {
	if this.conn, err = ftp.Dial(fmt.Sprintf("%s:%d", this.options.IP, this.options.Port), ftp.DialWithTimeout(5*time.Second)); err != nil {
		err = fmt.Errorf("ftp Dial err:%v", err)
		return
	}
	if err = this.conn.Login(this.options.User, this.options.Password); err != nil {
		err = fmt.Errorf("ftp Login err:%v", err)
	}
	return
}

func (this *Ftp) List(path string) (entries []*ftp.Entry, err error) {
	return this.conn.List(path)
}

func (this *Ftp) MakeDir(path string) error {
	return this.conn.MakeDir(path)
}

func (this *Ftp) RemoveDir(path string) error {
	return this.conn.RemoveDir(path)
}

func (this *Ftp) Close() error {
	return this.conn.Quit()
}
