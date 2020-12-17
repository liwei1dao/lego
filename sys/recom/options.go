package recom

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	RecomModel   RecomModel                    //推荐模型
	ItemIdsScore map[uint32]map[uint32]float64 //物品评分模型
	ItemIds      []uint32
	UserIds      []uint32
	Ratings      []float64
}

func SetRecomModel(v RecomModel) Option {
	return func(o *Options) {
		o.RecomModel = v
	}
}

func SetItemIdsScore(v map[uint32]map[uint32]float64) Option {
	return func(o *Options) {
		o.ItemIdsScore = v
	}
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
func SetRatings(v []float64) Option {
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
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return options
}
