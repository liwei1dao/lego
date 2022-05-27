package rpcx

import (
	"context"

	lgcore "github.com/liwei1dao/lego/core"
	"github.com/smallnest/rpcx/client"
)

func newClient(addr string) (c *Client, err error) {
	c = &Client{}
	return
}

type Client struct {
	client client.XClient
}

func (this *Client) Stop() (err error) {
	err = this.client.Close()
	return
}

func (this *Client) Call(ctx context.Context, serviceMethod lgcore.Rpc_Key, args interface{}, reply interface{}) (err error) {
	err = this.client.Call(ctx, string(serviceMethod), args, reply)
	return
}
func (this *Client) Go(ctx context.Context, serviceMethod lgcore.Rpc_Key, args interface{}, reply interface{}, done chan *client.Call) (*client.Call, error) {
	return this.client.Go(ctx, string(serviceMethod), args, reply, done)
}
