package rpc

import "github.com/smallnest/rpcx/errors"

// 插件容器
type PluginContainer interface {
	Add(plugin Plugin)
	Remove(plugin Plugin)
	All() []Plugin

	DoRegister(name string, fn interface{}, metadata string) error
	DoUnregister(name string) error
}

// 插件集合
type Plugin interface{}

type (
	RegisterPlugin interface {
		Register(name string, fn interface{}, metadata string) error
		Unregister(name string) error
	}
	RegisterFunctionPlugin interface {
		RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error
	}
)

type pluginContainer struct {
	plugins []Plugin
}

func (p *pluginContainer) Add(plugin Plugin) {
	p.plugins = append(p.plugins, plugin)
}

func (p *pluginContainer) Remove(plugin Plugin) {
	if p.plugins == nil {
		return
	}
	plugins := make([]Plugin, 0, len(p.plugins))
	for _, p := range p.plugins {
		if p != plugin {
			plugins = append(plugins, p)
		}
	}
	p.plugins = plugins
}

func (p *pluginContainer) All() []Plugin {
	return p.plugins
}

// DoRegister invokes DoRegister plugin.
func (p *pluginContainer) DoRegister(name string, fn interface{}, metadata string) error {
	var es []error
	for _, rp := range p.plugins {
		if plugin, ok := rp.(RegisterPlugin); ok {
			err := plugin.Register(name, fn, metadata)
			if err != nil {
				es = append(es, err)
			}
		}
	}

	if len(es) > 0 {
		return errors.NewMultiError(es)
	}
	return nil
}

// DoUnregister invokes RegisterPlugin.
func (p *pluginContainer) DoUnregister(name string) error {
	var es []error
	for _, rp := range p.plugins {
		if plugin, ok := rp.(RegisterPlugin); ok {
			err := plugin.Unregister(name)
			if err != nil {
				es = append(es, err)
			}
		}
	}

	if len(es) > 0 {
		return errors.NewMultiError(es)
	}
	return nil
}
