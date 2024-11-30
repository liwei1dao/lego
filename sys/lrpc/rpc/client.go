package rpc

import (
	"context"
	"sync"

	"github.com/liwei1dao/lego/sys/lrpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

type Client struct {
	mu      sync.Mutex
	seq     uint64
	pending map[uint64]*Call
}

func (this *Client) call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *Call, err error) {
	call = new(Call)
	call.ServicePath = servicePath
	call.ServiceMethod = serviceMethod
	call.Done = make(chan *Call, 10)
	call.Args = args
	call.Reply = reply
	this.mu.Lock()
	seq := this.seq
	this.seq++
	this.pending[seq] = call
	this.mu.Unlock()
	if cseq, ok := ctx.Value(ContextCallSeqKey).(*uint64); ok {
		*cseq = seq
	}
	var client rpccore.IConnClient
	if client, err = this.getclient(ctx, servicePath); err != nil {
		return
	}
	err = this.send(client, call, seq)
	return
}

func (this *Client) send(client rpccore.IConnClient, call *Call, seq uint64) (err error) {
	var (
		data    []byte
		allData *[]byte
	)

	req := protocol.GetPooledMsg()
	req.SetMessageType(protocol.Request)
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
