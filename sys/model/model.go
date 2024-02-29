package model

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/mgo"
	"github.com/liwei1dao/lego/sys/redis"
)

func newSys(options *Options) (sys *Model, err error) {
	sys = &Model{
		options: options,
	}
	err = sys.init()
	return
}

type Model struct {
	options *Options
	Redis   redis.ISys
	Mgo     mgo.ISys
}

func (this *Model) init() (err error) {
	if this.options.RedisIsCluster {
		this.Redis, err = redis.NewSys(
			redis.SetRedisType(redis.Redis_Cluster),
			redis.SetRedis_Cluster_Addr(this.options.RedisAddr),
			redis.SetRedis_Cluster_Password(this.options.RedisPassword))
	} else {
		this.Redis, err = redis.NewSys(
			redis.SetRedisType(redis.Redis_Single),
			redis.SetRedis_Single_Addr(this.options.RedisAddr[0]),
			redis.SetRedis_Single_Password(this.options.RedisPassword),
			redis.SetRedis_Single_DB(this.options.RedisDB),
		)
	}
	if err != nil {
		this.options.Log.Error(err.Error(), log.Field{Key: "options", Value: this.options})
		return
	}
	if this.Mgo, err = mgo.NewSys(
		mgo.SetMongodbUrl(this.options.MongodbUrl),
		mgo.SetMongodbDatabase(this.options.MongodbDatabase),
	); err != nil {
		this.options.Log.Error(err.Error(), log.Field{Key: "options", Value: this.options})
		return
	}
	return
}
