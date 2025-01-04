package discovery

import "time"



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
