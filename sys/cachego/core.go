package cachego

import (
	"io"
	"time"
)

/*
系统描述:本地缓存系统 可保存本地文件以及读取本地文件
*/

type (
	ISys interface {
		Flush()
		ItemCount() int
		Items() map[string]Item
		SaveFile(fname string) error
		Save(w io.Writer) (err error)
		Load(r io.Reader) error
		LoadFile(fname string) error
		Add(k string, x interface{}, d time.Duration) error
		Set(k string, x interface{}, d time.Duration)
		OnEvicted(f func(string, interface{}))
		Increment(k string, n int64) error
		IncrementFloat(k string, n float64) error
		IncrementInt(k string, n int) (int, error)
		IncrementInt8(k string, n int8) (int8, error)
		IncrementInt16(k string, n int16) (int16, error)
		IncrementInt32(k string, n int32) (int32, error)
		IncrementInt64(k string, n int64) (int64, error)
		IncrementUint(k string, n uint) (uint, error)
		IncrementUintptr(k string, n uintptr) (uintptr, error)
		IncrementUint8(k string, n uint8) (uint8, error)
		IncrementUint16(k string, n uint16) (uint16, error)
		IncrementUint32(k string, n uint32) (uint32, error)
		IncrementUint64(k string, n uint64) (uint64, error)
		IncrementFloat32(k string, n float32) (float32, error)
		IncrementFloat64(k string, n float64) (float64, error)
		Decrement(k string, n int64) error
		DecrementFloat(k string, n float64) error
		DecrementInt(k string, n int) (int, error)
		DecrementInt8(k string, n int8) (int8, error)
		DecrementInt16(k string, n int16) (int16, error)
		DecrementInt32(k string, n int32) (int32, error)
		DecrementInt64(k string, n int64) (int64, error)
		DecrementUint(k string, n uint) (uint, error)
		DecrementUintptr(k string, n uintptr) (uintptr, error)
		DecrementUint8(k string, n uint8) (uint8, error)
		DecrementUint16(k string, n uint16) (uint16, error)
		DecrementUint32(k string, n uint32) (uint32, error)
		DecrementUint64(k string, n uint64) (uint64, error)
		DecrementFloat32(k string, n float32) (float32, error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

func Flush() {
	defsys.Flush()
}
func ItemCount() int {
	return defsys.ItemCount()
}
func Items() map[string]Item {
	return defsys.Items()
}
func SaveFile(fname string) error {
	return defsys.SaveFile(fname)
}
func Save(w io.Writer) (err error) {
	return defsys.Save(w)
}
func Load(r io.Reader) error {
	return defsys.Load(r)
}
func LoadFile(fname string) error {
	return defsys.LoadFile(fname)
}
func Add(k string, x interface{}, d time.Duration) error {
	return defsys.Add(k, x, d)
}
func Set(k string, x interface{}, d time.Duration) {
	defsys.Set(k, x, d)
}
func OnEvicted(f func(string, interface{})) {
	defsys.OnEvicted(f)
}
func Increment(k string, n int64) error {
	return defsys.Increment(k, n)
}
func IncrementFloat(k string, n float64) error {
	return defsys.IncrementFloat(k, n)
}
func IncrementInt(k string, n int) (int, error) {
	return defsys.IncrementInt(k, n)
}
func IncrementInt8(k string, n int8) (int8, error) {
	return defsys.IncrementInt8(k, n)
}
func IncrementInt16(k string, n int16) (int16, error) {
	return defsys.IncrementInt16(k, n)
}
func IncrementInt32(k string, n int32) (int32, error) {
	return defsys.IncrementInt32(k, n)
}
func IncrementInt64(k string, n int64) (int64, error) {
	return defsys.IncrementInt64(k, n)
}
func IncrementUint(k string, n uint) (uint, error) {
	return defsys.IncrementUint(k, n)
}
func IncrementUintptr(k string, n uintptr) (uintptr, error) {
	return defsys.IncrementUintptr(k, n)
}
func IncrementUint8(k string, n uint8) (uint8, error) {
	return defsys.IncrementUint8(k, n)
}
func IncrementUint16(k string, n uint16) (uint16, error) {
	return defsys.IncrementUint16(k, n)
}
func IncrementUint32(k string, n uint32) (uint32, error) {
	return defsys.IncrementUint32(k, n)
}
func IncrementUint64(k string, n uint64) (uint64, error) {
	return defsys.IncrementUint64(k, n)
}
func IncrementFloat32(k string, n float32) (float32, error) {
	return defsys.IncrementFloat32(k, n)
}
func IncrementFloat64(k string, n float64) (float64, error) {
	return defsys.IncrementFloat64(k, n)
}
func Decrement(k string, n int64) error {
	return defsys.Decrement(k, n)
}
func DecrementFloat(k string, n float64) error {
	return defsys.DecrementFloat(k, n)
}
func DecrementInt(k string, n int) (int, error) {
	return defsys.DecrementInt(k, n)
}
func DecrementInt8(k string, n int8) (int8, error) {
	return defsys.DecrementInt8(k, n)
}
func DecrementInt16(k string, n int16) (int16, error) {
	return defsys.DecrementInt16(k, n)
}
func DecrementInt32(k string, n int32) (int32, error) {
	return defsys.DecrementInt32(k, n)
}
func DecrementInt64(k string, n int64) (int64, error) {
	return defsys.DecrementInt64(k, n)
}
func DecrementUint(k string, n uint) (uint, error) {
	return defsys.DecrementUint(k, n)
}
func DecrementUintptr(k string, n uintptr) (uintptr, error) {
	return defsys.DecrementUintptr(k, n)
}
func DecrementUint8(k string, n uint8) (uint8, error) {
	return defsys.DecrementUint8(k, n)
}
func DecrementUint16(k string, n uint16) (uint16, error) {
	return defsys.DecrementUint16(k, n)
}
func DecrementUint32(k string, n uint32) (uint32, error) {
	return defsys.DecrementUint32(k, n)
}
func DecrementUint64(k string, n uint64) (uint64, error) {
	return defsys.DecrementUint64(k, n)
}
func DecrementFloat32(k string, n float32) (float32, error) {
	return defsys.DecrementFloat32(k, n)
}
