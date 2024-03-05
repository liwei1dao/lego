package model

import (
	"context"
	"fmt"
	"reflect"
	"time"
	"unsafe"

	goredis "github.com/go-redis/redis/v8"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/redis"
	"github.com/liwei1dao/lego/sys/redis/pipe"
	"github.com/liwei1dao/lego/utils/codec"
	"github.com/liwei1dao/lego/utils/codec/codecore"
	"github.com/liwei1dao/lego/utils/codec/json"
	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var defconf = &codecore.Config{
	SortMapKeys:           true,
	IndentionStep:         1,
	OnlyTaggedField:       false,
	DisallowUnknownFields: false,
	CaseSensitive:         false,
	TagKey:                "json",
}

//文档模型
type Document struct {
	Model     *Model
	TableName string
	Expired   time.Duration //过期时间
}

func (this *Document) ReadKey(key string) string {
	return fmt.Sprintf("%s-%s:%s", this.Model.options.Tag, this.TableName, key)
}
func (this *Document) ReadListKey(uid string, id string) string {
	return fmt.Sprintf("%s-%s:%s-%s", this.Model.options.Tag, this.TableName, uid, id)
}

//创建锁对象
func (this *Document) NewRedisMutex(key string) (result *redis.RedisMutex, err error) {
	result, err = this.Model.Redis.NewRedisMutex(key)
	return
}

//创建锁对象 自定义过期时间
func (this *Document) NewRedisMutexForExpiry(key string, outtime int) (result *redis.RedisMutex, err error) {
	result, err = this.Model.Redis.NewRedisMutex(key)
	return
}

// 添加新的数据
func (this *Document) Add(uid string, data interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("DBModel Add", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "data", Value: data})
	key := this.ReadKey(uid)
	if err = this.Model.Redis.HMSet(this.ReadKey(key), data); err != nil {
		return
	}
	option := newDBOption(opt...)
	if option.IsMgoLog {
		err = this.Model.InsertModelLogs(this.TableName, uid, []interface{}{data})
	}
	if this.Expired > 0 {
		this.Model.UpDateModelExpired(key, nil, this.Expired)
	}
	return
}

// 添加新的数据
func (this *Document) Adds(data map[string]interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("DBModel Add", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "data", Value: data})
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	for k, v := range data {
		pipe.HMSet(this.ReadKey(k), v)
	}
	if _, err = pipe.Exec(); err != nil {
		return
	}
	option := newDBOption(opt...)
	if option.IsMgoLog {
		err = this.Model.InsertManyModelLogs(this.TableName, data)
	}
	if this.Expired > 0 {
		for k, _ := range data {
			this.Model.UpDateModelExpired(this.ReadKey(k), nil, this.Expired)
		}
	}
	return
}

// 添加新的数据到列表
func (this *Document) AddListOne(uid string, id string, data interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("AddListOne", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "_id", Value: id}, log.Field{Key: "data", Value: data})
	key := this.ReadListKey(uid, id)
	if err = this.Model.Redis.HMSet(key, data); err != nil {
		return
	}
	if err = this.Model.Redis.HSet(this.ReadKey(uid), id, key); err != nil {
		return
	}
	option := newDBOption(opt...)
	if option.IsMgoLog {
		err = this.Model.InsertModelLogs(this.TableName, uid, []interface{}{data})
	}
	if this.Expired > 0 {
		this.Model.UpDateModelExpired(this.ReadKey(uid), nil, this.Expired)
	}
	return
}

// 添加新的多个数据到列表 data map[string]type
func (this *Document) AddListsMany(uid string, data interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("AddListsMany", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "data", Value: data})
	vof := reflect.ValueOf(data)
	if !vof.IsValid() {
		return fmt.Errorf("Model_Comp: AddLists(nil)")
	}
	if vof.Kind() != reflect.Map {
		return fmt.Errorf("Model_Comp: AddLists(non-pointer %T)", data)
	}
	listskeys := make(map[string]string)
	keys := vof.MapKeys()
	lists := make([]interface{}, 0, len(keys))
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	for _, k := range keys {
		value := vof.MapIndex(k)
		keydata := k.Interface().(string)
		valuedata := value.Interface()
		key := this.ReadListKey(uid, keydata)
		pipe.HMSet(key, valuedata)
		listskeys[keydata] = key
		lists = append(lists, valuedata)
	}
	pipe.HMSetForMap(this.ReadKey(uid), listskeys)
	if _, err = pipe.Exec(); err != nil {
		return
	}
	option := newDBOption(opt...)
	if option.IsMgoLog {
		err = this.Model.InsertModelLogs(this.TableName, uid, lists)
	}
	if this.Expired > 0 {
		childs := make(map[string]struct{}, len(listskeys))
		for _, v := range listskeys {
			childs[v] = struct{}{}
		}
		this.Model.UpDateModelExpired(this.ReadKey(uid), childs, this.Expired)
	}
	return
}

// 添加队列
func (this *Document) AddQueues(key string, uplimit int64, data interface{}) (outkey []string, err error) {
	defer this.Model.options.Log.Debug("AddQueues", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "key", Value: key}, log.Field{Key: "data", Value: data})
	vof := reflect.ValueOf(data)
	if !vof.IsValid() {
		err = fmt.Errorf("Model_Comp: AddLists(nil)")
		return
	}
	if vof.Kind() != reflect.Map {
		err = fmt.Errorf("Model_Comp: AddLists(non-pointer %T)", data)
		return
	}
	keys := make([]string, 0)
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	for _, k := range vof.MapKeys() {
		value := vof.MapIndex(k)
		tkey := k.Interface().(string)
		valuedata := value.Interface()
		pipe.HMSet(tkey, valuedata)
		keys = append(keys, tkey)
	}
	pipe.RPushForStringSlice(key, keys...)
	lcmd := pipe.Llen(key)

	if _, err = pipe.Exec(); err == nil {
		if lcmd.Val() > uplimit*3 { //操作3倍上限移除多余数据
			off := uplimit - lcmd.Val()
			if outkey, err = this.Model.Redis.LRangeToStringSlice(key, 0, int(off-1)).Result(); err != nil {
				return
			}
			pipe.Ltrim(key, int(off), -1)
			for _, v := range outkey {
				pipe.Delete(v)
			}
			_, err = pipe.Exec()
		}
	}
	return
}

// 修改数据多个字段
func (this *Document) Change(uid string, data map[string]interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("Change", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "data", Value: data})
	if err = this.Model.Redis.HMSet(this.ReadKey(uid), data); err != nil {
		return
	}
	option := newDBOption(opt...)
	if option.IsMgoLog {
		err = this.Model.UpdateModelLogs(this.TableName, uid, bson.M{"uid": uid}, data)
	}
	if this.Expired > 0 {
		this.Model.UpDateModelExpired(this.ReadKey(uid), nil, this.Expired)
	}
	return nil
}

// 修改数据多个字段 主键id
func (this *Document) ChangeForId(id string, data map[string]interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("DBModel ChangeForId", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "id", Value: id}, log.Field{Key: "data", Value: data})
	if err = this.Model.Redis.HMSet(this.ReadKey(id), data); err != nil {
		return
	}
	option := newDBOption(opt...)
	if option.IsMgoLog {
		err = this.Model.UpdateModelLogs(this.TableName, "", bson.M{"_id": id}, data)
	}
	if this.Expired > 0 {
		this.Model.UpDateModelExpired(this.ReadKey(id), nil, this.Expired)
	}
	return nil
}

/*
修改列表中一个数据
*/
func (this *Document) ChangeListOne(uid string, _id string, data map[string]interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("ChangeListOne", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "_id", Value: _id}, log.Field{Key: "data", Value: data})
	if err = this.Model.Redis.HMSet(this.ReadListKey(uid, _id), data); err != nil {
		log.Error("DBModel ChangeList", log.Field{Key: "err", Value: err.Error()})
		return
	}

	option := newDBOption(opt...)
	if option.IsMgoLog {
		if uid == "" {
			err = this.Model.UpdateModelLogs(this.TableName, uid, bson.M{"_id": _id}, data)
		} else {
			err = this.Model.UpdateModelLogs(this.TableName, uid, bson.M{"_id": _id, "uid": uid}, data)
		}
	}
	if this.Expired > 0 {
		this.Model.UpDateModelExpired(this.ReadKey(uid), nil, this.Expired)
	}
	return nil
}

/*
修改列表中多个数据 datas key是 _id value是 这个数据对象
*/
func (this *Document) ChangeLists(uid string, datas map[string]interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("ChangeLists", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "datas", Value: datas})
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	for k, v := range datas {
		if err = pipe.HMSet(this.ReadListKey(uid, k), v); err != nil {
			this.Model.options.Log.Error("ChangeLists", log.Field{Key: "err", Value: err.Error()})
			return
		}
	}
	pipe.Exec()

	option := newDBOption(opt...)
	if option.IsMgoLog {
		if uid == "" {
			for k, v := range datas {
				err = this.Model.UpdateModelLogs(this.TableName, uid, bson.M{"_id": k}, v)
			}
		} else {
			for k, v := range datas {
				err = this.Model.UpdateModelLogs(this.TableName, uid, bson.M{"_id": k, "uid": uid}, v)
			}
		}
	}
	if this.Expired > 0 {
		this.Model.UpDateModelExpired(this.ReadKey(uid), nil, this.Expired)
	}
	return nil
}

//批量修改多个用户数据
func (this *Document) BatchChange(uids []string, datas map[string]interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("BatchChange", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uids", Value: uids}, log.Field{Key: "datas", Value: datas})
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	for _, uid := range uids {
		if err = pipe.HMSet(this.ReadKey(uid), datas); err != nil {
			this.Model.options.Log.Error("DBModel ChangeList", log.Field{Key: "err", Value: err.Error()})
			return
		}
	}
	pipe.Exec()

	option := newDBOption(opt...)
	if option.IsMgoLog {
		for _, uid := range uids {
			err = this.Model.UpdateModelLogs(this.TableName, uid, bson.M{"uid": uid}, datas)
		}
	}
	if this.Expired > 0 {
		for _, uid := range uids {
			this.Model.UpDateModelExpired(this.ReadKey(uid), nil, this.Expired)
		}
	}
	return nil
}

// 读取全部数据
func (this *Document) Get(uid string, data interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("Get", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "data", Value: data})
	if err = this.Model.Redis.HGetAll(this.ReadKey(uid), data); err != nil && err != redis.RedisNil {
		return
	}
	if err == redis.RedisNil {
		if err = this.Model.Mgo.FindOne(core.SqlTable(this.TableName), bson.M{"uid": uid}).Decode(data); err != nil {
			return
		}
		err = this.Model.Redis.HMSet(this.ReadKey(uid), data)
	}
	if this.Expired > 0 {
		this.Model.UpDateModelExpired(this.ReadKey(uid), nil, this.Expired)
	}
	return
}

/*
读取多个用户的数据
data 必须是切片指针 &[]interface
*/
func (this *Document) GetMany(uids []string, data interface{}, opt ...DBOption) (onfound []string, err error) {
	defer this.Model.options.Log.Debug("GetMany", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uids", Value: uids}, log.Field{Key: "data", Value: data}, log.Field{Key: "data", Value: data})
	var (
		dtype         reflect2.Type
		dkind         reflect.Kind
		sType         reflect2.Type
		sliceType     *reflect2.UnsafeSliceType
		sliceelemType reflect2.Type
		decoder       codecore.IDecoderMapJson
		encoder       codecore.IEncoderMapJson
		dptr          unsafe.Pointer
		elemPtr       unsafe.Pointer
		n             int
		ok            bool
		keys          map[string]string = make(map[string]string)
		tempdata      map[string]string
		pipe          *pipe.RedisPipe               = this.Model.Redis.RedisPipe(context.TODO())
		result        []*goredis.StringStringMapCmd = make([]*goredis.StringStringMapCmd, len(uids))
		c             *mongo.Cursor
	)
	onfound = make([]string, 0, len(uids))
	dptr = reflect2.PtrOf(data)
	dtype = reflect2.TypeOf(data)
	dkind = dtype.Kind()
	if dkind != reflect.Ptr {
		err = fmt.Errorf("GetMany(non-pointer %T)", data)
		return
	}
	sType = dtype.(*reflect2.UnsafePtrType).Elem()
	if sType.Kind() != reflect.Slice {
		err = fmt.Errorf("GetMany(data no slice %T)", data)
		return
	}
	sliceType = sType.(*reflect2.UnsafeSliceType)
	sliceelemType = sliceType.Elem()
	if sliceelemType.Kind() != reflect.Ptr {
		err = fmt.Errorf("GetMany(sliceelemType non-pointer %T)", data)
		return
	}
	if decoder, ok = codec.DecoderOf(sliceelemType, defconf).(codecore.IDecoderMapJson); !ok {
		err = fmt.Errorf("GetMany(data not support MarshalMapJson %T)", data)
		return
	}
	sliceelemType = sliceelemType.(*reflect2.UnsafePtrType).Elem()
	for i, v := range uids {
		result[i] = pipe.HGetAllToMapString(this.ReadKey(v))
	}
	if _, err = pipe.Exec(); err == nil {
		for i, v := range result {
			if tempdata, err = v.Result(); err == nil && len(tempdata) > 0 {
				sliceType.UnsafeGrow(dptr, n+1)
				elemPtr = sliceType.UnsafeGetIndex(dptr, n)
				if *((*unsafe.Pointer)(elemPtr)) == nil {
					newPtr := sliceelemType.UnsafeNew()
					if err = decoder.DecodeForMapJson(newPtr, json.GetReader([]byte{}, nil), tempdata); err != nil {
						log.Errorf("err:%v", err)
						return
					}
					*((*unsafe.Pointer)(elemPtr)) = newPtr
				} else {
					decoder.DecodeForMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetReader([]byte{}, nil), tempdata)
				}
				n++
			} else {
				onfound = append(onfound, uids[i])
			}
		}
	} else {
		onfound = uids
	}

	if len(onfound) > 0 {
		if c, err = this.Model.Mgo.Find(core.SqlTable(this.TableName), bson.M{"uid": bson.M{"$in": onfound}}); err != nil {
			return
		} else {
			if encoder, ok = codec.EncoderOf(sliceelemType, defconf).(codecore.IEncoderMapJson); !ok {
				err = fmt.Errorf("GetMany(data not support UnMarshalMapJson %T)", data)
				return
			}
			pipe := this.Model.Redis.RedisPipe(context.TODO())
			for c.Next(context.Background()) {
				_id := c.Current.Lookup("_id").StringValue()
				sliceType.UnsafeGrow(dptr, n+1)
				elemPtr = sliceType.UnsafeGetIndex(dptr, n)
				if *((*unsafe.Pointer)(elemPtr)) == nil {
					newPtr := sliceelemType.UnsafeNew()
					*((*unsafe.Pointer)(elemPtr)) = newPtr
				}
				elem := sliceType.GetIndex(data, n)
				if err = c.Decode(elem); err != nil {
					return
				}
				if tempdata, err = encoder.EncodeToMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetWriter(nil)); err != nil {
					return
				}
				key := this.ReadKey(_id)
				pipe.HMSetForMap(key, tempdata)
				keys[_id] = key
				n++
				for i, v := range onfound {
					if v == _id {
						onfound = append(onfound[:i], onfound[i+1:]...)
						break
					}
				}

			}
			if len(keys) > 0 {
				_, err = pipe.Exec()
			}
		}
	}
	if this.Expired > 0 {
		for _, v := range uids {
			this.Model.UpDateModelExpired(this.ReadKey(v), nil, this.Expired)
		}
	}
	return
}

/*
读取全部数据 唯一id
*/
func (this *Document) GetForId(id string, data interface{}, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("GetForId", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "id", Value: id}, log.Field{Key: "data", Value: data})
	if err = this.Model.Redis.HGetAll(this.ReadKey(id), data); err != nil && err != redis.RedisNil {
		log.Error("GetForId Err!", log.Field{Key: "ukey", Value: this.ReadKey(id)}, log.Field{Key: "err", Value: err.Error()})
		return
	}
	if err == redis.RedisNil {
		if err = this.Model.Mgo.FindOne(core.SqlTable(this.TableName), bson.M{"_id": id}).Decode(data); err != nil {
			return
		}
		err = this.Model.Redis.HMSet(this.ReadKey(id), data)
	}
	if this.Expired > 0 {
		this.Model.UpDateModelExpired(this.ReadKey(id), nil, this.Expired)
	}
	return
}

// 读取多个数据对象
func (this *Document) GetManyForId(ids []string, data interface{}, opt ...DBOption) (onfound []string, err error) {
	defer this.Model.options.Log.Debug("DBModel GetListObjs", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "ids", Value: ids}, log.Field{Key: "data", Value: data})
	var (
		dtype         reflect2.Type
		dkind         reflect.Kind
		sType         reflect2.Type
		sliceType     *reflect2.UnsafeSliceType
		sliceelemType reflect2.Type
		decoder       codecore.IDecoderMapJson
		encoder       codecore.IEncoderMapJson
		dptr          unsafe.Pointer
		elemPtr       unsafe.Pointer
		n             int
		ok            bool
		keys          map[string]string = make(map[string]string)
		tempdata      map[string]string
		pipe          *pipe.RedisPipe               = this.Model.Redis.RedisPipe(context.TODO())
		result        []*goredis.StringStringMapCmd = make([]*goredis.StringStringMapCmd, len(ids))
		c             *mongo.Cursor
	)
	onfound = make([]string, 0, len(ids))
	dptr = reflect2.PtrOf(data)
	dtype = reflect2.TypeOf(data)
	dkind = dtype.Kind()
	if dkind != reflect.Ptr {
		err = fmt.Errorf("GetManyForId(non-pointer %T)", data)
		return
	}
	sType = dtype.(*reflect2.UnsafePtrType).Elem()
	if sType.Kind() != reflect.Slice {
		err = fmt.Errorf("GetManyForId(data no slice %T)", data)
		return
	}
	sliceType = sType.(*reflect2.UnsafeSliceType)
	sliceelemType = sliceType.Elem()
	if sliceelemType.Kind() != reflect.Ptr {
		err = fmt.Errorf("GetManyForId(sliceelemType non-pointer %T)", data)
		return
	}
	if decoder, ok = codec.DecoderOf(sliceelemType, defconf).(codecore.IDecoderMapJson); !ok {
		err = fmt.Errorf("GetManyForId(data not support MarshalMapJson %T)", data)
		return
	}
	sliceelemType = sliceelemType.(*reflect2.UnsafePtrType).Elem()
	for i, v := range ids {
		result[i] = pipe.HGetAllToMapString(this.ReadKey(v))
	}
	if _, err = pipe.Exec(); err == nil {
		for i, v := range result {
			if tempdata, err = v.Result(); err == nil && len(tempdata) > 0 {
				sliceType.UnsafeGrow(dptr, n+1)
				elemPtr = sliceType.UnsafeGetIndex(dptr, n)
				if *((*unsafe.Pointer)(elemPtr)) == nil {
					newPtr := sliceelemType.UnsafeNew()
					if err = decoder.DecodeForMapJson(newPtr, json.GetReader([]byte{}, nil), tempdata); err != nil {
						log.Errorf("err:%v", err)
						return
					}
					*((*unsafe.Pointer)(elemPtr)) = newPtr
				} else {
					decoder.DecodeForMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetReader([]byte{}, nil), tempdata)
				}
				n++
			} else {
				onfound = append(onfound, ids[i])
			}
		}
	} else {
		onfound = ids
	}

	if len(onfound) > 0 {
		if c, err = this.Model.Mgo.Find(core.SqlTable(this.TableName), bson.M{"_id": bson.M{"$in": onfound}}); err != nil {
			return
		} else {
			if encoder, ok = codec.EncoderOf(sliceelemType, defconf).(codecore.IEncoderMapJson); !ok {
				err = fmt.Errorf("GetManyForId(data not support UnMarshalMapJson %T)", data)
				return
			}
			pipe := this.Model.Redis.RedisPipe(context.TODO())
			for c.Next(context.Background()) {
				_id := c.Current.Lookup("_id").StringValue()
				sliceType.UnsafeGrow(dptr, n+1)
				elemPtr = sliceType.UnsafeGetIndex(dptr, n)
				if *((*unsafe.Pointer)(elemPtr)) == nil {
					newPtr := sliceelemType.UnsafeNew()
					*((*unsafe.Pointer)(elemPtr)) = newPtr
				}
				elem := sliceType.GetIndex(data, n)
				if err = c.Decode(elem); err != nil {
					return
				}
				if tempdata, err = encoder.EncodeToMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetWriter(nil)); err != nil {
					return
				}
				key := this.ReadKey(_id)
				pipe.HMSetForMap(key, tempdata)
				keys[_id] = key
				n++
				for i, v := range onfound {
					if v == _id {
						onfound = append(onfound[:i], onfound[i+1:]...)
						break
					}
				}

			}
			if len(keys) > 0 {
				_, err = pipe.Exec()
			}
		}
	}
	if this.Expired > 0 {
		for _, v := range ids {
			this.Model.UpDateModelExpired(this.ReadKey(v), nil, this.Expired)
		}
	}
	return
}

/*
获取列表数据 注意 data 必须是 切片的指针 *[]type
*/
func (this *Document) GetList(uid string, data interface{}) (err error) {
	defer this.Model.options.Log.Debug("GetList", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "data", Value: data})
	var (
		dtype         reflect2.Type
		dkind         reflect.Kind
		sType         reflect2.Type
		sliceType     *reflect2.UnsafeSliceType
		sliceelemType reflect2.Type
		decoder       codecore.IDecoderMapJson
		encoder       codecore.IEncoderMapJson
		dptr          unsafe.Pointer
		elemPtr       unsafe.Pointer
		n             int
		ok            bool
		keys          map[string]string
		tempdata      map[string]string
		result        []*goredis.StringStringMapCmd
		c             *mongo.Cursor
	)
	keys = make(map[string]string)
	dptr = reflect2.PtrOf(data)
	dtype = reflect2.TypeOf(data)
	dkind = dtype.Kind()
	if dkind != reflect.Ptr {
		err = fmt.Errorf("GetList(non-pointer %T)", data)
		return
	}
	sType = dtype.(*reflect2.UnsafePtrType).Elem()
	if sType.Kind() != reflect.Slice {
		err = fmt.Errorf("GetList(data no slice %T)", data)
		return
	}
	sliceType = sType.(*reflect2.UnsafeSliceType)
	sliceelemType = sliceType.Elem()
	if sliceelemType.Kind() != reflect.Ptr {
		err = fmt.Errorf("GetList(sliceelemType non-pointer %T)", data)
		return
	}
	if decoder, ok = codec.DecoderOf(sliceelemType, defconf).(codecore.IDecoderMapJson); !ok {
		err = fmt.Errorf("GetList(data not support MarshalMapJson %T)", data)
		return
	}
	sliceelemType = sliceelemType.(*reflect2.UnsafePtrType).Elem()
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	if keys, err = this.Model.Redis.HGetAllToMapString(this.ReadKey(uid)); err == nil {
		result = make([]*goredis.StringStringMapCmd, 0)
		for _, v := range keys {
			cmd := pipe.HGetAllToMapString(v)
			result = append(result, cmd)
		}
		pipe.Exec()
		for _, v := range result {
			tempdata, err = v.Result()
			sliceType.UnsafeGrow(dptr, n+1)
			elemPtr = sliceType.UnsafeGetIndex(dptr, n)
			if *((*unsafe.Pointer)(elemPtr)) == nil {
				newPtr := sliceelemType.UnsafeNew()
				if err = decoder.DecodeForMapJson(newPtr, json.GetReader([]byte{}, nil), tempdata); err != nil {
					log.Errorf("err:%v", err)
					return
				}
				*((*unsafe.Pointer)(elemPtr)) = newPtr
			} else {
				decoder.DecodeForMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetReader([]byte{}, nil), tempdata)
			}
			n++
		}
	}
	if err == redis.RedisNil {
		var f = bson.M{}
		if uid != "" {
			f = bson.M{"uid": uid}
		}
		//query from mgo
		if c, err = this.Model.Mgo.Find(core.SqlTable(this.TableName), f); err != nil {
			return err
		} else {
			if encoder, ok = codec.EncoderOf(sliceelemType, defconf).(codecore.IEncoderMapJson); !ok {
				err = fmt.Errorf("GetList(data not support UnMarshalMapJson %T)", data)
				return
			}
			n = 0
			for c.Next(context.Background()) {
				_id := c.Current.Lookup("_id").StringValue()
				sliceType.UnsafeGrow(dptr, n+1)
				elemPtr = sliceType.UnsafeGetIndex(dptr, n)
				if *((*unsafe.Pointer)(elemPtr)) == nil {
					newPtr := sliceelemType.UnsafeNew()
					*((*unsafe.Pointer)(elemPtr)) = newPtr
				}
				elem := sliceType.GetIndex(data, n)
				if err = c.Decode(elem); err != nil {
					return
				}
				if tempdata, err = encoder.EncodeToMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetWriter(nil)); err != nil {
					return
				}
				key := this.ReadListKey(uid, _id)
				pipe.HMSetForMap(key, tempdata)
				keys[_id] = key
				n++
			}
			if len(keys) > 0 {
				pipe.HMSetForMap(this.ReadKey(uid), keys)
				_, err = pipe.Exec()
			}
		}
	}
	if this.Expired > 0 {
		childs := make(map[string]struct{}, len(keys))
		for _, v := range keys {
			childs[v] = struct{}{}
		}
		this.Model.UpDateModelExpired(this.ReadKey(uid), childs, this.Expired)
	}
	return err
}

// 随机读取列表数据
func (this *Document) RandGetList(uid string, num int32, data interface{}) (err error) {
	defer this.Model.options.Log.Debug("RandGetList", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "data", Value: data})
	var (
		dtype         reflect2.Type
		dkind         reflect.Kind
		sType         reflect2.Type
		sliceType     *reflect2.UnsafeSliceType
		sliceelemType reflect2.Type
		decoder       codecore.IDecoderMapJson
		dptr          unsafe.Pointer
		elemPtr       unsafe.Pointer
		n             int
		ok            bool
		keys          map[string]string
		tempdata      map[string]string
		result        []*goredis.StringStringMapCmd
	)
	keys = make(map[string]string)
	dptr = reflect2.PtrOf(data)
	dtype = reflect2.TypeOf(data)
	dkind = dtype.Kind()
	if dkind != reflect.Ptr {
		err = fmt.Errorf("RandGetList(non-pointer %T)", data)
		return
	}
	sType = dtype.(*reflect2.UnsafePtrType).Elem()
	if sType.Kind() != reflect.Slice {
		err = fmt.Errorf("RandGetList(data no slice %T)", data)
		return
	}
	sliceType = sType.(*reflect2.UnsafeSliceType)
	sliceelemType = sliceType.Elem()
	if sliceelemType.Kind() != reflect.Ptr {
		err = fmt.Errorf("RandGetList(sliceelemType non-pointer %T)", data)
		return
	}
	if decoder, ok = codec.DecoderOf(sliceelemType, defconf).(codecore.IDecoderMapJson); !ok {
		err = fmt.Errorf("RandGetList(data not support MarshalMapJson %T)", data)
		return
	}
	sliceelemType = sliceelemType.(*reflect2.UnsafePtrType).Elem()
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	if keys, err = this.Model.Redis.HGetAllToMapString(this.ReadKey(uid)); err == nil {
		result = make([]*goredis.StringStringMapCmd, 0)
		for _, v := range keys {
			cmd := pipe.HGetAllToMapString(v)
			result = append(result, cmd)
			if len(result) >= int(num) {
				break
			}
		}
		pipe.Exec()
		for _, v := range result {
			tempdata, err = v.Result()
			sliceType.UnsafeGrow(dptr, n+1)
			elemPtr = sliceType.UnsafeGetIndex(dptr, n)
			if *((*unsafe.Pointer)(elemPtr)) == nil {
				newPtr := sliceelemType.UnsafeNew()
				if err = decoder.DecodeForMapJson(newPtr, json.GetReader([]byte{}, nil), tempdata); err != nil {
					return
				}
				*((*unsafe.Pointer)(elemPtr)) = newPtr
			} else {
				decoder.DecodeForMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetReader([]byte{}, nil), tempdata)
			}
			n++
		}
	}
	if this.Expired > 0 {
		childs := make(map[string]struct{}, len(keys))
		for _, v := range keys {
			childs[v] = struct{}{}
		}
		this.Model.UpDateModelExpired(this.ReadKey(uid), childs, this.Expired)
	}
	return err
}

// 获取队列数据
func (this *Document) GetQueues(key string, count int, data interface{}) (err error) {
	defer this.Model.options.Log.Debug("GetQueues", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "key", Value: key}, log.Field{Key: "data", Value: data})
	var (
		dtype         reflect2.Type
		dkind         reflect.Kind
		sType         reflect2.Type
		sliceType     *reflect2.UnsafeSliceType
		sliceelemType reflect2.Type
		dptr          unsafe.Pointer
		elemPtr       unsafe.Pointer
		decoder       codecore.IDecoderMapJson
		ok            bool
		n             int
		keys          = make([]string, 0)
		result        = make([]*goredis.StringStringMapCmd, 0)
		tempdata      map[string]string
	)

	dptr = reflect2.PtrOf(data)
	dtype = reflect2.TypeOf(data)
	dkind = dtype.Kind()
	if dkind != reflect.Ptr {
		err = fmt.Errorf("GetQueues(non-pointer %T)", data)
		return
	}
	sType = dtype.(*reflect2.UnsafePtrType).Elem()
	if sType.Kind() != reflect.Slice {
		err = fmt.Errorf("GetQueues(data no slice %T)", data)
		return
	}
	sliceType = sType.(*reflect2.UnsafeSliceType)
	sliceelemType = sliceType.Elem()
	if sliceelemType.Kind() != reflect.Ptr {
		err = fmt.Errorf("GetQueues(sliceelemType non-pointer %T)", data)
		return
	}
	if decoder, ok = codec.DecoderOf(sliceelemType, defconf).(codecore.IDecoderMapJson); !ok {
		err = fmt.Errorf("GetQueues(data not support MarshalMapJson %T)", data)
		return
	}
	sliceelemType = sliceelemType.(*reflect2.UnsafePtrType).Elem()

	pipe := this.Model.Redis.RedisPipe(context.TODO())
	if keys, err = this.Model.Redis.LRangeToStringSlice(key, -1*count, -1).Result(); err == nil {
		result = make([]*goredis.StringStringMapCmd, 0)
		for _, v := range keys {
			cmd := pipe.HGetAllToMapString(v)
			result = append(result, cmd)
		}
		pipe.Exec()
		for _, v := range result {
			tempdata, err = v.Result()
			sliceType.UnsafeGrow(dptr, n+1)
			elemPtr = sliceType.UnsafeGetIndex(dptr, n)
			if *((*unsafe.Pointer)(elemPtr)) == nil {
				newPtr := sliceelemType.UnsafeNew()
				if err = decoder.DecodeForMapJson(newPtr, json.GetReader([]byte{}, nil), tempdata); err != nil {
					return
				}
				*((*unsafe.Pointer)(elemPtr)) = newPtr
			} else {
				decoder.DecodeForMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetReader([]byte{}, nil), tempdata)
			}
			n++
		}
	}
	return
}

// 读取单个数据中 多个字段数据
func (this *Document) GetForFields(uid string, data interface{}, fields ...string) (err error) {
	defer this.Model.options.Log.Debug("GetForFields", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "data", Value: data})
	if err = this.Model.Redis.HMGet(this.ReadKey(uid), data, fields...); err != nil && err != redis.RedisNil {
		return
	}
	if err == redis.RedisNil {
		if err = this.Model.Mgo.FindOne(core.SqlTable(this.TableName), bson.M{"uid": uid}).Decode(data); err != nil {
			return err
		}
		if err = this.Model.Redis.HMSet(this.ReadKey(uid), data); err != nil {
			return
		}
	}
	if this.Expired > 0 {
		this.Model.UpDateModelExpired(this.ReadKey(uid), nil, this.Expired)
	}
	return
}

// 读取List列表中单个数据中 多个字段数据
func (this *Document) GetListOneForFields(uid string, id string, data interface{}, fields ...string) (err error) {
	defer this.Model.options.Log.Debug("GetListOneForFields", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "id", Value: id}, log.Field{Key: "data", Value: data})
	var (
		keys     map[string]string
		tempdata map[string]string
	)
	if err = this.Model.Redis.HMGet(this.ReadListKey(uid, id), data, fields...); err != nil && err != redis.RedisNil {
		return
	}
	if err == redis.RedisNil {
		if err = this.Model.Mgo.FindOne(core.SqlTable(this.TableName), bson.M{"_id": id}).Decode(data); err != nil {
			return err
		} else {
			pipe := this.Model.Redis.RedisPipe(context.TODO())
			key := this.ReadListKey(uid, id)
			pipe.HMSetForMap(key, tempdata)
			keys[id] = key
			pipe.HMSetForMap(this.ReadKey(uid), keys)
			if _, err = pipe.Exec(); err != nil {
				return
			}
		}
	}
	if this.Expired > 0 {
		childs := make(map[string]struct{}, len(keys))
		for _, v := range keys {
			childs[v] = struct{}{}
		}
		this.Model.UpDateModelExpired(this.ReadKey(uid), childs, this.Expired)
	}
	return
}

// 读取列表数据中单个数据
func (this *Document) GetListOne(uid string, id string, data interface{}) (err error) {
	defer this.Model.options.Log.Debug("GetListOne", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "id", Value: id}, log.Field{Key: "data", Value: data})
	var (
		keys map[string]string
	)
	keys = make(map[string]string)
	if err = this.Model.Redis.HGetAll(this.ReadListKey(uid, id), data); err != nil && err != redis.RedisNil {
		return
	}
	if err == redis.RedisNil {
		if err = this.Model.Mgo.FindOne(core.SqlTable(this.TableName), bson.M{"_id": id}).Decode(data); err != nil {
			return err
		} else {
			pipe := this.Model.Redis.RedisPipe(context.TODO())
			key := this.ReadListKey(uid, id)
			pipe.HMSet(key, data)
			keys[id] = key
			pipe.HMSetForMap(this.ReadKey(uid), keys)
			if _, err = pipe.Exec(); err != nil {
				return
			}
		}
	}
	if this.Expired > 0 {
		childs := make(map[string]struct{}, len(keys))
		for _, v := range keys {
			childs[v] = struct{}{}
		}
		this.Model.UpDateModelExpired(this.ReadKey(uid), childs, this.Expired)
	}
	return
}

// 读取列表数据中多条数据
func (this *Document) GetListMany(uid string, ids []string, data interface{}) (err error) {
	defer this.Model.options.Log.Debug("GetListMany", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "ids", Value: ids}, log.Field{Key: "data", Value: data})
	var (
		dtype         reflect2.Type
		dkind         reflect.Kind
		sType         reflect2.Type
		sliceType     *reflect2.UnsafeSliceType
		sliceelemType reflect2.Type
		decoder       codecore.IDecoderMapJson
		encoder       codecore.IEncoderMapJson
		dptr          unsafe.Pointer
		elemPtr       unsafe.Pointer
		n             int
		ok            bool
		keys          map[string]string = make(map[string]string)
		tempdata      map[string]string
		onfound       []string                      = make([]string, 0, len(ids))
		pipe          *pipe.RedisPipe               = this.Model.Redis.RedisPipe(context.TODO())
		result        []*goredis.StringStringMapCmd = make([]*goredis.StringStringMapCmd, len(ids))
		c             *mongo.Cursor
	)
	dptr = reflect2.PtrOf(data)
	dtype = reflect2.TypeOf(data)
	dkind = dtype.Kind()
	if dkind != reflect.Ptr {
		err = fmt.Errorf("GetListMany(non-pointer %T)", data)
		return
	}
	sType = dtype.(*reflect2.UnsafePtrType).Elem()
	if sType.Kind() != reflect.Slice {
		err = fmt.Errorf("GetListMany(data no slice %T)", data)
		return
	}
	sliceType = sType.(*reflect2.UnsafeSliceType)
	sliceelemType = sliceType.Elem()
	if sliceelemType.Kind() != reflect.Ptr {
		err = fmt.Errorf("GetListMany(sliceelemType non-pointer %T)", data)
		return
	}
	if decoder, ok = codec.DecoderOf(sliceelemType, defconf).(codecore.IDecoderMapJson); !ok {
		err = fmt.Errorf("GetListMany(data not support MarshalMapJson %T)", data)
		return
	}
	sliceelemType = sliceelemType.(*reflect2.UnsafePtrType).Elem()
	for i, v := range ids {
		result[i] = pipe.HGetAllToMapString(this.ReadListKey(uid, v))
	}
	if _, err = pipe.Exec(); err == nil {
		for i, v := range result {
			if tempdata, err = v.Result(); err == nil && len(tempdata) > 0 {
				sliceType.UnsafeGrow(dptr, n+1)
				elemPtr = sliceType.UnsafeGetIndex(dptr, n)
				if *((*unsafe.Pointer)(elemPtr)) == nil {
					newPtr := sliceelemType.UnsafeNew()
					if err = decoder.DecodeForMapJson(newPtr, json.GetReader([]byte{}, nil), tempdata); err != nil {
						log.Errorf("err:%v", err)
						return
					}
					*((*unsafe.Pointer)(elemPtr)) = newPtr
				} else {
					decoder.DecodeForMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetReader([]byte{}, nil), tempdata)
				}
				n++
			} else {
				onfound = append(onfound, ids[i])
			}
		}
	} else {
		onfound = ids
	}

	if len(onfound) > 0 {
		if c, err = this.Model.Mgo.Find(core.SqlTable(this.TableName), bson.M{"_id": bson.M{"$in": onfound}}); err != nil {
			return err
		} else {
			if encoder, ok = codec.EncoderOf(sliceelemType, defconf).(codecore.IEncoderMapJson); !ok {
				err = fmt.Errorf("MCompModel: GetList(data not support UnMarshalMapJson %T)", data)
				return
			}
			pipe := this.Model.Redis.RedisPipe(context.TODO())
			for c.Next(context.Background()) {
				_id := c.Current.Lookup("_id").StringValue()
				sliceType.UnsafeGrow(dptr, n+1)
				elemPtr = sliceType.UnsafeGetIndex(dptr, n)
				if *((*unsafe.Pointer)(elemPtr)) == nil {
					newPtr := sliceelemType.UnsafeNew()
					*((*unsafe.Pointer)(elemPtr)) = newPtr
				}
				elem := sliceType.GetIndex(data, n)
				if err = c.Decode(elem); err != nil {
					return
				}
				if tempdata, err = encoder.EncodeToMapJson(*((*unsafe.Pointer)(elemPtr)), json.GetWriter(nil)); err != nil {
					return
				}
				for i, v := range onfound {
					if v == _id {
						onfound = append(onfound[0:i], onfound[i+1:]...)
					}
				}
				key := this.ReadListKey(uid, _id)
				pipe.HMSetForMap(key, tempdata)
				keys[_id] = key
				n++
			}
			if len(keys) > 0 {
				pipe.HMSetForMap(this.ReadKey(uid), keys)
				_, err = pipe.Exec()
			}
		}
	}
	if len(onfound) > 0 {
		err = fmt.Errorf("onfound:%v data!", onfound)
	}
	if this.Expired > 0 {
		childs := make(map[string]struct{}, len(keys))
		for _, v := range keys {
			childs[v] = struct{}{}
		}
		this.Model.UpDateModelExpired(this.ReadKey(uid), childs, this.Expired)
	}
	return
}

// 删除用户数据
func (this *Document) Del(uid string, opt ...DBOption) (err error) {
	//defer log.Debug("DBModel Del", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid})
	err = this.Model.Redis.Delete(this.ReadKey(uid))
	if err != nil {
		return err
	}
	option := newDBOption(opt...)
	if option.IsMgoLog {
		err = this.Model.DeleteModelLogs(this.TableName, uid, bson.M{"uid": uid})
	}
	return nil
}

// 删除目标数据 主键id
func (this *Document) DelForId(id string, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("Del", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "id", Value: id})
	err = this.Model.Redis.Delete(this.ReadKey(id))
	if err != nil {
		return err
	}
	option := newDBOption(opt...)
	if option.IsMgoLog {
		err = this.Model.DeleteModelLogs(this.TableName, "", bson.M{"_id": id})
	}
	return nil
}

// 批量删除数据
func (this *Document) DelLists(uid string, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("DelLists", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid})
	var keys map[string]string
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	if keys, err = this.Model.Redis.HGetAllToMapString(this.ReadKey(uid)); err == nil {
		for _, v := range keys {
			pipe.Delete(v)
		}
		pipe.Delete(this.ReadKey(uid))
		pipe.Exec()
	}
	option := newDBOption(opt...)
	if option.IsMgoLog {
		err = this.Model.DeleteModelLogs(this.TableName, uid, bson.M{"uid": uid})
	}
	return
}

// 删除多条数据
func (this *Document) DelListForIds(uid string, ids []string, opt ...DBOption) (err error) {
	defer this.Model.options.Log.Debug("DelListForIds", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uid", Value: uid}, log.Field{Key: "ids", Value: ids})
	listkey := this.ReadKey(uid)
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	for _, v := range ids {
		key := this.ReadListKey(uid, v)
		pipe.Delete(key)
	}
	pipe.HDel(listkey, ids...)
	if _, err = pipe.Exec(); err == nil {
		option := newDBOption(opt...)
		if option.IsMgoLog {
			err = this.Model.DeleteModelLogs(this.TableName, uid, bson.M{"_id": bson.M{"$in": ids}})
		}
	}
	return
}

// 删除卸载redis
func (this *Document) RedisDels(uids []string) (err error) {
	defer this.Model.options.Log.Debug("RedisDels", log.Field{Key: "TableName", Value: this.TableName}, log.Field{Key: "uids", Value: uids})
	pipe := this.Model.Redis.RedisPipe(context.TODO())
	for _, uid := range uids {
		if err = pipe.Delete(this.ReadKey(uid)); err != nil {
			log.Error("RedisDels", log.Field{Key: "err", Value: err.Error()})
			return
		}
	}
	_, err = pipe.Exec()
	return
}
