package consul

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/liwei1dao/lego/sys/discovery/dcore"
	"github.com/liwei1dao/lego/sys/log"
)

const (
	// DefaultWatchWaitTime is how long we block for at a
	// time to check if the watched key has changed. This
	// affects the minimum time it takes to cancel a watch.
	DefaultWatchWaitTime = 15 * time.Second

	// RenewSessionRetryMax is the number of time we should try
	// to renew the session before giving up and throwing an error
	RenewSessionRetryMax = 5

	// MaxSessionDestroyAttempts is the maximum times we will try
	// to explicitely destroy the session attached to a lock after
	// the connectivity to the store has been lost
	MaxSessionDestroyAttempts = 5

	// defaultLockTTL is the default ttl for the consul lock
	defaultLockTTL = 20 * time.Second
)

var (
	// ErrMultipleEndpointsUnsupported is thrown when there are
	// multiple endpoints specified for Consul
	ErrMultipleEndpointsUnsupported = errors.New("consul does not support multiple endpoints")

	// ErrSessionRenew is thrown when the session can't be
	// renewed because the Consul version does not support sessions
	ErrSessionRenew = errors.New("cannot set or renew session for ttl, unable to operate on sessions")
)

func NewConsulStore(address []string, options *dcore.Config) (store *ConsulStore, err error) {
	if len(address) > 1 {
		return nil, ErrMultipleEndpointsUnsupported
	}

	config := api.DefaultConfig()
	config.HttpClient = http.DefaultClient
	config.Address = address[0]
	config.Scheme = "http"
	if options != nil {
		if options.TLS != nil {
			config.HttpClient.Transport = &http.Transport{
				TLSClientConfig: options.TLS,
			}
		}
		if options.ConnectionTimeout != 0 {
			config.WaitTime = options.ConnectionTimeout
		}
	}
	var client *api.Client
	client, err = api.NewClient(config)
	if err != nil {
		return
	}
	store = &ConsulStore{
		config: config,
		client: client,
	}
	return
}

type ConsulStore struct {
	sync.Mutex
	config *api.Config
	client *api.Client
}

func (this *ConsulStore) Get(key string) (*dcore.KVPair, error) {
	options := &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}

	pair, meta, err := this.client.KV().Get(this.normalize(key), options)
	if err != nil {
		return nil, err
	}

	// If pair is nil then the key does not exist
	if pair == nil {
		return nil, dcore.ErrKeyNotFound
	}

	return &dcore.KVPair{Key: pair.Key, Value: pair.Value, LastIndex: meta.LastIndex}, nil
}
func (this *ConsulStore) Put(key string, value []byte, opts *dcore.WriteOptions) error {
	key = this.normalize(key)

	p := &api.KVPair{
		Key:   key,
		Value: value,
		Flags: api.LockFlagValue,
	}

	if opts != nil && opts.TTL > 0 {
		// Create or renew a session holding a TTL. Operations on sessions
		// are not deterministic: creating or renewing a session can fail
		for retry := 1; retry <= RenewSessionRetryMax; retry++ {
			err := this.renewSession(p, opts.TTL)
			if err == nil {
				break
			}
			log.Errorln(err)
			if retry == RenewSessionRetryMax {
				return ErrSessionRenew
			}
		}
	}

	_, err := this.client.KV().Put(p, nil)
	return err
}
func (this *ConsulStore) Delete(key string) error {
	if _, err := this.Get(key); err != nil {
		return err
	}
	_, err := this.client.KV().Delete(this.normalize(key), nil)
	return err
}
func (this *ConsulStore) Exists(key string) (bool, error) {
	_, err := this.Get(key)
	if err != nil {
		if err == dcore.ErrKeyNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (this *ConsulStore) List(directory string) ([]*dcore.KVPair, error) {
	pairs, _, err := this.client.KV().List(this.normalize(directory), &api.QueryOptions{WaitTime: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	if len(pairs) == 0 {
		return nil, dcore.ErrKeyNotFound
	}

	kv := []*dcore.KVPair{}

	for _, pair := range pairs {
		if pair.Key == directory {
			continue
		}
		kv = append(kv, &dcore.KVPair{
			Key:       pair.Key,
			Value:     pair.Value,
			LastIndex: pair.ModifyIndex,
		})
	}

	return kv, nil
}
func (this *ConsulStore) DeleteTree(directory string) error {
	if _, err := this.List(directory); err != nil {
		return err
	}
	_, err := this.client.KV().DeleteTree(this.normalize(directory), nil)
	return err
}
func (this *ConsulStore) WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*dcore.KVPair, error) {
	kv := this.client.KV()
	watchCh := make(chan []*dcore.KVPair)

	go func() {
		defer close(watchCh)

		// Use a wait time in order to check if we should quit
		// from time to time.
		opts := &api.QueryOptions{WaitTime: DefaultWatchWaitTime}
		for {
			// Check if we should quit
			select {
			case <-stopCh:
				return
			default:
			}

			// Get all the childrens
			pairs, meta, err := kv.List(directory, opts)
			if err != nil {
				return
			}

			// If LastIndex didn't change then it means `Get` returned
			// because of the WaitTime and the child keys didn't change.
			if opts.WaitIndex == meta.LastIndex {
				continue
			}
			opts.WaitIndex = meta.LastIndex

			// Return children KV pairs to the channel
			kvpairs := []*dcore.KVPair{}
			for _, pair := range pairs {
				if pair.Key == directory {
					continue
				}
				kvpairs = append(kvpairs, &dcore.KVPair{
					Key:       pair.Key,
					Value:     pair.Value,
					LastIndex: pair.ModifyIndex,
				})
			}
			watchCh <- kvpairs
		}
	}()

	return watchCh, nil
}
func (this *ConsulStore) Close() {
	return
}
func (this *ConsulStore) normalize(key string) string {
	key = dcore.Normalize(key)
	return strings.TrimPrefix(key, "/")
}
func (this *ConsulStore) renewSession(pair *api.KVPair, ttl time.Duration) error {
	// Check if there is any previous session with an active TTL
	session, err := this.getActiveSession(pair.Key)
	if err != nil {
		return err
	}

	if session == "" {
		entry := &api.SessionEntry{
			Behavior:  api.SessionBehaviorDelete, // Delete the key when the session expires
			TTL:       (ttl / 2).String(),        // Consul multiplies the TTL by 2x
			LockDelay: 1 * time.Millisecond,      // Virtually disable lock delay
		}

		// Create the key session
		session, _, err = this.client.Session().Create(entry, nil)
		if err != nil {
			return err
		}

		lockOpts := &api.LockOptions{
			Key:     pair.Key,
			Session: session,
		}

		// Lock and ignore if lock is held
		// It's just a placeholder for the
		// ephemeral behavior
		lock, _ := this.client.LockOpts(lockOpts)
		if lock != nil {
			lock.Lock(nil)
		}
	}

	_, _, err = this.client.Session().Renew(session, nil)
	return err
}
func (this *ConsulStore) getActiveSession(key string) (string, error) {
	pair, _, err := this.client.KV().Get(key, nil)
	if err != nil {
		return "", err
	}
	if pair != nil && pair.Session != "" {
		return pair.Session, nil
	}
	return "", nil
}
