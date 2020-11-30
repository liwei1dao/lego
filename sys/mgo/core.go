package mgo

import (
	"context"

	"github.com/liwei1dao/lego/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IMongodb interface {
		CountDocuments(sqltable core.SqlTable, filter interface{}, opts ...*options.CountOptions) (int64, error)
		Find(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
		FindOne(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
		FindOneAndUpdate(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
		InsertOne(sqltable core.SqlTable, data interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
		InsertMany(sqltable core.SqlTable, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
		InsertManyByCtx(sqltable core.SqlTable, ctx context.Context, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
		UpdateOne(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		UpdateMany(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		UpdateManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		FindOneAndDelete(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
		DeleteOne(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
		DeleteMany(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
		DeleteManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	}
)

var (
	defsys IMongodb
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IMongodb, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func CountDocuments(sqltable core.SqlTable, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return defsys.CountDocuments(sqltable, filter, opts...)
}

func Find(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return defsys.Find(sqltable, filter, opts...)
}

func FindOne(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return defsys.FindOne(sqltable, filter, opts...)
}

func FindOneAndUpdate(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return defsys.FindOneAndUpdate(sqltable, filter, update, opts...)
}

func InsertOne(sqltable core.SqlTable, data interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return defsys.InsertOne(sqltable, data, opts...)
}

func InsertMany(sqltable core.SqlTable, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return defsys.InsertMany(sqltable, data, opts...)
}

func InsertManyByCtx(sqltable core.SqlTable, ctx context.Context, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return defsys.InsertManyByCtx(sqltable, ctx, data, opts...)
}

func UpdateOne(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return defsys.UpdateOne(sqltable, filter, update, opts...)
}

func UpdateMany(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return defsys.UpdateMany(sqltable, filter, update, opts...)
}

func UpdateManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return defsys.UpdateManyByCtx(sqltable, ctx, filter, update, opts...)
}

func FindOneAndDelete(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return defsys.FindOneAndDelete(sqltable, filter, opts...)
}

func DeleteOne(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return defsys.DeleteOne(sqltable, filter, opts...)
}

func DeleteMany(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return defsys.DeleteMany(sqltable, filter, opts...)
}

func DeleteManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return defsys.DeleteManyByCtx(sqltable, ctx, filter, opts...)
}
