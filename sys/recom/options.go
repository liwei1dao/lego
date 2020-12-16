package recom

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	UserIds []uint32
	ItemIds []uint32
	Ratings []uint8
}

func SetUserIds(v []uint32) Option {
	return func(o *Options) {
		o.UserIds = v
	}
}
func SetItemIds(v []uint32) Option {
	return func(o *Options) {
		o.ItemIds = v
	}
}
func SetRatings(v []uint8) Option {
	return func(o *Options) {
		o.Ratings = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		UserIds: make([]uint32, 0),
		ItemIds: make([]uint32, 0),
		Ratings: make([]uint8, 0),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
