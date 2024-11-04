package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	v9 "github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"

	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/sys/log"

	"github.com/liwei1dao/lego/sys/redis"
)

// redis的备份数据
type RedisBackup struct {
	Type  string            `json:"type"`
	Value string            `json:"value"`
	List  []string          `json:"list,omitempty"`
	Set   []string          `json:"set,omitempty"`
	Hash  map[string]string `json:"hash,omitempty"`
}

/*
服务类型:工具库代码
服务描述:集成日志导出工具
*/
var (
	addrs string //数据库地址
	pw    string //密码
	tls   int32
	file  string
)
var logoutCmd = &cobra.Command{
	Use:   "backup",
	Short: "备份数据",
	Run: func(cmd *cobra.Command, args []string) {
		lego.Recover("backup")
		backup()
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "还原数据",
	Run: func(cmd *cobra.Command, args []string) {
		lego.Recover("restore")
		restore()
	},
}

func emptyRun(*cobra.Command, []string) {}

var RootCmd = &cobra.Command{
	Use:   "a11",
	Short: "命令行",
	Long:  "命令行工具",
	Run:   emptyRun,
}

// 初始化自定义cmd
func init() {
	RootCmd.PersistentFlags().StringVarP(&addrs, "add", "a", "127.0.0.1:6379", "Redis的地址")
	RootCmd.PersistentFlags().StringVarP(&pw, "pw", "p", "", "密码")
	RootCmd.PersistentFlags().Int32VarP(&tls, "tls", "t", 0, "tls 开关")
	RootCmd.PersistentFlags().StringVarP(&file, "file", "f", "./redis.json", "写入和读取文件")
	RootCmd.AddCommand(logoutCmd)
	RootCmd.AddCommand(restoreCmd)
}

func main() {
	flag.Parse()
	if err := log.OnInit(nil,
		log.SetFileName("./redis.log"),
		log.SetLoglevel(log.DebugLevel),
		log.SetIsDebug(true)); err != nil {
		panic(fmt.Sprintf("Sys log Init err:%v !", err))
	} else {
		log.Infof("Sys log Init success !")
	}
	Execute()
}

// 执行命令
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

// 输出日志
func backup() {
	var (
		addr   []string
		client v9.UniversalClient
		err    error
	)
	addr = strings.Split(addrs, ",")
	if err = redis.OnInit(nil,
		redis.SetRedisAddr(addr),
		redis.SetRedisPassword(pw),
		redis.SetRedisTLS(tls == 1),
	); err != nil {
		log.Errorln("err:%v", err)
		return
	} else {
		client = redis.GetClient()
		_, err = client.Ping(context.Background()).Result()
		if err != nil {
			log.Errorln("addr:%s pw:%s ttl:%f err:%v", addrs, pw, tls, err)
			return
		}
		data := make(map[string]RedisBackup)
		var cursor uint64
		for {
			keys, newCursor, err := client.Scan(context.Background(), cursor, "*", 0).Result()
			if err != nil {
				log.Errorln("err:%v", err)
				return
			}
			cursor = newCursor
			for _, key := range keys {
				keyType, err := client.Type(context.Background(), key).Result()
				if err != nil {
					log.Errorln("err:%v", err)
					return
				}
				switch keyType {
				case "string":
					value, err := client.Get(context.Background(), key).Result()
					if err != nil {
						log.Errorln("err:%v", err)
						return
					}
					data[key] = RedisBackup{Type: "string", Value: value}
				case "list":
					values, err := client.LRange(context.Background(), key, 0, -1).Result()
					if err != nil {
						log.Errorln("err:%v", err)
						return
					}
					data[key] = RedisBackup{Type: "list", List: values}
				case "set":
					values, err := client.SMembers(context.Background(), key).Result()
					if err != nil {
						log.Errorln("err:%v", err)
						return
					}
					data[key] = RedisBackup{Type: "set", Set: values}
				case "hash":
					values, err := client.HGetAll(context.Background(), key).Result()
					if err != nil {
						log.Errorln("err:%v", err)
						return
					}
					data[key] = RedisBackup{Type: "hash", Hash: values}
				// 可以添加更多类型支持，例如 hash 等
				default:
					log.Panicf("Unsupported key type: %s for key: %s\n", keyType, key)
				}
			}
			if cursor == 0 {
				break
			}
		}

		filebyte, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Errorln("err:%v", err)
			return
		}
		ioutil.WriteFile(file, filebyte, 0644)
		if err != nil {
			log.Errorln("err:%v", err)
			return
		}
	}
}

func restore() {
	var (
		addr   []string
		client v9.UniversalClient
		err    error
	)
	addr = strings.Split(addrs, ",")
	if err = redis.OnInit(nil,
		redis.SetRedisAddr(addr),
		redis.SetRedisPassword(pw),
		redis.SetRedisTLS(tls == 1),
	); err != nil {
		log.Errorln("err:%v", err)
		return
	} else {
		client = redis.GetClient()
		_, err = client.Ping(context.Background()).Result()
		if err != nil {
			log.Errorln("addr:%s pw:%s ttl:%f err:%v", addrs, pw, tls, err)
			return
		}
		filedata, err := ioutil.ReadFile(file)
		if err != nil {
			log.Errorln("err:%v", err)
			return
		}
		var data map[string]RedisBackup
		if err := json.Unmarshal(filedata, &data); err != nil {
			log.Errorln("err:%v", err)
			return
		}
		for key, backup := range data {
			switch backup.Type {
			case "string":
				if err := client.Set(context.Background(), key, backup.Value, 0).Err(); err != nil {
					log.Errorln("err:%v", err)
					return
				}
			case "list":
				if err := client.RPush(context.Background(), key, backup.List).Err(); err != nil {
					log.Errorln("err:%v", err)
					return
				}
			case "set":
				if err := client.SAdd(context.Background(), key, backup.Set).Err(); err != nil {
					log.Errorln("err:%v", err)
					return
				}
			case "hash":
				if err := client.HSet(context.Background(), key, backup.Hash).Err(); err != nil {
					log.Errorln("err:%v", err)
					return
				}
			default:
				log.Panicf("Unsupported key type: %s for key: %s\n", backup.Type, key)
			}
		}

		filebyte, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Errorln("err:%v", err)
			return
		}
		ioutil.WriteFile(file, filebyte, 0644)
		if err != nil {
			log.Errorln("err:%v", err)
			return
		}
	}
}
