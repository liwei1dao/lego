package console

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/event"
	"github.com/liwei1dao/lego/sys/mgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBComp struct {
	cbase.ModuleCompBase
	module IConsole
	mgo    mgo.IMongodb
}

func (this *DBComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.module = module.(IConsole)
	if this.mgo, err = mgo.NewSys(mgo.SetMongodbUrl(this.module.Options().GetMongodbUrl()), mgo.SetMongodbDatabase(this.module.Options().GetMongodbDatabase())); err == nil {
		err = this.checkDbInit()
	}
	return
}

//校验数据库初始化工作是否完成
func (this DBComp) checkDbInit() (err error) {
	count, err := this.mgo.CountDocuments(Sql_ConsoleUserDataIdTable, bson.M{})
	if err != nil || count == 0 {
		//批量插入数据
		data := make([]interface{}, this.module.Options().GetInitUserIdNum())
		r := rand.New(rand.NewSource(time.Now().Unix()))
		n := 0
		for _, i := range r.Perm(this.module.Options().GetInitUserIdNum()) {
			data[n] = bson.M{"_id": i}
			n++
		}
		var (
			err error
		)
		if _, err = this.mgo.InsertManyByCtx(Sql_ConsoleUserDataIdTable, context.TODO(), data); err != nil {
			return fmt.Errorf("CheckDbInit  err=%s", err.Error())
		}
	}
	return
}

func (this *DBComp) GetMgo() mgo.IMongodb {
	return this.mgo
}

/*				 User相关接口
 * _______________#########_______________________
 * ______________############_____________________
 * ______________#############____________________
 * _____________##__###########___________________
 * ____________###__######_#####__________________
 * ____________###_#######___####_________________
 * ___________###__##########_####________________
 * __________####__###########_####_______________
 * ________#####___###########__#####_____________
 * _______######___###_########___#####___________
 * _______#####___###___########___######_________
 * ______######___###__###########___######_______
 * _____######___####_##############__######______
 * ____#######__#####################_#######_____
 * ____#######__##############################____
 * ___#######__######_#################_#######___
 * ___#######__######_######_#########___######___
 * ___#######____##__######___######_____######___
 * ___#######________######____#####_____#####____
 * ____######________#####_____#####_____####_____
 * _____#####________####______#####_____###______
 * ______#####______;###________###______#________
 * ________##_______####________####______________
 */

//查询用户 Id
func (this *DBComp) QueryUserDataById(uid uint32) (result *DB_UserData, err error) {
	result = &DB_UserData{}
	err = this.mgo.FindOne(Sql_ConsoleUserDataTable, bson.M{"_id": uid}).Decode(result)
	return
}

//查询用户 Id
func (this *DBComp) QueryUserDataByPhonOrEmail(phonoremail string) (result *DB_UserData, err error) {
	result = &DB_UserData{}
	err = this.mgo.FindOne(Sql_ConsoleUserDataTable, bson.M{"phonoremail": phonoremail}).Decode(result)
	return
}

//查询用户 账号
func (this *DBComp) LoginUserDataByPhonOrEmail(data *DB_UserData) (result *DB_UserData, err error) {
	result = &DB_UserData{}
	if err := this.mgo.FindOne(Sql_ConsoleUserDataTable, bson.M{"phonoremail": data.PhonOrEmail}).Decode(result); err != nil {
		return this.registeredUser(data)
	}
	return
}

//更新用户信息
func (this *DBComp) UpDataUserData(data *DB_UserData) (err error) {
	_, err = this.mgo.UpdateOne(Sql_ConsoleUserDataTable,
		bson.M{"_id": data.Id},
		bson.M{"$set": bson.M{
			"nickname": data.NickName,
			"headurl":  data.HeadUrl,
			"userrole": data.UserRole,
			"token":    data.Token,
		}},
		new(options.UpdateOptions).SetUpsert(true))
	return
}

//注册账号
func (this *DBComp) registeredUser(data *DB_UserData) (result *DB_UserData, err error) {
	var newuId ConsoleUserId
	if err = this.mgo.FindOneAndDelete(Sql_ConsoleUserDataIdTable, bson.D{}).Decode(&newuId); err != nil {
		return
	}
	result = &DB_UserData{
		Id:          newuId.Id,
		PhonOrEmail: data.PhonOrEmail,
		Password:    data.Password,
		NickName:    data.NickName,
		HeadUrl:     data.HeadUrl,
	}
	if _, err = this.mgo.InsertOne(Sql_ConsoleUserDataTable, result); err == nil {
		event.TriggerEvent(Event_Registered, result)
	}
	return
}
