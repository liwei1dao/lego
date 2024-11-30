package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/liwei1dao/lego/sys/discovery/dcore"
	lgredis "github.com/liwei1dao/lego/sys/redis"
	"github.com/redis/go-redis/v9"
	"github.com/rpcxio/libkv/store"
)

var (
	// ErrMultipleEndpointsUnsupported is thrown when there are
	// multiple endpoints specified for Redis
	ErrMultipleEndpointsUnsupported = errors.New("redis: does not support multiple endpoints")

	// ErrTLSUnsupported is thrown when tls config is given
	ErrTLSUnsupported = errors.New("redis does not support tls")

	// ErrAbortTryLock is thrown when a user stops trying to seek the lock
	// by sending a signal to the stop chan, this is used to verify if the
	// operation succeeded
	ErrAbortTryLock = errors.New("redis: lock operation aborted")
)

func New(address []string, options *dcore.Config) (*RedisStore, error) {
	var password string
	if len(address) > 1 {
		return nil, ErrMultipleEndpointsUnsupported
	}
	if options != nil && options.TLS != nil {
		return nil, ErrTLSUnsupported
	}
	if options != nil && options.Password != "" {
		password = options.Password
	}

	dbIndex := 0
	if options != nil {
		dbIndex, _ = strconv.Atoi(options.Bucket)
	}

	return newRedis(address, password, dbIndex)
}
func newRedis(endpoints []string, password string, dbIndex int) (*RedisStore, error) {
	// TODO: use *redis.ClusterClient if we support miltiple endpoints
	client, err := lgredis.NewSys()
	if err != nil {
		return nil, err
	}
	// Listen to Keyspace events
	client.GetClient().ConfigSet(context.Background(), "notify-keyspace-events", "KEA")

	return &RedisStore{
		redis:  client,
		script: redis.NewScript(luaScript()),
		codec:  defaultCodec{},
	}, nil
}

const (
	noExpiration   = time.Duration(0)
	defaultLockTTL = 60 * time.Second
)

type defaultCodec struct{}

func (this defaultCodec) encode(kv *dcore.KVPair) (string, error) {
	b, err := json.Marshal(kv)
	return string(b), err
}
func (this defaultCodec) decode(b string, kv *dcore.KVPair) error {
	return json.Unmarshal([]byte(b), kv)
}
func regexWatch(key string, withChildren bool) string {
	var regex string
	if withChildren {
		regex = fmt.Sprintf("__keyspace*:%s*", key)
		// for all database and keys with $key prefix
	} else {
		regex = fmt.Sprintf("__keyspace*:%s", key)
		// for all database and keys with $key
	}
	return regex
}

type getter func() (interface{}, error)
type pusher func(interface{})

func watchLoop(msgCh chan *redis.Message, stopCh <-chan struct{}, get getter, push pusher) error {
	// deliver the original data before we setup any events
	pair, err := get()
	if err != nil {
		return err
	}
	push(pair)

	for m := range msgCh {
		// retrieve and send back
		pair, err := get()
		if err != nil && err != store.ErrKeyNotFound {
			return err
		}

		// in case of watching a key that has been expired or deleted return and empty KV
		if err == store.ErrKeyNotFound && (m.Payload == "expire" || m.Payload == "del") {
			push(&dcore.KVPair{})
		} else {
			push(pair)
		}
	}

	return nil
}

type subscribe struct {
	pubsub  *redis.PubSub
	closeCh chan struct{}
}

func newSubscribe(redis lgredis.ISys, regex string) (*subscribe, error) {
	ch := redis.GetClient().PSubscribe(context.Background(), regex)
	return &subscribe{
		pubsub:  ch,
		closeCh: make(chan struct{}),
	}, nil
}
func (this *subscribe) Close() error {
	close(this.closeCh)
	return this.pubsub.Close()
}

func (this *subscribe) Receive(stopCh <-chan struct{}) chan *redis.Message {
	msgCh := make(chan *redis.Message)
	go this.receiveLoop(msgCh, stopCh)
	return msgCh
}

func (this *subscribe) receiveLoop(msgCh chan *redis.Message, stopCh <-chan struct{}) {
	defer close(msgCh)

	for {
		select {
		case <-this.closeCh:
			return
		case <-stopCh:
			return
		default:
			msg, err := this.pubsub.ReceiveMessage(context.Background())
			if err != nil {
				return
			}
			if msg != nil {
				msgCh <- msg
			}
		}
	}
}

type RedisStore struct {
	redis  lgredis.ISys
	script *redis.Script
	codec  defaultCodec
}

func (this *RedisStore) Get(key string) (*dcore.KVPair, error) {
	return this.get(normalize(key))
}
func (this *RedisStore) List(directory string) ([]*dcore.KVPair, error) {
	return this.list(normalize(directory))
}
func (this *RedisStore) Put(key string, value []byte, opts *dcore.WriteOptions) error {
	expirationAfter := noExpiration
	if opts != nil && opts.TTL != 0 {
		expirationAfter = opts.TTL
	}

	return this.setTTL(normalize(key), &dcore.KVPair{
		Key:       key,
		Value:     value,
		LastIndex: sequenceNum(),
	}, expirationAfter)
}
func (this *RedisStore) Watch(key string, stopCh <-chan struct{}) (<-chan *dcore.KVPair, error) {
	watchCh := make(chan *dcore.KVPair)
	nKey := normalize(key)

	get := getter(func() (interface{}, error) {
		pair, err := this.get(nKey)
		if err != nil {
			return nil, err
		}
		return pair, nil
	})

	push := pusher(func(v interface{}) {
		if val, ok := v.(*dcore.KVPair); ok {
			watchCh <- val
		}
	})

	sub, err := newSubscribe(this.redis, regexWatch(nKey, false))
	if err != nil {
		return nil, err
	}

	go func(sub *subscribe, stopCh <-chan struct{}, get getter, push pusher) {
		defer sub.Close()
		msgCh := sub.Receive(stopCh)
		if err := watchLoop(msgCh, stopCh, get, push); err != nil {
			log.Printf("watchLoop in Watch err:%v\n", err)
		}
	}(sub, stopCh, get, push)

	return watchCh, nil
}
func (this *RedisStore) WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*dcore.KVPair, error) {
	watchCh := make(chan []*dcore.KVPair)
	nKey := normalize(directory)

	get := getter(func() (interface{}, error) {
		pair, err := this.list(nKey)
		if err != nil {
			return nil, err
		}
		return pair, nil
	})

	push := pusher(func(v interface{}) {
		if _, ok := v.([]*dcore.KVPair); !ok {
			return
		}
		watchCh <- v.([]*dcore.KVPair)
	})

	sub, err := newSubscribe(this.redis, regexWatch(nKey, true))
	if err != nil {
		return nil, err
	}

	go func(sub *subscribe, stopCh <-chan struct{}, get getter, push pusher) {
		defer sub.Close()

		msgCh := sub.Receive(stopCh)
		if err := watchLoop(msgCh, stopCh, get, push); err != nil {
			log.Printf("watchLoop in WatchTree err:%v\n", err)
		}
	}(sub, stopCh, get, push)

	return watchCh, nil
}
func (this *RedisStore) AtomicDelete(key string, previous *dcore.KVPair) (bool, error) {
	if err := this.cad(normalize(key), previous); err != nil {
		return false, err
	}
	return true, nil
}

func (this *RedisStore) Close() {
	this.redis.GetClient().Close()
}

func (this *RedisStore) get(key string) (*dcore.KVPair, error) {
	reply, err := this.redis.GetClient().Get(context.Background(), key).Bytes()
	if err != nil {
		if err == lgredis.RedisNil {
			return nil, dcore.ErrKeyNotFound
		}
		return nil, err
	}
	val := dcore.KVPair{}
	if err := this.codec.decode(string(reply), &val); err != nil {
		return nil, err
	}
	return &val, nil
}

func (this *RedisStore) list(directory string) ([]*dcore.KVPair, error) {

	var allKeys []string
	regex := scanRegex(directory) // for all keyed with $directory
	allKeys, err := this.keys(regex)
	if err != nil {
		return nil, err
	}
	// TODO: need to handle when #key is too large
	return this.mget(directory, allKeys...)
}
func (this *RedisStore) keys(regex string) ([]string, error) {
	const (
		startCursor  = 0
		endCursor    = 0
		defaultCount = 10
	)

	var allKeys []string

	keys, nextCursor, err := this.redis.GetClient().Scan(context.Background(), startCursor, regex, defaultCount).Result()
	if err != nil {
		return nil, err
	}
	allKeys = append(allKeys, keys...)
	for nextCursor != endCursor {
		keys, nextCursor, err = this.redis.GetClient().Scan(context.Background(), nextCursor, regex, defaultCount).Result()
		if err != nil {
			return nil, err
		}

		allKeys = append(allKeys, keys...)
	}
	if len(allKeys) == 0 {
		return nil, store.ErrKeyNotFound
	}
	return allKeys, nil
}

func (this *RedisStore) cad(key string, old *dcore.KVPair) error {
	oldVal, err := this.codec.encode(old)
	if err != nil {
		return err
	}

	return this.runScript(
		cmdCAD,
		key,
		oldVal,
	)
}

func (this *RedisStore) runScript(args ...interface{}) error {
	err := this.script.Run(
		context.Background(),
		this.redis.GetClient(),
		nil,
		args...,
	).Err()
	if err != nil && strings.Contains(err.Error(), "redis: key is not found") {
		return store.ErrKeyNotFound
	}
	if err != nil && strings.Contains(err.Error(), "redis: value has been changed") {
		return store.ErrKeyModified
	}
	return err
}
func (this *RedisStore) mget(directory string, keys ...string) ([]*dcore.KVPair, error) {
	replies, err := this.redis.GetClient().MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}

	pairs := []*dcore.KVPair{}
	for _, reply := range replies {
		var sreply string
		if _, ok := reply.(string); ok {
			sreply = reply.(string)
		}
		if sreply == "" {
			// empty reply
			continue
		}

		newkv := &dcore.KVPair{}
		if err := this.codec.decode(sreply, newkv); err != nil {
			return nil, err
		}
		if normalize(newkv.Key) != directory {
			pairs = append(pairs, newkv)
		}
	}
	return pairs, nil
}

// Delete the value at the specified key
func (this *RedisStore) Delete(key string) error {
	return this.redis.GetClient().Del(context.Background(), normalize(key)).Err()
}

// Exists verify if a Key exists in the store
func (this *RedisStore) Exists(key string) (bool, error) {
	i, err := this.redis.GetClient().Exists(context.Background(), normalize(key)).Result()
	if err != nil {
		return false, err
	}
	return i == 1, nil
}

func normalize(key string) string {
	key = dcore.Normalize(key)
	return strings.TrimSuffix(key, "/")
}

func sequenceNum() uint64 {
	// TODO: use uuid if we concerns collision probability of this number
	return uint64(time.Now().Nanosecond())
}
func scanRegex(directory string) string {
	return fmt.Sprintf("%s*", directory)
}
func (this *RedisStore) setTTL(key string, val *dcore.KVPair, ttl time.Duration) error {
	valStr, err := this.codec.encode(val)
	if err != nil {
		return err
	}

	return this.redis.GetClient().Set(context.Background(), key, valStr, ttl).Err()
}
