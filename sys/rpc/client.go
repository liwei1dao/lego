package rpc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/discovery"
	"github.com/liwei1dao/lego/sys/discovery/dcore"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

func NewXClient(service string, selector rpccore.ISelector, discovery discovery.IServiceDiscovery, cpools rpccore.IConnPool) IClient {
	client := &Client{

		selector:  selector,
		discovery: discovery,
		service:   service,
	}

	pairs := discovery.GetServices()
	sort.Slice(pairs, func(i, j int) bool {
		return strings.Compare(pairs[i].Key, pairs[j].Key) <= 0
	})
	servers := make(map[string]string, len(pairs))
	for _, p := range pairs {
		servers[p.Key] = p.Value
	}
	filterByStateAndGroup(servers)

	client.servers = servers

	ch := client.discovery.WatchService()
	if ch != nil {
		client.ch = ch
		go client.watch(ch)
	}

	return client
}

type Client struct {
	service   string
	options   *Options
	mu        sync.RWMutex
	servers   map[string]string
	cpools    rpccore.IConnPool
	discovery discovery.IServiceDiscovery
	selector  rpccore.ISelector
	ch        chan []*dcore.KV
	mutex     sync.Mutex
	seq       uint64
	pending   map[uint64]*MessageCall
}

// 同步执行
func (this *Client) Call(ctx context.Context, req interface{}, reply interface{}) (err error) { //同步调用 等待结果
	seq := new(uint64)
	ctx = rpccore.WithValue(ctx, rpccore.CallSeqKey, seq)
	stime := time.Now()
	// this.options.Log.Debug("Call Start", log.Field{Key: "servicePath", Value: servicePath}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req})
	defer func() {
		this.options.Log.Debug("RPC Call",
			log.Field{Key: "t", Value: time.Since(stime).Milliseconds()},
			log.Field{Key: "servicePath", Value: this.service},
			log.Field{Key: "req", Value: req},
			log.Field{Key: "reply", Value: reply},
		)
	}()
	var call *MessageCall
	call, err = this.call(ctx, req, reply)
	select {
	case <-ctx.Done(): // cancel by context
		this.mutex.Lock()
		call := this.pending[*seq]
		delete(this.pending, *seq)
		this.mutex.Unlock()
		if call != nil {
			call.Error = ctx.Err()
			call.done(this.options.Log)
		}
		return ctx.Err()
	case call := <-call.Done:
		err = call.Error
	}
	return
}

// 异步执行 异步返回
func (this *Client) Go(ctx context.Context, req interface{}, reply interface{}) (call *MessageCall, err error) { //异步调用 异步返回
	seq := new(uint64)
	ctx = rpccore.WithValue(ctx, rpccore.CallSeqKey, seq)
	stime := time.Now()
	// this.options.Log.Debug("Go start!", log.Field{Key: "servicePath", Value: servicePath}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req})
	defer func() {
		this.options.Log.Debug("RPC Go",
			log.Field{Key: "t", Value: time.Since(stime).Milliseconds()},
			log.Field{Key: "service", Value: this.service},
			log.Field{Key: "req", Value: req},
			log.Field{Key: "reply", Value: reply},
		)
	}()
	call, err = this.call(ctx, req, reply)
	return
}

func (this *Client) Broadcast(ctx context.Context, args interface{}) (err error) {
	return
}

func (this *Client) call(ctx context.Context, args interface{}, reply interface{}) (call *MessageCall, err error) {
	call = new(MessageCall)
	call.Service = this.service
	call.Done = make(chan *MessageCall, 10)
	call.Args = args
	call.Reply = reply
	this.mutex.Lock()
	seq := this.seq
	this.seq++
	this.pending[seq] = call
	this.mutex.Unlock()
	if cseq, ok := ctx.Value(rpccore.CallSeqKey).(*uint64); ok {
		*cseq = seq
	}
	var client rpccore.IConnClient
	if client, err = this.getclient(ctx); err != nil {
		return
	}
	err = this.send(client, call, seq)
	return
}

func (this *Client) getclient(ctx context.Context) (client rpccore.IConnClient, err error) {
	nodes := this.selector.Select(ctx)
	if nodes == nil || len(nodes) == 0 {
		err = fmt.Errorf("no found any service:%s", this.service)
		return
	}
	if client, err = this.cpools.GetClient(nodes[0]); err != nil {
		return
	}
	return
}

func (this *Client) watch(ch chan []*dcore.KV) {
	for pairs := range ch {
		sort.Slice(pairs, func(i, j int) bool {
			return strings.Compare(pairs[i].Key, pairs[j].Key) <= 0
		})
		servers := make(map[string]string, len(pairs))
		for _, p := range pairs {
			servers[p.Key] = p.Value
		}
		this.mu.Lock()
		filterByStateAndGroup(servers)
		this.servers = servers

		if this.selector != nil {
			this.selector.UpdateServer(servers)
		}

		this.mu.Unlock()
	}
}

func (this *Client) send(client rpccore.IConnClient, call *MessageCall, seq uint64) (err error) {
	var (
		data    []byte
		allData *[]byte
	)
	req := protocol.GetPooledMsg()
	req.SetVersion(this.options.ProtoVersion)
	req.SetMessageType(rpccore.Request)
	req.SetSeq(seq)
	if call.Reply == nil {
		req.SetOneway(true)
	}
	if call.Metadata != nil {
		req.SetMetadata(call.Metadata)
	}
	req.SetFrom(this.options.ServiceNode)
	req.SetSerializeType(this.options.SerializeType)
	codec := codecs[this.options.SerializeType]
	if codec == nil {
		err = rpccore.ErrUnsupportedCodec
		return
	}
	data, err = codec.Marshal(call.Args)
	if err != nil {
		return
	}
	if len(data) > 1024 && this.options.CompressType != rpccore.CompressNone {
		req.SetCompressType(this.options.CompressType)
	}
	req.SetPayload(data)
	allData = req.EncodeSlicePointer()
	defer func() {
		protocol.PutData(allData)
		protocol.FreeMsg(req)
	}()
	err = client.Write(*allData)
	return
}

func (this *Client) handleresponse(ctx context.Context, res rpccore.IMessage) {
	var call *MessageCall
	seq := res.Seq()
	isServerMessage := (res.MessageType() == rpccore.Request && !res.IsHeartbeat() && res.IsOneway())
	if !isServerMessage {
		this.mutex.Lock()
		call = this.pending[seq]
		delete(this.pending, seq)
		this.mutex.Unlock()
	}
	switch {
	case call == nil:
		this.options.Log.Warnf("call is nil res:%v", res)
	case res.MessageStatusType() == rpccore.Error:
		if len(res.Metadata()) > 0 {
			call.ResMetadata = res.Metadata()
			call.Error = errors.New(res.Metadata()[rpccore.ServiceError])
		}
		if len(res.Payload()) > 0 {
			data := res.Payload()
			codec := codecs[res.SerializeType()]
			if codec != nil {
				_ = codec.Unmarshal(data, call.Reply)
			}
			call.done(this.options.Log)
		}
	default:
		data := res.Payload()
		if len(data) > 0 {
			codec := codecs[res.SerializeType()]
			if codec == nil {
				call.Error = rpccore.ErrUnsupportedCodec
			} else {
				if call.Reply == nil {
					call.Error = fmt.Errorf("%s reply is null no cant Unmarshal !", res.GetService())
				} else {
					err := codec.Unmarshal(data, call.Reply)
					if err != nil {
						call.Error = err
					}
				}
			}
		}
		if len(res.Metadata()) > 0 {
			call.ResMetadata = res.Metadata()
		}
		call.done(this.options.Log)
	}
}

// 获取请求消息对象
func (this *Client) getMessage(service string, args interface{}, reply interface{}) (call *MessageCall, req *protocol.Message, err error) {
	var data []byte
	call = new(MessageCall)
	call.Service = service
	call.Done = make(chan *MessageCall, 10)
	call.Args = args
	call.Reply = reply
	req = protocol.GetPooledMsg()
	req.SetVersion(this.options.ProtoVersion)
	if call.Reply != nil {
		this.mutex.Lock()
		seq := this.seq
		this.seq++
		this.pending[seq] = call
		this.mutex.Unlock()
		req.SetSeq(seq)
		req.SetOneway(true)
	} else {
		req.SetOneway(false)
	}

	req.SetService(call.Service)
	req.SetFrom(this.options.ServiceNode)
	req.SetMessageType(rpccore.Request)
	req.SetSerializeType(this.options.SerializeType)
	data, err = codecs[this.options.SerializeType].Marshal(call.Args)
	if err != nil {
		return
	}
	if len(data) > 1024 && this.options.CompressType != rpccore.CompressNone {
		req.SetCompressType(this.options.CompressType)
	}
	req.SetPayload(data)
	return
}

func filterByStateAndGroup(servers map[string]string) {
	for k, v := range servers {
		if values, err := url.ParseQuery(v); err == nil {
			if state := values.Get("state"); state == "inactive" {
				delete(servers, k)
			}
			// if group != "" && group != values.Get("group") {
			// 	delete(servers, k)
			// }
		}
	}
}
