package mgo

import (
	"context"
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Sql_UserDataTable core.SqlTable = "userdata" //用户
)

type TestData struct {
	Id       uint32 `bson:"_id"` //用户id
	NiceName string //昵称
	Sex      int32  //性别	1男 2女
	HeadUrl  string //头像Id
}

//测试 系统
func Test_sys(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"MongodbUrl":      "mongodb://127.0.0.1:10001",
		"MongodbDatabase": "testdb",
		"Replset":         "replset0",
	}); err == nil {
		if _, err := UpdateOne(Sql_UserDataTable,
			bson.M{"_id": 10001},
			bson.M{"$set": bson.M{
				"nicename": "liwei1dao",
				"headurl":  "http://test.web.com",
				"sex":      1,
			}},
			new(options.UpdateOptions).SetUpsert(true)); err != nil {
			fmt.Printf("UpdateOne errr:%v", err)
		}

		data := &TestData{}
		if err := FindOne(Sql_UserDataTable, bson.M{"_id": 10002}).Decode(data); err != nil {
			fmt.Printf("FindOne errr:%v", err)
		} else {
			fmt.Printf("FindOne data:%+v", data)
		}
	} else {
		fmt.Printf("FindOne errr:%v", err)
	}
}

//测试 事务
func Test_Affair(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"MongodbUrl":      "mongodb://127.0.0.1:27017",
		"MongodbDatabase": "testdb",
	}); err == nil {
		err = UseSession(func(sessionContext mongo.SessionContext) error {
			err := sessionContext.StartTransaction()
			if err != nil {
				fmt.Println(err)
				return err
			}
			col := Collection(Sql_UserDataTable)

			//在事务内写一条id为“222”的记录
			_, err = col.InsertOne(sessionContext, bson.M{"_id": 10002, "nicename": "liwei2dao", "headurl": "http://test1.web.com", "sex": 2})
			if err != nil {
				fmt.Println(err)
				return err
			}
			//在事务内写一条id为“333”的记录
			_, err = col.InsertOne(sessionContext, bson.M{"_id": 10003, "nicename": "liwei3dao", "headurl": "http://test1.web.com", "sex": 2})
			if err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			} else {
				sessionContext.CommitTransaction(sessionContext)
			}
			return nil
		})
		fmt.Printf("FindOne errr:%v", err)
	} else {
		fmt.Printf("FindOne errr:%v", err)
	}
}

//测试创建索引 地图索引
func Test_CreateIndex(t *testing.T) {
	sys, err := NewSys(SetMongodbUrl("mongodb://47.90.84.157:9094"), SetMongodbDatabase("square"))
	if err != nil {
		fmt.Printf("start sys Fail err:%v", err)
		return
	}
	str, err := sys.CreateIndex(core.SqlTable("dynamics"), bson.M{"location": "2dsphere"}, new(options.IndexOptions).SetBits(11111))
	if err != nil {
		fmt.Printf("CreateIndex  err:%v", err)
	} else {
		fmt.Printf("CreateIndex  str:%v", str)
	}
}

//创建复合索引
func Test_CreateCompoundIndex(t *testing.T) {
	sys, err := NewSys(SetMongodbUrl("mongodb://47.90.84.157:9094"), SetMongodbDatabase("lego_yl"))
	if err != nil {
		fmt.Printf("start sys Fail err:%v", err)
		return
	}
	str, err := mgo.CreateIndex(core.SqlTable("unreadmsg"), bson.M{"channeltype": 1, "targetid": 1, "uid": 1, "unreadmsg.sendtime": 1}, nil)
	if err != nil {
		fmt.Printf("CreateIndex  err:%v", err)
	} else {
		fmt.Printf("CreateIndex  str:%v", str)
	}
}

//测试删除索引 地图索引
func Test_DeleteIndex(t *testing.T) {
	sys, err := NewSys(SetMongodbUrl("mongodb://47.90.84.157:9094"), SetMongodbDatabase("square"))
	if err != nil {
		fmt.Printf("start sys Fail err:%v", err)
		return
	}
	_, err = sys.DeleteIndex(core.SqlTable("dynamics"), "location_2dsphere", nil)
	if err != nil {
		fmt.Printf("DeleteIndex err:%v", err)
	} else {
		fmt.Printf("DeleteIndex succ")
	}
}

//测试附近查找
/*
db.<collection>.find( { <location field> :
						{ $near :
							{ $geometry :
								{ type : "Point" ,coordinates : [ <longitude> , <latitude> ] } ,
								$maxDistance : <distance in meters>
								}
							}
						})
*/
func Test_QueryRound(t *testing.T) {
	sys, err := NewSys(SetMongodbUrl("mongodb://47.90.84.157:9094"), SetMongodbDatabase("square"))
	if err != nil {
		fmt.Printf("start sys Fail err:%v", err)
		return
	}
	filter := bson.M{
		"location": bson.D{
			{"$near", bson.D{
				{"$geometry", bson.M{"type": "Point", "coordinates": []float64{113.867584, 22.56532}}},
				{"$maxDistance", 10000},
			}},
		},
	}
	_, err = sys.Find(core.SqlTable("dynamics"), filter)
	if err != nil {
		fmt.Printf("CreateIndex err:%v", err)
	} else {
		fmt.Printf("CreateIndex succ")
	}
	return
}

//聚合查询演示
func Test_Aggregate(t *testing.T) {
	sys, err := NewSys(SetMongodbUrl("mongodb://47.90.84.157:9094"), SetMongodbDatabase("square"))
	if err != nil {
		fmt.Printf("start sys Fail err:%v", err)
		return
	}
	result := make([]bson.M, 0)
	var cursor *mongo.Cursor
	matchStage := bson.D{{"$match", bson.M{"_id": 979332}}}
	unwindStage := bson.D{{"$unwind", "$chatmsg"}}
	match2Stage := bson.D{{"$match", bson.M{"chatmsg.sendtime": bson.M{"$gt": 1611890940100, "$lt": 1611893135929}}}}
	groupStage := bson.D{{"$group", bson.M{"_id": 979332, "chatmsg": bson.M{"$push": "$chatmsg"}}}}
	projectStage := bson.D{{"$project", bson.M{"_id": 1, "chatmsg": 1}}}
	if cursor, err = sys.Aggregate(core.SqlTable("grouprecorddata"), mongo.Pipeline{matchStage, unwindStage, match2Stage, groupStage, projectStage}); err == nil {
		err = cursor.All(context.Background(), &result)
	}
	return
}
