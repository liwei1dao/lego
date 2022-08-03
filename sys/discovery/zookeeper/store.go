package zookeeper

import (
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/liwei1dao/lego/sys/discovery/dcore"
)

const (
	// SOH control character
	SOH = "\x01"

	defaultTimeout = 10 * time.Second
)

type ZookeeperStore struct {
	timeout time.Duration
	client  *zk.Conn
}

func (this *ZookeeperStore) Get(key string) (pair *dcore.KVPair, err error) {
	resp, meta, err := this.client.Get(this.normalize(key))

	if err != nil {
		if err == zk.ErrNoNode {
			return nil, dcore.ErrKeyNotFound
		}
		return nil, err
	}

	if string(resp) == SOH {
		return this.Get(dcore.Normalize(key))
	}

	pair = &dcore.KVPair{
		Key:       key,
		Value:     resp,
		LastIndex: uint64(meta.Version),
	}

	return pair, nil
}

func (this *ZookeeperStore) Put(key string, value []byte, opts *dcore.WriteOptions) error {
	fkey := this.normalize(key)

	exists, err := this.Exists(key)
	if err != nil {
		return err
	}

	if !exists {
		if opts != nil && opts.TTL > 0 {
			this.createFullPath(dcore.SplitKey(strings.TrimSuffix(key, "/")), true)
		} else {
			this.createFullPath(dcore.SplitKey(strings.TrimSuffix(key, "/")), false)
		}
	}

	_, err = this.client.Set(fkey, value, -1)
	return err
}

func (this *ZookeeperStore) Delete(key string) error {
	err := this.client.Delete(this.normalize(key), -1)
	if err == zk.ErrNoNode {
		return dcore.ErrKeyNotFound
	}
	return err
}

func (this *ZookeeperStore) Exists(key string) (bool, error) {
	exists, _, err := this.client.Exists(this.normalize(key))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (this *ZookeeperStore) Watch(key string, stopCh <-chan struct{}) (<-chan *dcore.KVPair, error) {
	pair, err := this.Get(key)
	if err != nil {
		return nil, err
	}

	watchCh := make(chan *dcore.KVPair)
	go func() {
		defer close(watchCh)
		watchCh <- pair
		for {
			_, _, eventCh, err := this.client.GetW(this.normalize(key))
			if err != nil {
				return
			}
			select {
			case e := <-eventCh:
				if e.Type == zk.EventNodeDataChanged {
					if entry, err := this.Get(key); err == nil {
						watchCh <- entry
					}
				}
			case <-stopCh:
				return
			}
		}
	}()

	return watchCh, nil
}
func (this *ZookeeperStore) List(directory string) ([]*dcore.KVPair, error) {
	keys, stat, err := this.client.Children(this.normalize(directory))
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, dcore.ErrKeyNotFound
		}
		return nil, err
	}

	kv := []*dcore.KVPair{}

	// FIXME Costly Get request for each child key..
	for _, key := range keys {
		pair, err := this.Get(strings.TrimSuffix(directory, "/") + this.normalize(key))
		if err != nil {
			// If node is not found: List is out of date, retry
			if err == dcore.ErrKeyNotFound {
				return this.List(directory)
			}
			return nil, err
		}

		kv = append(kv, &dcore.KVPair{
			Key:       key,
			Value:     []byte(pair.Value),
			LastIndex: uint64(stat.Version),
		})
	}

	return kv, nil
}
func (this *ZookeeperStore) WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*dcore.KVPair, error) {
	entries, err := this.List(directory)
	if err != nil {
		return nil, err
	}

	watchCh := make(chan []*dcore.KVPair)
	go func() {
		defer close(watchCh)
		watchCh <- entries
		for {
			_, _, eventCh, err := this.client.ChildrenW(this.normalize(directory))
			if err != nil {
				return
			}
			select {
			case e := <-eventCh:
				if e.Type == zk.EventNodeChildrenChanged {
					if kv, err := this.List(directory); err == nil {
						watchCh <- kv
					}
				}
			case <-stopCh:
				return
			}
		}
	}()

	return watchCh, nil
}

func (this *ZookeeperStore) normalize(key string) string {
	key = dcore.Normalize(key)
	return strings.TrimSuffix(key, "/")
}

func (this *ZookeeperStore) createFullPath(path []string, ephemeral bool) error {
	for i := 1; i <= len(path); i++ {
		newpath := "/" + strings.Join(path[:i], "/")
		if i == len(path) && ephemeral {
			_, err := this.client.Create(newpath, []byte{}, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			return err
		}
		_, err := this.client.Create(newpath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			// Skip if node already exists
			if err != zk.ErrNodeExists {
				return err
			}
		}
	}
	return nil
}
