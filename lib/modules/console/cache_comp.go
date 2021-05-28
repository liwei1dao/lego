package console

import (
	"fmt"
	reflect "reflect"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/redis"
)

type CacheComp struct {
	cbase.ModuleCompBase
	module IConsole
	redis  redis.IRedis
}

func (this *CacheComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.module = module.(IConsole)
	if this.redis, err = redis.NewSys(
		redis.SetRedisUrl(this.module.Options().GetRedisUrl()),
		redis.SetRedisDB(this.module.Options().GetRedisDB()),
		redis.SetRedisPassword(this.module.Options().GetRedisPassword()),
	); err != nil {
		err = fmt.Errorf("redis[%s]err:%v", this.module.Options().GetRedisUrl(), err)
	}
	return
}

func (this *CacheComp) GetRedis() redis.IRedis {
	return this.redis
}

/*				Token相关接口
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

//查询用户数据
func (this *CacheComp) QueryToken(token string) (uId uint32, err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleToken), token)
	err = this.redis.Get(core.Redis_Key(Id), &uId)
	return
}

//写入Token
func (this *CacheComp) WriteToken(token string, uId uint32) (err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleToken), token)
	err = this.redis.Set(core.Redis_Key(Id), uId, time.Second*time.Duration(this.module.Options().GetTokenCacheExpirationDate()))
	return
}

//清理Token
func (this *CacheComp) CleanToken(token string) (err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleToken), token)
	err = this.redis.Delete(core.Redis_Key(Id))
	return
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

//查询用户数据
func (this *CacheComp) QueryUserData(uId uint32) (result *Cache_UserData, err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleUsers), uId)
	result = &Cache_UserData{}
	err = this.redis.Get(core.Redis_Key(Id), result)
	return
}

//同步用户数据到缓存
func (this *CacheComp) synchronizeUserToCache(uId uint32) (result *Cache_UserData, err error) {
	var user *DB_UserData
	if user, err = this.module.DB().QueryUserDataById(uId); err == nil {
		result = &Cache_UserData{
			Db_UserData: user,
			IsOnLine:    false,
		}
		this.writeUserDataByEx(result)
	}
	return
}

//离线用户缓存读取之后保存10分钟
func (this *CacheComp) writeUserDataByEx(result *Cache_UserData) (err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleUsers), result.Db_UserData.Id)
	err = this.redis.Set(core.Redis_Key(Id), result, time.Second*time.Duration(this.module.Options().GetUserCacheExpirationDate()))
	return
}

//登录用户缓存信息长期驻留
func (this *CacheComp) WriteUserData(data *Cache_UserData) (err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleUsers), data.Db_UserData.Id)
	err = this.redis.Set(core.Redis_Key(Id), data, 0)
	return
}

//清理用户缓存
func (this *CacheComp) CleanUserData(uid uint32) (err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleUsers), uid)
	err = this.redis.Delete(core.Redis_Key(Id))
	return
}

/*			ClusterMonitor相关接口
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

//添加新的ClusterMonitor
func (this *CacheComp) AddNewClusterMonitor(data map[string]*ClusterMonitor) {
	for k, v := range data {
		id := fmt.Sprintf(string(Cache_ConsoleClusterMonitor), k)
		this.redis.RPush(core.Redis_Key(id), v)
		if len, err := this.redis.Llen(core.Redis_Key(id)); err == nil && len > this.module.Options().GetMonitorTotalTime() {
			this.redis.LPop(Cache_ConsoleClusterMonitor, v)
		}
	}
}

//添加新的ClusterMonitor
func (this *CacheComp) GetClusterMonitor(sIs string, timeleng int) (result []*ClusterMonitor, err error) {
	var values []interface{}
	result = make([]*ClusterMonitor, 0)
	id := fmt.Sprintf(string(Cache_ConsoleClusterMonitor), sIs)
	values, err = this.redis.LRange(core.Redis_Key(id), 0, timeleng, reflect.TypeOf(&ClusterMonitor{}))
	if err == nil && values != nil && len(values) > 0 {
		result = make([]*ClusterMonitor, len(values))
		for i, v := range values {
			result[i] = v.(*ClusterMonitor)
		}
	}
	return
}

/*			  HostMonitor相关接口
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

//添加新的HostMonitor
func (this *CacheComp) AddNewHostMonitor(data *HostMonitor) {
	this.redis.RPush(Cache_ConsoleHostMonitor, data)
	if len, err := this.redis.Llen(Cache_ConsoleHostMonitor); err == nil && len > this.module.Options().GetMonitorTotalTime() {
		this.redis.LPop(Cache_ConsoleHostMonitor, data)
	}
}

//添加新的HostMonitor
func (this *CacheComp) GetHostMonitor(timeleng int) (result []*HostMonitor, err error) {
	var values []interface{}
	result = make([]*HostMonitor, 0)
	values, err = this.redis.LRange(Cache_ConsoleHostMonitor, 0, timeleng, reflect.TypeOf(&HostMonitor{}))
	if values != nil && len(values) > 0 {
		result = make([]*HostMonitor, len(values))
		for i, v := range values {
			result[i] = v.(*HostMonitor)
		}
	}
	return
}
