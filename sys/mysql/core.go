package mysql

import "gorm.io/gorm"

/*
系统描述:mgo数据库驱动系统
*/
type (
	ISys interface {
		Table(tName string) (tx *gorm.DB)
		Begin() (tx *gorm.DB)
		Raw(sql string, values ...interface{}) (tx *gorm.DB)
		Exec(sql string, values ...interface{}) (tx *gorm.DB)
		CreateTable(tName string, model any) (err error)
		FindOne(tName string, model any, where map[string]interface{}) (err error)
		Find(tName string, models any, where map[string]interface{}) (err error)
		Insert(tName string, model any) (err error)
		Save(tName string, model any) (err error)
		Update(tName string, where, change map[string]interface{}) (err error)
		Delete(tName string, model any, where map[string]interface{}) (err error)
		DropTable(tName string) (err error)
	}
)

var (
	defsys         ISys
	ErrNoDocuments = gorm.ErrRecordNotFound
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func Table(tName string) (tx *gorm.DB) {
	return defsys.Table(tName)
}
func Exec(sql string, values ...interface{}) (tx *gorm.DB) {
	return defsys.Exec(sql, values...)
}
func Raw(sql string, values ...interface{}) (tx *gorm.DB) {
	return defsys.Raw(sql, values...)
}
func CreateTable(tName string, model any) (err error) {
	return defsys.CreateTable(tName, model)
}
func FindOne(tName string, model any, where map[string]interface{}) (err error) {
	return defsys.FindOne(tName, model, where)
}
func Find(tName string, models any, where map[string]interface{}) (err error) {
	return defsys.Find(tName, models, where)
}
func Insert(tName string, model any) (err error) {
	return defsys.Insert(tName, model)
}
func Update(tName string, where, change map[string]interface{}) (err error) {
	return defsys.Update(tName, where, change)
}
func Delete(tName string, model any, where map[string]interface{}) (err error) {
	return defsys.Delete(tName, model, where)
}

func DropTable(tName string) (err error) {
	return defsys.DropTable(tName)
}
