package discovery

import (
	"crypto/tls"
	"errors"
	"strings"
	"time"

	"github.com/liwei1dao/lego/core"
)

type (
	IDiscovery interface {
		Start(node core.IServiceNode) error
		Stop() error
		GetServices() []*KV
		WatchService() chan []*KV
	}
	ServiceDiscoveryFilter func(kvp *KV) bool
	ServiceDiscovery       interface {
		GetServices() []*KV
		WatchService() chan []*KV
		RemoveWatcher(ch chan []*KV)
		Clone(servicePath string) (ServiceDiscovery, error)
		SetFilter(ServiceDiscoveryFilter)
		Close()
	}
)

/*
发现系统:服务发现
*/

var (
	// ErrBackendNotSupported is thrown when the backend k/v store is not supported by libkv
	ErrBackendNotSupported = errors.New("Backend storage not supported yet, please choose one of")
	// ErrCallNotSupported is thrown when a method is not implemented/supported by the current backend
	ErrCallNotSupported = errors.New("The current call is not supported with this backend")
	// ErrNotReachable is thrown when the API cannot be reached for issuing common store operations
	ErrNotReachable = errors.New("Api not reachable")
	// ErrCannotLock is thrown when there is an error acquiring a lock on a key
	ErrCannotLock = errors.New("Error acquiring the lock")
	// ErrKeyModified is thrown during an atomic operation if the index does not match the one in the store
	ErrKeyModified = errors.New("Unable to complete atomic operation, key modified")
	// ErrKeyNotFound is thrown when the key is not found in the store during a Get operation
	ErrKeyNotFound = errors.New("Key not found in store")
	// ErrPreviousNotSpecified is thrown when the previous value is not specified for an atomic operation
	ErrPreviousNotSpecified = errors.New("Previous K/V pair should be provided for the Atomic operation")
	// ErrKeyExists is thrown when the previous value exists in the case of an AtomicPut
	ErrKeyExists = errors.New("Previous K/V pair exists, cannot complete Atomic operation")
)

type Backend string

const (
	// CONSUL backend
	CONSUL Backend = "consul"
	// ZK backend
	ZK Backend = "zk"
	// REDIS backend
	REDIS Backend = "redis"
)

const (
	// ETCD backend
	ETCD Backend = "etcd"
	// ETCDV3 backend
	ETCDV3 Backend = "etcdv3"
	// ETCDV3 Single backend
	ETCDV3_SINGLE Backend = "etcdv3_single"
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
	List(directory string) ([]*KVPair, error)
	NewLock(key string, options *LockOptions) (Locker, error)
	WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*KVPair, error)
	DeleteTree(directory string) error
	AtomicPut(key string, value []byte, previous *KVPair, options *WriteOptions) (bool, *KVPair, error)
	AtomicDelete(key string, previous *KVPair) (bool, error)
	Close()
}
type KV struct {
	Key   string
	Value string
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
