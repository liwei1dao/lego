package rpcx

import (
	"context"

	"github.com/smallnest/rpcx/client"
)

func newClient(addr string, sId string) (c *Client, err error) {
	c = &Client{}
	d, err := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	c.xclient = client.NewXClient(sId, client.Failfast, client.RandomSelect, d, client.DefaultOption)
	return
}

type Client struct {
	xclient client.XClient
}

func (this *Client) Stop() (err error) {
	err = this.xclient.Close()
	return
}

func (this *Client) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) (err error) {
	err = this.xclient.Call(ctx, string(serviceMethod), args, reply)
	return
}
func (this *Client) Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (*client.Call, error) {
	return this.xclient.Go(ctx, string(serviceMethod), args, reply, done)
}
