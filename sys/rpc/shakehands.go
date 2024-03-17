package rpc

import (
	"context"
	"fmt"

	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

type ShakehandsReq struct {
}

type ShakehandsResp struct {
}

//握手请求
func (this *rpc) ShakehandsRequest(ctx context.Context, client rpccore.IConnClient) (err error) {
	var (
		call *MessageCall
		req  *protocol.Message
	)
	if call, req, err = this.getMessage("Shakehands", &ShakehandsReq{}, &ShakehandsResp{}); err != nil {
		return
	}
	req.SetShakeHands(true)
	allData := req.EncodeSlicePointer()
	err = client.Write(*allData)
	protocol.PutData(allData)
	protocol.FreeMsg(req)
	if err != nil {
		return
	}
	select {
	case <-ctx.Done(): // cancel by context
		this.pendingmutex.Lock()
		call := this.pending[req.Seq()]
		delete(this.pending, req.Seq())
		this.pendingmutex.Unlock()
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

//握手请求 回应
func (this *rpc) ShakehandsResponse(ctx context.Context, client rpccore.IConnClient, req rpccore.IMessage) (err error) {
	res := req.Clone()
	defer func() {
		data := res.EncodeSlicePointer()
		client.Write(*data)
		protocol.PutData(data)
	}()
	res.SetMessageType(rpccore.Response)
	codec := codecs[req.SerializeType()]
	if codec == nil {
		err = fmt.Errorf("can not find codec for %d", req.SerializeType())
		handleError(res, err)
		return
	}
	data, err := codec.Marshal(&ShakehandsResp{})
	if err != nil {
		handleError(res, err)
		return
	}
	res.SetPayload(data)
	if len(res.Payload()) > 1024 && res.CompressType() != rpccore.CompressNone {
		res.SetCompressType(res.CompressType())
	}
	this.cpool.AddClient(client, req.From())
	return
}
