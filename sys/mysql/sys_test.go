package mysql_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/mysql"

	msql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// 测试 系统 直连模式
func Test_mgsql(t *testing.T) {
	dsn := "root:li13451234@tcp(127.0.0.1:3306)/admin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(msql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 迁移 schema
	// 创建表
	db.Table("test").AutoMigrate(&Product{})

	// Create
	db.Table("test").Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.Table("test").First(&product, 1)                 // 根据整型主键查找
	db.Table("test").First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Table("test").Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	db.Table("test").Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Table("test").Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Table("test").Delete(&product, 1)
}

type DBAgent struct {
	AgentId  string `gorm:"id:50;unique"`               // 唯一索引
	AgentKey string `gorm:"key:255;index:idx_age_name"` // 组合索引的一部分
	Addr     string `gorm:"addr:255"`                   // 组合索引的一部分
}

func Test_sys(t *testing.T) {
	if sys, err := mysql.NewSys(mysql.SetMySQLDsn("root:li13451234@tcp(127.0.0.1:3306)/admin?charset=utf8mb4&parseTime=True&loc=Local")); err != nil {
		fmt.Println(err)
	} else {
		// if err = sys.CreateTable("agent", &DBAgent{}); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// if err = sys.Insert("agent", &DBAgent{
		// 	AgentId:  "1001",
		// 	AgentKey: "D3PX7iaNEZN5FdkG2wfb0w==",
		// 	Addr:     "https://baidu.com",
		// }); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// if err = sys.Insert("agent", []*DBAgent{
		// 	{
		// 		AgentId:  "1002",
		// 		AgentKey: "D3PX7iaNEZN5FdkG2wfb0w==",
		// 		Addr:     "https://baidu.com",
		// 	},
		// 	{
		// 		AgentId:  "1003",
		// 		AgentKey: "D3PX7iaNEZN5FdkG2wfb0w==",
		// 		Addr:     "https://baidu.com",
		// 	},
		// }); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// var agent DBAgent
		// if err = sys.Find("agent", &agent, map[string]interface{}{"agent_id": 1001}); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		var agent []*DBAgent = make([]*DBAgent, 0)
		if err = sys.Find("agent", &agent, map[string]interface{}{}); err != nil {
			fmt.Println(err)
			return
		}
		return
	}
}
