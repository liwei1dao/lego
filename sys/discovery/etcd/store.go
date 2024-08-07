package etcd

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/discovery/dcore"
	"github.com/rpcxio/libkv/store"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const defaultTTL = 30

var EtcdConfigAutoSyncInterval = time.Minute * 5
var (
	// ErrMultipleEndpointsUnsupported is thrown when there are
	// multiple endpoints specified for Consul
	ErrMultipleEndpointsUnsupported = errors.New("etcd does not support multiple endpoints")

	// ErrSessionRenew is thrown when the session can't be
	// renewed because the Consul version does not support sessions
	ErrSessionRenew = errors.New("cannot set or renew session for ttl, unable to operate on sessions")
)

func NewEtcdStore(address []string, options *dcore.Config) (store *ETCDV3Store, err error) {
	store = &ETCDV3Store{
		done:           make(chan struct{}),
		startKeepAlive: make(chan struct{}),
		ttl:            defaultTTL,
	}

	cfg := clientv3.Config{
		Endpoints: address,
	}

	if options != nil {
		store.timeout = options.ConnectionTimeout
		cfg.DialTimeout = options.ConnectionTimeout
		cfg.DialKeepAliveTimeout = options.ConnectionTimeout
		cfg.TLS = options.TLS
		cfg.Username = options.Username
		cfg.Password = options.Password
		cfg.AutoSyncInterval = EtcdConfigAutoSyncInterval
	}
	if store.timeout == 0 {
		store.timeout = 10 * time.Second
	}
	store.cfg = cfg
	if err = store.init(true); err != nil {
		return nil, err
	}
	go store.keepAlive()
	return
}

type ETCDV3Store struct {
	timeout        time.Duration
	client         *clientv3.Client
	leaseID        clientv3.LeaseID
	cfg            clientv3.Config
	done           chan struct{}
	startKeepAlive chan struct{}

	AllowKeyNotFound bool

	mu  sync.RWMutex
	ttl int64
}

func (this *ETCDV3Store) init(grant bool) error {
	cli, err := clientv3.New(this.cfg)
	if err != nil {
		return err
	}

	this.client = cli

	if grant {

		this.mu.RLock()
		err = this.grant(this.ttl)
		this.mu.RUnlock()
	}

	return err
}

func (this *ETCDV3Store) keepAlive() {
	var ch <-chan *clientv3.LeaseKeepAliveResponse
	var err error
rekeepalive:
	for {
		if this.leaseID != 0 {
			ch, err = this.client.KeepAlive(context.Background(), this.leaseID)
		}
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	// KeepAlive success
	for {
		select {
		case <-this.done:
			return
		case resp := <-ch: // KeepAlive channel is closed
			if resp == nil { // connection is closed
				this.client.Close()
				for {
					select {
					case <-this.done:
						return
					default:
						err = this.init(false)
						if err != nil {
							time.Sleep(time.Second)
							continue
						}

						this.mu.RLock()
						err = this.grant(this.ttl)
						this.mu.RUnlock()
						if err != nil {
							this.client.Close()
							time.Sleep(time.Second)
							continue
						}
						goto rekeepalive
					}
				}
			}
		}
	}
}

// grant a lease.
func (s *ETCDV3Store) grant(ttl int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	resp, err := s.client.Grant(ctx, ttl)
	cancel()
	if err == nil {
		s.leaseID = resp.ID
	}
	return err
}

// Put a value at the specified key
func (s *ETCDV3Store) Put(key string, value []byte, options *dcore.WriteOptions) error {
	var ttl int64
	if options != nil {
		ttl = int64(options.TTL.Seconds())
	}
	if ttl == 0 {
		ttl = defaultTTL
	}
	s.mu.Lock()
	s.ttl = ttl
	s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	_, err := s.client.Put(ctx, key, string(value), clientv3.WithLease(s.leaseID))
	cancel()

	return err
}

// Get a value given its key
func (s *ETCDV3Store) Get(key string) (*dcore.KVPair, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	resp, err := s.client.Get(ctx, key)
	cancel()
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, store.ErrKeyNotFound
	}

	pair := &dcore.KVPair{
		Key:       key,
		Value:     resp.Kvs[0].Value,
		LastIndex: uint64(resp.Kvs[0].Version),
	}

	return pair, nil
}

// Delete the value at the specified key
func (s *ETCDV3Store) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	_, err := s.client.Delete(ctx, key)
	cancel()

	return err
}

// Exists verifies if a Key exists in the store
func (s *ETCDV3Store) Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	resp, err := s.client.Get(ctx, key)
	cancel()
	if err != nil {
		return false, err
	}

	return len(resp.Kvs) != 0, nil
}

// Watch for changes on a key.
func (this *ETCDV3Store) Watch(key string, stopCh <-chan struct{}) (<-chan *dcore.KVPair, error) {
	watchCh := make(chan *dcore.KVPair)

	go func() {
		defer close(watchCh)

		// put the current value into returned channel before watch
		pair, err := this.Get(key)
		if err != nil {
			return
		}
		watchCh <- pair

		rch := this.client.Watch(context.Background(), key)
		for {
			select {
			case <-this.done:
				return
			case wresp, ok := <-rch:
				if !ok || wresp.Canceled { // watch is canceled
					return
				}
				for _, event := range wresp.Events {
					watchCh <- &dcore.KVPair{
						Key:       string(event.Kv.Key),
						Value:     event.Kv.Value,
						LastIndex: uint64(event.Kv.Version),
					}
				}
			}
		}
	}()

	return watchCh, nil
}

// WatchTree watches for changes on child nodes under a given directory
func (this *ETCDV3Store) WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*dcore.KVPair, error) {
	watchCh := make(chan []*dcore.KVPair)
	list, err := this.List(directory)
	if err != nil {
		if !this.AllowKeyNotFound || err != store.ErrKeyNotFound {
			return watchCh, err
		}
	}
	go func() {
		defer close(watchCh)

		watchCh <- list

		rch := this.client.Watch(context.Background(), directory, clientv3.WithPrefix())
		for {
			select {
			case <-this.done:
				return
			case resp := <-rch:
				if resp.Canceled { // watch is canceled
					return
				}

				list, err := this.List(directory)
				if err != nil {
					if !this.AllowKeyNotFound || err != store.ErrKeyNotFound {
						continue
					}
				}
				watchCh <- list
			}
		}
	}()

	return watchCh, nil
}

// List the content of a given prefix
func (this *ETCDV3Store) List(directory string) ([]*dcore.KVPair, error) {
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()

	resp, err := this.client.Get(ctx, directory, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	kvpairs := make([]*dcore.KVPair, 0, len(resp.Kvs))

	if len(resp.Kvs) == 0 {
		return nil, store.ErrKeyNotFound
	}

	for _, kv := range resp.Kvs {
		pair := &dcore.KVPair{
			Key:       string(kv.Key),
			Value:     kv.Value,
			LastIndex: uint64(kv.Version),
		}
		kvpairs = append(kvpairs, pair)
	}

	return kvpairs, nil
}

// DeleteTree deletes a range of keys under a given directory
func (s *ETCDV3Store) DeleteTree(directory string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	_, err := s.client.Delete(ctx, directory, clientv3.WithPrefix())
	cancel()

	return err
}

// AtomicPut CAS operation on a single value.
// Pass previous = nil to create a new key.
func (s *ETCDV3Store) AtomicPut(key string, value []byte, previous *dcore.KVPair, options *dcore.WriteOptions) (bool, *dcore.KVPair, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var revision int64
	var presp *clientv3.PutResponse
	var txresp *clientv3.TxnResponse
	var err error
	if previous == nil {
		if exist, err := s.Exists(key); err != nil { // not atomicput
			return false, nil, err
		} else if !exist {
			presp, err = s.client.Put(ctx, key, string(value))
			if err != nil {
				return false, nil, err
			}
			if presp != nil {
				revision = presp.Header.GetRevision()
			}
		} else {
			return false, nil, store.ErrKeyExists
		}
	} else {

		cmps := []clientv3.Cmp{
			clientv3.Compare(clientv3.Value(key), "=", string(previous.Value)),
			clientv3.Compare(clientv3.Version(key), "=", int64(previous.LastIndex)),
		}
		txresp, err = s.client.Txn(ctx).If(cmps...).
			Then(clientv3.OpPut(key, string(value))).
			Commit()
		if txresp != nil {
			if txresp.Succeeded {
				revision = txresp.Header.GetRevision()
			} else {
				err = errors.New("key's version not matched!")
			}
		}
	}

	if err != nil {
		return false, nil, err
	}

	pair := &dcore.KVPair{
		Key:       key,
		Value:     value,
		LastIndex: uint64(revision),
	}

	return true, pair, nil
}

// AtomicDelete cas deletes a single value
func (s *ETCDV3Store) AtomicDelete(key string, previous *dcore.KVPair) (bool, error) {
	deleted := false
	var err error
	var txresp *clientv3.TxnResponse
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if previous == nil {
		return false, errors.New("key's version info is needed!")
	} else {
		cmps := []clientv3.Cmp{
			clientv3.Compare(clientv3.Value(key), "=", string(previous.Value)),
			clientv3.Compare(clientv3.Version(key), "=", int64(previous.LastIndex)),
		}
		txresp, err = s.client.Txn(ctx).If(cmps...).
			Then(clientv3.OpDelete(key)).
			Commit()

		deleted = txresp.Succeeded
		if !deleted {
			err = errors.New("conflicts!")
		}
	}

	if err != nil {
		return false, err
	}

	return deleted, nil
}

// Close closes the client connection
func (s *ETCDV3Store) Close() {
	defer func() {
		if recover() != nil {
			// close of closed channel panic occur
		}
	}()
	close(s.done)
	s.client.Close()
}
