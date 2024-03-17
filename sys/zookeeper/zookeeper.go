package zookeeper

import (
	"time"

	"github.com/go-zookeeper/zk"
)

func newSys(options Options) (sys *Zookeeper, err error) {
	sys = &Zookeeper{options: options}
	err = sys.init()
	return
}

type Zookeeper struct {
	options Options
	conn    *zk.Conn
}

func (this *Zookeeper) init() (err error) {
	if this.conn, _, err = zk.Connect(this.options.Adds, time.Second*time.Duration(this.options.TimeOut)); err != nil {
		return
	}
	return
}
func (this *Zookeeper) Exists(path string) (bool, *zk.Stat, error) {
	return this.conn.Exists(path)
}

func (this *Zookeeper) Delete(path string, version int32) error {
	return this.conn.Delete(path, version)
}

func (this *Zookeeper) Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	return this.conn.Create(path, data, flags, acl)
}

func (this *Zookeeper) Get(path string) ([]byte, *zk.Stat, error) {
	return this.conn.Get(path)
}

func (this *Zookeeper) Set(path string, data []byte, version int32) (*zk.Stat, error) {
	return this.conn.Set(path, data, version)
}

func (this *Zookeeper) GetACL(path string) ([]zk.ACL, *zk.Stat, error) {
	return this.conn.GetACL(path)
}

func (this *Zookeeper) SetACL(path string, acl []zk.ACL, version int32) (*zk.Stat, error) {
	return this.conn.SetACL(path, acl, version)
}

func (this *Zookeeper) AddAuth(scheme string, auth []byte) error {
	return this.conn.AddAuth(scheme, auth)
}

func (this *Zookeeper) ExistsW(path string) (bool, *zk.Stat, <-chan zk.Event, error) {
	return this.conn.ExistsW(path)
}
