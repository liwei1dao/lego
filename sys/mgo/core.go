package mgo

import (
	"context"

	"github.com/liwei1dao/lego/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defmgodb IMongodb
)

type (
	IMongodb interface {
		CountDocuments(sqltable core.SqlTable, filter interface{}, opts ...*options.CountOptions) (int64, error)
		CreateIndex(sqltable core.SqlTable, keys interface{}, options *options.IndexOptions) (string, error)
		Find(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
		FindByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
		FindOne(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
		FindOneAndUpdate(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
		InsertOne(sqltable core.SqlTable, data interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
		InsertMany(sqltable core.SqlTable, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
		InsertManyByCtx(sqltable core.SqlTable, ctx context.Context, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
		UpdateOne(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		UpdateOneByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		UpdateMany(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		UpdateManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		FindOneAndDelete(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
		DeleteOne(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
		DeleteMany(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
		DeleteManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	}
)

func OnInit(s core.IService, opt ...Option) (err error) {
	defmgodb, err = newMongodb(opt...)
	return
}

func NewMongodb(opt ...Option) (mgodb IMongodb, err error) {
	mgodb, err = newMongodb(opt...)
	return
}

func CountDocuments(sqltable core.SqlTable, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return defmgodb.CountDocuments(sqltable, filter, opts...)
}

func Find(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return defmgodb.Find(sqltable, filter, opts...)
}

func FindByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return defmgodb.FindByCtx(sqltable, ctx, filter, opts...)
}

func FindOne(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return defmgodb.FindOne(sqltable, filter, opts...)
}

func FindOneAndUpdate(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return defmgodb.FindOneAndUpdate(sqltable, filter, update, opts...)
}

func InsertOne(sqltable core.SqlTable, data interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return defmgodb.InsertOne(sqltable, data, opts...)
}

func InsertMany(sqltable core.SqlTable, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return defmgodb.InsertMany(sqltable, data, opts...)
}

func InsertManyByCtx(sqltable core.SqlTable, ctx context.Context, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return defmgodb.InsertManyByCtx(sqltable, ctx, data, opts...)
}

func UpdateOne(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return defmgodb.UpdateOne(sqltable, filter, update, opts...)
}

func UpdateOneByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return defmgodb.UpdateOneByCtx(sqltable, ctx, filter, update, opts...)
}

func UpdateMany(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return defmgodb.UpdateMany(sqltable, filter, update, opts...)
}

func UpdateManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return defmgodb.UpdateManyByCtx(sqltable, ctx, filter, update, opts...)
}

func FindOneAndDelete(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return defmgodb.FindOneAndDelete(sqltable, filter, opts...)
}

func DeleteOne(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return defmgodb.DeleteOne(sqltable, filter, opts...)
}

func DeleteMany(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return defmgodb.DeleteMany(sqltable, filter, opts...)
}

func DeleteManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return defmgodb.DeleteManyByCtx(sqltable, ctx, filter, opts...)
}

func CreateIndex(sqltable core.SqlTable, keys interface{}, options *options.IndexOptions) (string, error) {
	return defmgodb.CreateIndex(sqltable, keys, options)
}
