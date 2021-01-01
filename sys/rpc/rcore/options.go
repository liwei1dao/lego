package rcore

import "github.com/nats-io/nats.go"

type sOption func(*sOptions)
type sOptions struct {
	sId          string
	Nats         *nats.Conn
	MaxCoroutine int
}

func SId(v string) sOption {
	return func(o *sOptions) {
		o.sId = v
	}
}
func SNats(v *nats.Conn) sOption {
	return func(o *sOptions) {
		o.Nats = v
	}
}
func SMaxCoroutine(v int) sOption {
	return func(o *sOptions) {
		o.MaxCoroutine = v
	}
}

func newSOptions(opts ...sOption) sOptions {
	opt := sOptions{}
	for _, o := range opts {
		o(&opt)
	}

	return opt
}

type cOption func(*cOptions)
type cOptions struct {
	sId        string
	rpcId      string
	Nats       *nats.Conn
	RpcExpired int
	Log        bool
}

func CId(v string) cOption {
	return func(o *cOptions) {
		o.sId = v
	}
}
func CrpcId(v string) cOption {
	return func(o *cOptions) {
		o.rpcId = v
	}
}
func CNats(v *nats.Conn) cOption {
	return func(o *cOptions) {
		o.Nats = v
	}
}
func CRpcExpired(v int) cOption {
	return func(o *cOptions) {
		o.RpcExpired = v
	}
}

func newCOptions(opts ...cOption) cOptions {
	opt := cOptions{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
