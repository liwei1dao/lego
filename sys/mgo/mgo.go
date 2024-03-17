package mgo

import (
	"context"
	"fmt"

	"github.com/liwei1dao/lego/core"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newSys(options Options) (sys *Mongodb, err error) {
	sys = &Mongodb{options: options}
	err = sys.init()
	return
}

type Mongodb struct {
	options  Options
	Client   *mongo.Client
	Database *mongo.Database
}

func (this *Mongodb) init() (err error) {
	opt := options.Client().ApplyURI(this.options.MongodbUrl)
	opt.SetMaxPoolSize(this.options.MaxPoolSize) //使用最大的连接数
	if client, err := mongo.Connect(this.getContext(), opt); err != nil {
		return fmt.Errorf("连接数据库错误 err=%s", err.Error())
	} else {
		this.Client = client
		if err = client.Ping(this.getContext(), readpref.Primary()); err != nil {
			return fmt.Errorf("数据库不可用 err=%s", err.Error())
		}
		this.Database = client.Database(this.options.MongodbDatabase)
	}
	return
}

func (this *Mongodb) Close() (err error) {
	err = this.Client.Disconnect(this.getContext())
	return
}

func (this *Mongodb) ListCollectionNames(filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	return this.Database.ListCollectionNames(this.getContext(), filter, opts...)
}

func (this *Mongodb) Collection(sqltable core.SqlTable) *mongo.Collection {
	return this.Client.Database(this.options.MongodbDatabase).Collection(string(sqltable))
}

func (this *Mongodb) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.options.TimeOut)
	return
}

func (this *Mongodb) CreateIndex(sqltable core.SqlTable, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	return this.Collection(sqltable).Indexes().CreateOne(this.getContext(), model, opts...)
}

func (this *Mongodb) DeleteIndex(sqltable core.SqlTable, name string, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	return this.Collection(sqltable).Indexes().DropOne(this.getContext(), name, opts...)
}

func (this *Mongodb) UseSession(fn func(sessionContext mongo.SessionContext) error) error {
	return this.Client.UseSession(this.getContext(), fn)
}

func (this *Mongodb) CountDocuments(sqltable core.SqlTable, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return this.Collection(sqltable).CountDocuments(this.getContext(), filter, opts...)
}
func (this *Mongodb) CountDocumentsByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return this.Collection(sqltable).CountDocuments(ctx, filter, opts...)
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
