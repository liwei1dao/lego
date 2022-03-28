package zookeeper

import "github.com/go-zookeeper/zk"

type (
	ISys interface {
		Exists(path string) (bool, *zk.Stat, error)
		Delete(path string, version int32) error
		Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
		Get(path string) ([]byte, *zk.Stat, error)
		Set(path string, data []byte, version int32) (*zk.Stat, error)
		GetACL(path string) ([]zk.ACL, *zk.Stat, error)
		SetACL(path string, acl []zk.ACL, version int32) (*zk.Stat, error)
		AddAuth(scheme string, auth []byte) error
		ExistsW(path string) (bool, *zk.Stat, <-chan zk.Event, error)
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
func Exists(path string) (bool, *zk.Stat, error) {
	return defsys.Exists(path)
}
func Delete(path string, version int32) error {
	return defsys.Delete(path, version)
}
func Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	return defsys.Create(path, data, flags, acl)
}
func Get(path string) ([]byte, *zk.Stat, error) {
	return defsys.Get(path)
}
func Set(path string, data []byte, version int32) (*zk.Stat, error) {
	return defsys.Set(path, data, version)
}
func GetACL(path string) ([]zk.ACL, *zk.Stat, error) {
	return defsys.GetACL(path)
}
func SetACL(path string, acl []zk.ACL, version int32) (*zk.Stat, error) {
	return defsys.SetACL(path, acl, version)
}
func AddAuth(scheme string, auth []byte) error {
	return defsys.AddAuth(scheme, auth)
}
func ExistsW(path string) (bool, *zk.Stat, <-chan zk.Event, error) {
	return defsys.ExistsW(path)
}
