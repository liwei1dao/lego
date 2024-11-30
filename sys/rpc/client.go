package rpc

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
	"github.com/smallnest/rpcx/share"
)

var (
	ErrShutdown = errors.New("connection is shut down")
)

type seqKey struct{}
type Client struct {
	options      *Options
	Conn         rpccore.IConnClient
	node         core.IServiceNode
	mutex        sync.Mutex
	seq          uint64
	pending      map[uint64]*MessageCall
	closing      bool
	shutdown     bool
	pluginClosed bool
}

func (this *Client) ServiceNode() core.IServiceNode {
	return this.options.ServiceNode
}
func (client *Client) IsClosing() bool {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return client.closing
}

// IsShutdown client is shutdown or not.
func (client *Client) IsShutdown() bool {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return client.shutdown
}
func (client *Client) Call(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) error {
	return client.call(ctx, servicePath, serviceMethod, args, reply)
}

func (client *Client) Go(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *MessageCall) *MessageCall {
	call := new(MessageCall)
	call.ServicePath = servicePath
	call.ServiceMethod = serviceMethod
	meta := ctx.Value(rpccore.ReqMetaDataKey)
	if meta != nil {
		call.Metadata = meta.(map[string]string)
	}

	if !rpccore.IsShareContext(ctx) {
		ctx = rpccore.NewContext(ctx)
	}

	call.Args = args
	call.Reply = reply
	if done == nil {
		done = make(chan *MessageCall, 10)
	} else {
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done

	if share.Trace {
		log.Debugf("client.Go send request for %s.%s, args: %+v in case of client call", servicePath, serviceMethod, args)
	}

	go client.send(ctx, call)

	return call
}
func (client *Client) call(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) error {
	seq := new(uint64)
	ctx = context.WithValue(ctx, seqKey{}, seq)

	if share.Trace {
		log.Debugf("client.call for %s.%s, args: %+v in case of client call", servicePath, serviceMethod, args)
		defer func() {
			log.Debugf("client.call done for %s.%s, args: %+v in case of client call", servicePath, serviceMethod, args)
		}()
	}

	Done := client.Go(ctx, servicePath, serviceMethod, args, reply, make(chan *MessageCall, 1)).Done

	var err error
	select {
	case <-ctx.Done(): // cancel by context
		client.mutex.Lock()
		call := client.pending[*seq]
		delete(client.pending, *seq)
		client.mutex.Unlock()
		if call != nil {
			call.Error = ctx.Err()
			call.done()
		}

		return ctx.Err()
	case call := <-Done:
		err = call.Error
		meta := ctx.Value(share.ResMetaDataKey)
		if meta != nil && len(call.ResMetadata) > 0 {
			resMeta := meta.(map[string]string)
			locker, ok := ctx.Value(share.ContextTagsLock).(*sync.Mutex)
			if ok {

				locker.Lock()
				for k, v := range call.ResMetadata {
					resMeta[k] = v
				}
				resMeta[share.ServerAddress] = client.Conn.ServiceNode().Addr()
				locker.Unlock()

			} else {
				for k, v := range call.ResMetadata {
					resMeta[k] = v
				}
				resMeta[share.ServerAddress] = client.Conn.ServiceNode().Addr()
			}
		}
	}

	return err
}

func (this *Client) send(ctx context.Context, call *MessageCall) (err error) {
	var (
		data    []byte
		allData *[]byte
	)
	this.mutex.Lock()
	if this.shutdown || this.closing {
		call.Error = ErrShutdown
		this.mutex.Unlock()
		call.done()
		return
	}

	isHeartbeat := call.ServicePath == "" && call.ServiceMethod == ""
	serializeType := this.options.SerializeType
	if isHeartbeat {
		serializeType = rpccore.MsgPack
	}
	codec := codecs[serializeType]
	if codec == nil {
		call.Error = rpccore.ErrUnsupportedCodec
		this.mutex.Unlock()
		call.done()
		return
	}

	if this.pending == nil {
		this.pending = make(map[uint64]*MessageCall)
	}

	seq := this.seq
	this.seq++
	this.pending[seq] = call
	this.mutex.Unlock()

	if cseq, ok := ctx.Value(rpccore.ServiceSeqKey).(*uint64); ok {
		*cseq = seq
	}
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

	req.SetServiceMethod(call.ServiceMethod)
	req.SetFrom(this.options.ServiceNode)
	req.SetSerializeType(this.options.SerializeType)
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
	err = this.Conn.Write(*allData)
	return
}

func (client *Client) heartbeat() {
	t := time.NewTicker(time.Second * time.Duration(client.options.HeartbeatInterval))

	if client.options.MaxWaitForHeartbeat == 0 {
		client.options.MaxWaitForHeartbeat = 30
	}

	for range t.C {
		if client.IsShutdown() || client.IsClosing() {
			t.Stop()
			return
		}

		request := time.Now().UnixNano()
		reply := int64(0)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.options.MaxWaitForHeartbeat))
		err := client.Call(ctx, "", "", &request, &reply)
		abnormal := false
		if ctx.Err() != nil {
			// log.Warnf("failed to heartbeat to %s, context err: %v", client.Conn.RemoteAddr().String(), ctx.Err())
			abnormal = true
		}
		cancel()
		if err != nil {
			// log.Warnf("failed to heartbeat to %s: %v", client.Conn.RemoteAddr().String(), err)
			abnormal = true
		}

		if reply != request {
			// log.Warnf("reply %d in heartbeat to %s is different from request %d", reply, client.Conn.RemoteAddr().String(), request)
		}

		if abnormal {
			client.Close()
		}
	}
}

func (client *Client) Close() error {
	client.mutex.Lock()

	for seq, call := range client.pending {
		delete(client.pending, seq)
		if call != nil {
			call.Error = ErrShutdown
			call.done()
		}
	}

	var err error
	if !client.pluginClosed {
		client.pluginClosed = true
		err = client.Conn.Close()
	}

	if client.closing || client.shutdown {
		client.mutex.Unlock()
		return ErrShutdown
	}

	client.closing = true
	client.mutex.Unlock()
	return err
}
