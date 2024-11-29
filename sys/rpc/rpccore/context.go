package rpccore

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

const (
	ServiceSeqKey   = "__seqKey__"       //请求序号
	ServiceError    = "__rpcx_error__"   //服务错误信息字段
	ServerTimeout   = "__ServerTimeout"  //服务超时字段
	ReqMetaDataKey  = "__req_metadata"   //请求元数据字段
	ResMetaDataKey  = "__res_metadata"   //返回元数据字段
	ServiceAddrKey  = "__service_addr__" //服务端地址
	CallSeqKey      = "__call_seq__"     //客户端请求id存储key
	ContextTagsLock = "_tagsLock"
	isShareContext  = "_isShareContext"
)

type Context struct {
	tagsLock *sync.Mutex
	tags     map[interface{}]interface{}
	context.Context
}

func NewContext(ctx context.Context) *Context {
	tagsLock := &sync.Mutex{}
	ctx = context.WithValue(ctx, ContextTagsLock, tagsLock)
	return &Context{
		tagsLock: tagsLock,
		Context:  ctx,
		tags:     map[interface{}]interface{}{},
	}
}

func (c *Context) Lock() {
	c.tagsLock.Lock()
}

func (c *Context) Unlock() {
	c.tagsLock.Unlock()
}

func (c *Context) Value(key interface{}) interface{} {
	c.tagsLock.Lock()
	defer c.tagsLock.Unlock()
	if c.tags == nil {
		c.tags = make(map[interface{}]interface{})
	}

	if v, ok := c.tags[key]; ok {
		return v
	}
	return c.Context.Value(key)
}

func (c *Context) SetValue(key, val interface{}) {
	c.tagsLock.Lock()
	defer c.tagsLock.Unlock()

	if c.tags == nil {
		c.tags = make(map[interface{}]interface{})
	}
	c.tags[key] = val
}

// DeleteKey delete the kv pair by key.
func (c *Context) DeleteKey(key interface{}) {
	c.tagsLock.Lock()
	defer c.tagsLock.Unlock()

	if c.tags == nil || key == nil {
		return
	}
	delete(c.tags, key)
}

func (c *Context) String() string {
	return fmt.Sprintf("%v.WithValue(%v)", c.Context, c.tags)
}

func WithValue(parent context.Context, key, val interface{}) *Context {
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	tags := make(map[interface{}]interface{})
	tags[key] = val
	return &Context{Context: parent, tags: tags, tagsLock: &sync.Mutex{}}
}

func WithLocalValue(ctx *Context, key, val interface{}) *Context {
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	if ctx.tags == nil {
		ctx.tags = make(map[interface{}]interface{})
	}

	ctx.tags[key] = val
	return ctx
}

func IsShareContext(ctx context.Context) bool {
	ok := ctx.Value(isShareContext)
	return ok != nil
}
