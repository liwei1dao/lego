package console

import (
	"fmt"
	reflect "reflect"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/redis"
)

type CacheComp struct {
	cbase.ModuleCompBase
	module IConsole
	redis  redis.IRedisFactory
}

func (this *CacheComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.module = module.(IConsole)
	this.redis, err = redis.NewSys(redis.SetRedisUrl(this.module.Options().GetRedisUrl()))
	return
}

func (this *CacheComp) GetPool() *redis.RedisPool {
	return this.redis.GetPool()
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
	pool := this.redis.GetPool()
	err = pool.GetKeyForValue(Id, &uId)
	return
}

//写入Token
func (this *CacheComp) WriteToken(token string, uId uint32) (err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleToken), token)
	pool := this.redis.GetPool()
	err = pool.SetExKeyForValue(Id, uId, this.module.Options().GetTokenCacheExpirationDate())
	return
}

//清理Token
func (this *CacheComp) CleanToken(token string) (err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleToken), token)
	pool := this.redis.GetPool()
	err = pool.Delete(Id)
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
	pool := this.redis.GetPool()
	result = &Cache_UserData{}
	err = pool.GetKeyForValue(Id, result)
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
	pool := this.redis.GetPool()
	err = pool.SetExKeyForValue(Id, result, this.module.Options().GetUserCacheExpirationDate())
	return
}

//登录用户缓存信息长期驻留
func (this *CacheComp) WriteUserData(data *Cache_UserData) (err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleUsers), data.Db_UserData.Id)
	pool := this.redis.GetPool()
	err = pool.SetKeyForValue(Id, data)
	return
}

//清理用户缓存
func (this *CacheComp) CleanUserData(uid uint32) (err error) {
	Id := fmt.Sprintf(string(Cache_ConsoleUsers), uid)
	pool := this.redis.GetPool()
	err = pool.Delete(Id)
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
	pool := this.redis.GetPool()
	for k, v := range data {
		id := fmt.Sprintf(string(Cache_ConsoleClusterMonitor), k)
		pool.SetListByRPush(id, []interface{}{v})
		if len, err := pool.GetListCount(id); err == nil && len > this.module.Options().GetMonitorTotalTime() {
			pool.GetListByLPop(string(Cache_ConsoleClusterMonitor), v)
		}
	}
}

//添加新的ClusterMonitor
func (this *CacheComp) GetClusterMonitor(sIs string, timeleng int32) (result []*ClusterMonitor, err error) {
	var values []interface{}
	result = make([]*ClusterMonitor, 0)
	id := fmt.Sprintf(string(Cache_ConsoleClusterMonitor), sIs)
	pool := this.redis.GetPool()
	values, err = pool.GetListByLrange(id, 0, timeleng, reflect.TypeOf(&ClusterMonitor{}))
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
	pool := this.redis.GetPool()
	pool.SetListByRPush(string(Cache_ConsoleHostMonitor), []interface{}{data})
	if len, err := pool.GetListCount(string(Cache_ConsoleHostMonitor)); err == nil && len > this.module.Options().GetMonitorTotalTime() {
		pool.GetListByLPop(string(Cache_ConsoleHostMonitor), data)
	}
}

//添加新的HostMonitor
func (this *CacheComp) GetHostMonitor(timeleng int32) (result []*HostMonitor, err error) {
	var values []interface{}
	result = make([]*HostMonitor, 0)
	pool := this.redis.GetPool()
	values, err = pool.GetListByLrange(string(Cache_ConsoleHostMonitor), 0, timeleng, reflect.TypeOf(&HostMonitor{}))
	if values != nil && len(values) > 0 {
		result = make([]*HostMonitor, len(values))
		for i, v := range values {
			result[i] = v.(*HostMonitor)
		}
	}
	return
}
