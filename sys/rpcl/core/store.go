package core

import (
	"crypto/tls"
	"strings"
	"time"
)

type Config struct {
	ClientTLS         *ClientTLSConfig
	TLS               *tls.Config
	ConnectionTimeout time.Duration
	Bucket            string
	PersistConnection bool
	Username          string
	Password          string
}

type ClientTLSConfig struct {
	CertFile   string
	KeyFile    string
	CACertFile string
}

type IStore interface {
	Put(key string, value []byte, options *WriteOptions) error
	Get(key string) (*KVPair, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*KVPair, error)
	Close()
}

type KVPair struct {
	Key       string
	Value     []byte
	LastIndex uint64
}
type WriteOptions struct {
	IsDir bool
	TTL   time.Duration
}
type LockOptions struct {
	Value     []byte
	TTL       time.Duration
	RenewLock chan struct{}
}
type Locker interface {
	Lock(stopChan chan struct{}) (<-chan struct{}, error)
	Unlock() error
}

func Normalize(key string) string {
	return "/" + join(SplitKey(key))
}

func SplitKey(key string) (path []string) {
	if strings.Contains(key, "/") {
		path = strings.Split(key, "/")
	} else {
		path = []string{key}
	}
	return path
}

func join(parts []string) string {
	return strings.Join(parts, "/")
}
