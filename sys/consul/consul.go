package consul

import (
	"github.com/hashicorp/consul/api"
	consul "github.com/hashicorp/consul/api"
)

func newSys(options Options) (sys *Consul, err error) {
	sys = &Consul{options: options}
	err = sys.init()
	return
}

type Consul struct {
	options Options
	client  *consul.Client
}

func (this *Consul) init() (err error) {
	if this.client, err = api.NewClient(api.DefaultConfig()); err != nil {
		return
	}
	return
}

func (this *Consul) Services() (map[string]*consul.AgentService, error) {
	return this.client.Agent().Services()
}

func (this *Consul) ServicesWithFilter(filter string) (map[string]*consul.AgentService, error) {
	return this.client.Agent().ServicesWithFilter(filter)
}

func (this *Consul) ServiceRegister(service *consul.AgentServiceRegistration) error {
	return this.client.Agent().ServiceRegister(service)
}

func (this *Consul) ServiceDeregister(serviceID string) error {
	return this.client.Agent().ServiceDeregister(serviceID)
}

func (this *Consul) LockKey(key string) (*consul.Lock, error) {
	return this.client.LockKey(key)
}

func (this *Consul) Get(key string, q *consul.QueryOptions) (*consul.KVPair, *consul.QueryMeta, error) {
	return this.client.KV().Get(key, q)
}

func (this *Consul) Put(p *consul.KVPair, q *consul.WriteOptions) (*consul.WriteMeta, error) {
	return this.client.KV().Put(p, q)
}
