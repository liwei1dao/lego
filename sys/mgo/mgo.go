package mgo

import (
	"context"
	"fmt"
	"time"

	"github.com/liwei1dao/lego/core"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func newMongodb(opt ...Option) (*Mongodb, error) {
	mogodb := new(Mongodb)
	mogodb.opts = newOptions(opt...)
	err := mogodb.init()
	return mogodb, err
}

type Mongodb struct {
	opts     Options
	Client   *mongo.Client
	Database *mongo.Database
}

func (this *Mongodb) init() (err error) {
	want, err := readpref.New(readpref.SecondaryMode) //表示只使用辅助节点
	if err != nil {
		return fmt.Errorf("数据库设置辅助节点 err=%s", err.Error())
	}
	wc := writeconcern.New(writeconcern.WMajority())
	readconcern.Majority()
	//链接mongo服务
	opt := options.Client().ApplyURI(this.opts.MongodbUrl)
	opt.SetLocalThreshold(3 * time.Second)     //只使用与mongo操作耗时小于3秒的
	opt.SetMaxConnIdleTime(5 * time.Second)    //指定连接可以保持空闲的最大毫秒数
	opt.SetMaxPoolSize(this.opts.MaxPoolSize)  //使用最大的连接数
	opt.SetReadPreference(want)                //表示只使用辅助节点
	opt.SetReadConcern(readconcern.Majority()) //指定查询应返回实例的最新数据确认为，已写入副本集中的大多数成员
	opt.SetWriteConcern(wc)                    //请求确认写操作传播到大多数mongod实例
	if client, err := mongo.Connect(this.getContext(), opt); err != nil {
		return fmt.Errorf("连接数据库错误 err=%s", err.Error())
	} else {
		this.Client = client
		if err = client.Ping(this.getContext(), readpref.Primary()); err != nil {
			return fmt.Errorf("数据库不可用 err=%s", err.Error())
		}
		this.Database = client.Database(this.opts.MongodbDatabase)
	}
	return
}

func (this *Mongodb) Collection(sqltable core.SqlTable) *mongo.Collection {
	return this.Client.Database(this.opts.MongodbDatabase).Collection(string(sqltable))
}

func (this *Mongodb) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.opts.TimeOut)
	return
}

func (this *Mongodb) CreateIndex(sqltable core.SqlTable, keys interface{}, options *options.IndexOptions) (string, error) {
	return this.Collection(sqltable).Indexes().CreateOne(this.getContext(), mongo.IndexModel{Keys: keys, Options: options})
}

func (this *Mongodb) CountDocuments(sqltable core.SqlTable, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return this.Collection(sqltable).CountDocuments(this.getContext(), filter, opts...)
}

func (this *Mongodb) Find(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return this.Collection(sqltable).Find(this.getContext(), filter, opts...)
}

func (this *Mongodb) FindByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return this.Collection(sqltable).Find(ctx, filter, opts...)
}

func (this *Mongodb) FindOne(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return this.Collection(sqltable).FindOne(this.getContext(), filter, opts...)
}

func (this *Mongodb) FindOneAndUpdate(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return this.Collection(sqltable).FindOneAndUpdate(this.getContext(), filter, update, opts...)
}

func (this *Mongodb) InsertOne(sqltable core.SqlTable, data interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return this.Collection(sqltable).InsertOne(this.getContext(), data, opts...)
}

func (this *Mongodb) InsertMany(sqltable core.SqlTable, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return this.Collection(sqltable).InsertMany(this.getContext(), data, opts...)
}

func (this *Mongodb) InsertManyByCtx(sqltable core.SqlTable, ctx context.Context, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return this.Collection(sqltable).InsertMany(ctx, data, opts...)
}

func (this *Mongodb) UpdateOne(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return this.Collection(sqltable).UpdateOne(this.getContext(), filter, update, opts...)
}

func (this *Mongodb) UpdateOneByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return this.Collection(sqltable).UpdateOne(ctx, filter, update, opts...)
}

func (this *Mongodb) UpdateMany(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return this.Collection(sqltable).UpdateMany(this.getContext(), filter, update, opts...)
}

func (this *Mongodb) UpdateManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return this.Collection(sqltable).UpdateMany(ctx, filter, update, opts...)
}

func (this *Mongodb) FindOneAndDelete(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return this.Collection(sqltable).FindOneAndDelete(this.getContext(), filter, opts...)
}

func (this *Mongodb) DeleteOne(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return this.Collection(sqltable).DeleteOne(this.getContext(), filter, opts...)
}

func (this *Mongodb) DeleteMany(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return this.Collection(sqltable).DeleteMany(this.getContext(), filter, opts...)
}

func (this *Mongodb) DeleteManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return this.Collection(sqltable).DeleteMany(ctx, filter, opts...)
}

func (this *Mongodb) Aggregate(sqltable core.SqlTable, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return this.Collection(sqltable).Aggregate(this.getContext(), pipeline, opts...)
}
