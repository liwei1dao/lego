package mgo

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userId struct {
	ID uint32 `bson:"_id"`
}

type userData struct {
	Id               uint32 `bson:"_id"` //用户id
	RegistrationTime int64  //注册时间
	Source           uint8  //账户来源 0测试账号 1手机注册 2微信 ...
	Account          string //账号名称
	Password         string //密码
	NiceName         string //昵称
	HeadId           uint8  //头像Id
	Sex              uint8  //性别	1男 2女
}

func TestInsertMany(t *testing.T) {
	err := OnInit(nil,
		MongodbUrl("mongodb://127.0.0.1:27017"),
		MongodbDatabase("test"),
		TimeOut(60*time.Second),
	)
	if err != nil {
		t.Error(err)
		return
	}
	//批量插入数据
	leng := 1000000
	cIds := make([]interface{}, leng)
	for i, _ := range cIds {
		cIds[i] = 10000 + i
	}
	data := make([]interface{}, leng)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	n := 0
	for _, i := range r.Perm(leng) {
		data[n] = bson.M{"_id": i}
		n++
	}
	begin := time.Now()
	if _, err = InsertMany("userid", data); err != nil {
		t.Errorf("执行 InitUserIdTable 数据库操作 1 错误 err=%s", err.Error())
	}
	t.Logf("执行 InitUserIdTable 成功 耗时 %v", time.Now().Sub(begin))
}

func TestConcurrent(t *testing.T) {
	err := OnInit(nil,
		MongodbUrl("mongodb://127.0.0.1:27017"),
		MongodbDatabase("test"),
		TimeOut(5*time.Second),
	)
	if err != nil {
		t.Error(err)
		return
	}
	//wg := new(sync.WaitGroup)
	//wg.Add(10000)
	begin := time.Now()
	for i := 0; i < 10000; i++ {
		//go func(account,password string) {
		//	defer wg.Done()
		if count, err := CountDocuments("userdata", bson.M{"account": fmt.Sprintf("liwei%ddao", i)}, new(options.CountOptions).SetLimit(1)); err != nil {
			t.Errorf("执行 RegisterTestAccount 数据库操作 1 错误 err=%s", err.Error())
		} else {
			if count != 0 {
				t.Errorf("执行 RegisterTestAccount 数据库操作 2 错误 count=%d", count)
				return
			}
		}
		//var newuId userId
		//if err := FindOneAndDelete("userid", bson.D{}).Decode(&newuId); err != nil {
		//	t.Errorf("执行 RegisterTestAccount 数据库操作 3 错误 err=%s", err.CodeError())
		//	return
		//}
		//udata := &userData{
		//	Id:       newuId.ID,
		//	Source:   0,
		//	Account:  account,
		//	Password: password,
		//	NiceName: "测试用户_" + strconv.Itoa(int(newuId.ID)%100),
		//	HeadId:   1,
		//	Sex:      1,
		//}
		//if _, err := InsertOne("userdata", udata,new(options.InsertOneOptions).SetBypassDocumentValidation(true)); err != nil {
		//	t.Errorf("执行 RegisterTestAccount 数据库操作 4 错误 err=%s", err.CodeError())
		//	return
		//}
		//return
		//}(fmt.Sprintf("liwei%ddao",i),"111111")
	}

	//wg.Wait()
	t.Logf("执行 高并发注册 成功 耗时 %v", time.Now().Sub(begin))
}
