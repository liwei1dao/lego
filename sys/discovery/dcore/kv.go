package dcore

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rpcxio/libkv/store"
)

// Initialize creates a new Store object, initializing the client
type Initialize func(addrs []string, options *Config) (IStore, error)

var (
	// Backend initializers
	initializers = make(map[Backend]Initialize)

	supportedBackend = func() string {
		keys := make([]string, 0, len(initializers))
		for k := range initializers {
			keys = append(keys, string(k))
		}
		sort.Strings(keys)
		return strings.Join(keys, ", ")
	}()
)

// NewStore creates an instance of store
func NewStore(backend Backend, addrs []string, options *Config) (IStore, error) {
	if init, exists := initializers[backend]; exists {
		return init(addrs, options)
	}

	return nil, fmt.Errorf("%s %s", store.ErrBackendNotSupported.Error(), supportedBackend)
}

// AddStore adds a new store backend to libkv
func AddStore(store Backend, init Initialize) {
	initializers[store] = init
}
