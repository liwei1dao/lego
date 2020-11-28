package mgo

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/core"
	"go.mongodb.org/mongo-driver/bson"
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

func Test_sys(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"MongodbUrl":      "mongodb://127.0.0.1:27017",
		"MongodbDatabase": "testdb",
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
		if err := FindOne(Sql_UserDataTable, bson.M{"_id": uId}).Decode(data); err != nil {
			fmt.Printf("FindOne errr:%v", err)
		} else {
			fmt.Printf("FindOne data:%+v", data)
		}
	}
}
