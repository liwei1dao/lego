package mgo

import (
	"context"

	"github.com/liwei1dao/lego/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	service core.IService
	mgodb   *Mongodb
)

func OnInit(s core.IService, opt ...Option) (err error) {
	service = s
	mgodb, err = newMongodb(opt...)
	return
}

func GetService() core.IService {
	return service
}

func CountDocuments(sqltable core.SqlTable, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return mgodb.Collection(sqltable).CountDocuments(mgodb.getContext(), filter, opts...)
}

func Find(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return mgodb.Collection(sqltable).Find(mgodb.getContext(), filter, opts...)
}

func FindOne(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return mgodb.Collection(sqltable).FindOne(mgodb.getContext(), filter, opts...)
}

func FindOneAndUpdate(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return mgodb.Collection(sqltable).FindOneAndUpdate(mgodb.getContext(), filter, update, opts...)
}

func InsertOne(sqltable core.SqlTable, data interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return mgodb.Collection(sqltable).InsertOne(mgodb.getContext(), data, opts...)
}

func InsertMany(sqltable core.SqlTable, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return mgodb.Collection(sqltable).InsertMany(mgodb.getContext(), data, opts...)
}

func InsertManyByCtx(sqltable core.SqlTable, ctx context.Context, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return mgodb.Collection(sqltable).InsertMany(ctx, data, opts...)
}

func UpdateOne(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mgodb.Collection(sqltable).UpdateOne(mgodb.getContext(), filter, update, opts...)
}

func UpdateMany(sqltable core.SqlTable, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mgodb.Collection(sqltable).UpdateMany(mgodb.getContext(), filter, update, opts...)
}

func UpdateManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mgodb.Collection(sqltable).UpdateMany(ctx, filter, update, opts...)
}

func FindOneAndDelete(sqltable core.SqlTable, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return mgodb.Collection(sqltable).FindOneAndDelete(mgodb.getContext(), filter, opts...)
}

func DeleteOne(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return mgodb.Collection(sqltable).DeleteOne(mgodb.getContext(), filter, opts...)
}

func DeleteMany(sqltable core.SqlTable, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return mgodb.Collection(sqltable).DeleteMany(mgodb.getContext(), filter, opts...)
}

func DeleteManyByCtx(sqltable core.SqlTable, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return mgodb.Collection(sqltable).DeleteMany(ctx, filter, opts...)
}
