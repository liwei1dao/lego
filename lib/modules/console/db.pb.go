// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.4
// source: db.proto

package console

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

//验证码类型
type CaptchaType int32

const (
	CaptchaType_LoginCaptcha CaptchaType = 0 //登录验证码
	CaptchaType_BindCaptcha  CaptchaType = 1 //绑定验证码
)

// Enum value maps for CaptchaType.
var (
	CaptchaType_name = map[int32]string{
		0: "LoginCaptcha",
		1: "BindCaptcha",
	}
	CaptchaType_value = map[string]int32{
		"LoginCaptcha": 0,
		"BindCaptcha":  1,
	}
)

func (x CaptchaType) Enum() *CaptchaType {
	p := new(CaptchaType)
	*p = x
	return p
}

func (x CaptchaType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CaptchaType) Descriptor() protoreflect.EnumDescriptor {
	return file_db_proto_enumTypes[0].Descriptor()
}

func (CaptchaType) Type() protoreflect.EnumType {
	return &file_db_proto_enumTypes[0]
}

func (x CaptchaType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CaptchaType.Descriptor instead.
func (CaptchaType) EnumDescriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{0}
}

//用户角色
type UserRole int32

const (
	UserRole_Guester   UserRole = 0 //游客
	UserRole_Operator  UserRole = 1 //运营
	UserRole_Developer UserRole = 2 //开发
	UserRole_Master    UserRole = 3 //超级管理员
)

// Enum value maps for UserRole.
var (
	UserRole_name = map[int32]string{
		0: "Guester",
		1: "Operator",
		2: "Developer",
		3: "Master",
	}
	UserRole_value = map[string]int32{
		"Guester":   0,
		"Operator":  1,
		"Developer": 2,
		"Master":    3,
	}
)

func (x UserRole) Enum() *UserRole {
	p := new(UserRole)
	*p = x
	return p
}

func (x UserRole) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserRole) Descriptor() protoreflect.EnumDescriptor {
	return file_db_proto_enumTypes[1].Descriptor()
}

func (UserRole) Type() protoreflect.EnumType {
	return &file_db_proto_enumTypes[1]
}

func (x UserRole) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserRole.Descriptor instead.
func (UserRole) EnumDescriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{1}
}

type DB_UserData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint32   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty" bson:"_id"`
	PhonOrEmail string   `protobuf:"bytes,2,opt,name=PhonOrEmail,proto3" json:"PhonOrEmail,omitempty"`
	Password    string   `protobuf:"bytes,3,opt,name=Password,proto3" json:"Password,omitempty"`
	NickName    string   `protobuf:"bytes,4,opt,name=NickName,proto3" json:"NickName,omitempty"`
	HeadUrl     string   `protobuf:"bytes,5,opt,name=HeadUrl,proto3" json:"HeadUrl,omitempty"`
	UserRole    UserRole `protobuf:"varint,6,opt,name=UserRole,proto3,enum=UserRole" json:"UserRole,omitempty"`
	Token       string   `protobuf:"bytes,7,opt,name=Token,proto3" json:"Token,omitempty"`
}

func (x *DB_UserData) Reset() {
	*x = DB_UserData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DB_UserData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DB_UserData) ProtoMessage() {}

func (x *DB_UserData) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DB_UserData.ProtoReflect.Descriptor instead.
func (*DB_UserData) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{0}
}

func (x *DB_UserData) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DB_UserData) GetPhonOrEmail() string {
	if x != nil {
		return x.PhonOrEmail
	}
	return ""
}

func (x *DB_UserData) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *DB_UserData) GetNickName() string {
	if x != nil {
		return x.NickName
	}
	return ""
}

func (x *DB_UserData) GetHeadUrl() string {
	if x != nil {
		return x.HeadUrl
	}
	return ""
}

func (x *DB_UserData) GetUserRole() UserRole {
	if x != nil {
		return x.UserRole
	}
	return UserRole_Guester
}

func (x *DB_UserData) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

//用户缓存
type Cache_UserData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Db_UserData *DB_UserData `protobuf:"bytes,1,opt,name=Db_UserData,json=DbUserData,proto3" json:"Db_UserData,omitempty"` //用户信息
	IsOnLine    bool         `protobuf:"varint,2,opt,name=IsOnLine,proto3" json:"IsOnLine,omitempty"`                      //是否在线
}

func (x *Cache_UserData) Reset() {
	*x = Cache_UserData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cache_UserData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cache_UserData) ProtoMessage() {}

func (x *Cache_UserData) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cache_UserData.ProtoReflect.Descriptor instead.
func (*Cache_UserData) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{1}
}

func (x *Cache_UserData) GetDb_UserData() *DB_UserData {
	if x != nil {
		return x.Db_UserData
	}
	return nil
}

func (x *Cache_UserData) GetIsOnLine() bool {
	if x != nil {
		return x.IsOnLine
	}
	return false
}

//cpu信息
type CpuInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CPU        int32   `protobuf:"varint,1,opt,name=CPU,proto3" json:"CPU,omitempty"`              //cpu编号
	VendorID   string  `protobuf:"bytes,2,opt,name=VendorID,proto3" json:"VendorID,omitempty"`     //供应商ID
	Family     string  `protobuf:"bytes,3,opt,name=Family,proto3" json:"Family,omitempty"`         //家庭
	Model      string  `protobuf:"bytes,4,opt,name=Model,proto3" json:"Model,omitempty"`           //模型
	Stepping   int32   `protobuf:"varint,5,opt,name=Stepping,proto3" json:"Stepping,omitempty"`    //步进 表示生产工艺较小的改进
	PhysicalID string  `protobuf:"bytes,6,opt,name=PhysicalID,proto3" json:"PhysicalID,omitempty"` //物理ID
	CoreID     string  `protobuf:"bytes,7,opt,name=CoreID,proto3" json:"CoreID,omitempty"`         //核心ID
	Cores      int32   `protobuf:"varint,8,opt,name=Cores,proto3" json:"Cores,omitempty"`          //核心数
	ModelName  string  `protobuf:"bytes,9,opt,name=ModelName,proto3" json:"ModelName,omitempty"`   //模块名
	Mhz        float64 `protobuf:"fixed64,10,opt,name=Mhz,proto3" json:"Mhz,omitempty"`            //兆赫
	CacheSize  int32   `protobuf:"varint,11,opt,name=CacheSize,proto3" json:"CacheSize,omitempty"` //缓存大小
	Flags      string  `protobuf:"bytes,12,opt,name=Flags,proto3" json:"Flags,omitempty"`          //标志
	Microcode  string  `protobuf:"bytes,13,opt,name=Microcode,proto3" json:"Microcode,omitempty"`  //微码
}

func (x *CpuInfo) Reset() {
	*x = CpuInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CpuInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CpuInfo) ProtoMessage() {}

func (x *CpuInfo) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CpuInfo.ProtoReflect.Descriptor instead.
func (*CpuInfo) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{2}
}

func (x *CpuInfo) GetCPU() int32 {
	if x != nil {
		return x.CPU
	}
	return 0
}

func (x *CpuInfo) GetVendorID() string {
	if x != nil {
		return x.VendorID
	}
	return ""
}

func (x *CpuInfo) GetFamily() string {
	if x != nil {
		return x.Family
	}
	return ""
}

func (x *CpuInfo) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *CpuInfo) GetStepping() int32 {
	if x != nil {
		return x.Stepping
	}
	return 0
}

func (x *CpuInfo) GetPhysicalID() string {
	if x != nil {
		return x.PhysicalID
	}
	return ""
}

func (x *CpuInfo) GetCoreID() string {
	if x != nil {
		return x.CoreID
	}
	return ""
}

func (x *CpuInfo) GetCores() int32 {
	if x != nil {
		return x.Cores
	}
	return 0
}

func (x *CpuInfo) GetModelName() string {
	if x != nil {
		return x.ModelName
	}
	return ""
}

func (x *CpuInfo) GetMhz() float64 {
	if x != nil {
		return x.Mhz
	}
	return 0
}

func (x *CpuInfo) GetCacheSize() int32 {
	if x != nil {
		return x.CacheSize
	}
	return 0
}

func (x *CpuInfo) GetFlags() string {
	if x != nil {
		return x.Flags
	}
	return ""
}

func (x *CpuInfo) GetMicrocode() string {
	if x != nil {
		return x.Microcode
	}
	return ""
}

//内存信息
type MemoryInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total       uint64  `protobuf:"varint,1,opt,name=Total,proto3" json:"Total,omitempty"`
	Available   uint64  `protobuf:"varint,2,opt,name=Available,proto3" json:"Available,omitempty"`
	Used        uint64  `protobuf:"varint,3,opt,name=Used,proto3" json:"Used,omitempty"`
	UsedPercent float64 `protobuf:"fixed64,4,opt,name=UsedPercent,proto3" json:"UsedPercent,omitempty"`
	Free        uint64  `protobuf:"varint,5,opt,name=Free,proto3" json:"Free,omitempty"`
	Active      uint64  `protobuf:"varint,6,opt,name=Active,proto3" json:"Active,omitempty"`
	Inactive    uint64  `protobuf:"varint,7,opt,name=Inactive,proto3" json:"Inactive,omitempty"`
	Wired       uint64  `protobuf:"varint,8,opt,name=Wired,proto3" json:"Wired,omitempty"`
	Laundry     uint64  `protobuf:"varint,9,opt,name=Laundry,proto3" json:"Laundry,omitempty"`
}

func (x *MemoryInfo) Reset() {
	*x = MemoryInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemoryInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemoryInfo) ProtoMessage() {}

func (x *MemoryInfo) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemoryInfo.ProtoReflect.Descriptor instead.
func (*MemoryInfo) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{3}
}

func (x *MemoryInfo) GetTotal() uint64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *MemoryInfo) GetAvailable() uint64 {
	if x != nil {
		return x.Available
	}
	return 0
}

func (x *MemoryInfo) GetUsed() uint64 {
	if x != nil {
		return x.Used
	}
	return 0
}

func (x *MemoryInfo) GetUsedPercent() float64 {
	if x != nil {
		return x.UsedPercent
	}
	return 0
}

func (x *MemoryInfo) GetFree() uint64 {
	if x != nil {
		return x.Free
	}
	return 0
}

func (x *MemoryInfo) GetActive() uint64 {
	if x != nil {
		return x.Active
	}
	return 0
}

func (x *MemoryInfo) GetInactive() uint64 {
	if x != nil {
		return x.Inactive
	}
	return 0
}

func (x *MemoryInfo) GetWired() uint64 {
	if x != nil {
		return x.Wired
	}
	return 0
}

func (x *MemoryInfo) GetLaundry() uint64 {
	if x != nil {
		return x.Laundry
	}
	return 0
}

//主机信息
type HostInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HostID               string `protobuf:"bytes,1,opt,name=HostID,proto3" json:"HostID,omitempty"`                              //主机id
	Hostname             string `protobuf:"bytes,2,opt,name=Hostname,proto3" json:"Hostname,omitempty"`                          //主机名称
	Uptime               uint64 `protobuf:"varint,3,opt,name=Uptime,proto3" json:"Uptime,omitempty"`                             //正常运行时间
	BootTime             uint64 `protobuf:"varint,4,opt,name=BootTime,proto3" json:"BootTime,omitempty"`                         //开机时间
	Procs                uint64 `protobuf:"varint,5,opt,name=Procs,proto3" json:"Procs,omitempty"`                               //进程数
	OS                   string `protobuf:"bytes,6,opt,name=OS,proto3" json:"OS,omitempty"`                                      //内核系统 例如:freebsd, linux
	Platform             string `protobuf:"bytes,7,opt,name=Platform,proto3" json:"Platform,omitempty"`                          //操作系统 例如:ubuntu, centos
	PlatformFamily       string `protobuf:"bytes,8,opt,name=PlatformFamily,proto3" json:"PlatformFamily,omitempty"`              //主机系统系列 ex: debian, rhel
	PlatformVersion      string `protobuf:"bytes,9,opt,name=PlatformVersion,proto3" json:"PlatformVersion,omitempty"`            //主机系统版本
	KernelArch           string `protobuf:"bytes,10,opt,name=KernelArch,proto3" json:"KernelArch,omitempty"`                     //Cpu架构
	VirtualizationSystem string `protobuf:"bytes,11,opt,name=VirtualizationSystem,proto3" json:"VirtualizationSystem,omitempty"` //虚拟系统
	VirtualizationRole   string `protobuf:"bytes,12,opt,name=VirtualizationRole,proto3" json:"VirtualizationRole,omitempty"`     //虚拟身份  guest or host
}

func (x *HostInfo) Reset() {
	*x = HostInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HostInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HostInfo) ProtoMessage() {}

func (x *HostInfo) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HostInfo.ProtoReflect.Descriptor instead.
func (*HostInfo) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{4}
}

func (x *HostInfo) GetHostID() string {
	if x != nil {
		return x.HostID
	}
	return ""
}

func (x *HostInfo) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *HostInfo) GetUptime() uint64 {
	if x != nil {
		return x.Uptime
	}
	return 0
}

func (x *HostInfo) GetBootTime() uint64 {
	if x != nil {
		return x.BootTime
	}
	return 0
}

func (x *HostInfo) GetProcs() uint64 {
	if x != nil {
		return x.Procs
	}
	return 0
}

func (x *HostInfo) GetOS() string {
	if x != nil {
		return x.OS
	}
	return ""
}

func (x *HostInfo) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *HostInfo) GetPlatformFamily() string {
	if x != nil {
		return x.PlatformFamily
	}
	return ""
}

func (x *HostInfo) GetPlatformVersion() string {
	if x != nil {
		return x.PlatformVersion
	}
	return ""
}

func (x *HostInfo) GetKernelArch() string {
	if x != nil {
		return x.KernelArch
	}
	return ""
}

func (x *HostInfo) GetVirtualizationSystem() string {
	if x != nil {
		return x.VirtualizationSystem
	}
	return ""
}

func (x *HostInfo) GetVirtualizationRole() string {
	if x != nil {
		return x.VirtualizationRole
	}
	return ""
}

//主机监控信息
type HostMonitor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CpuUsageRate    []float64 `protobuf:"fixed64,1,rep,packed,name=CpuUsageRate,proto3" json:"CpuUsageRate,omitempty"`       //Cpus使用率统计
	MemoryUsageRate []float64 `protobuf:"fixed64,2,rep,packed,name=MemoryUsageRate,proto3" json:"MemoryUsageRate,omitempty"` //内存使用率统计
}

func (x *HostMonitor) Reset() {
	*x = HostMonitor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HostMonitor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HostMonitor) ProtoMessage() {}

func (x *HostMonitor) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HostMonitor.ProtoReflect.Descriptor instead.
func (*HostMonitor) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{5}
}

func (x *HostMonitor) GetCpuUsageRate() []float64 {
	if x != nil {
		return x.CpuUsageRate
	}
	return nil
}

func (x *HostMonitor) GetMemoryUsageRate() []float64 {
	if x != nil {
		return x.MemoryUsageRate
	}
	return nil
}

//集群监控信息
type ClusterMonitor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CpuUsageRate    []float64 `protobuf:"fixed64,1,rep,packed,name=CpuUsageRate,proto3" json:"CpuUsageRate,omitempty"`       //Cpus使用率统计
	MemoryUsageRate []float64 `protobuf:"fixed64,2,rep,packed,name=MemoryUsageRate,proto3" json:"MemoryUsageRate,omitempty"` //内存使用率统计
	GoroutineUsed   []float64 `protobuf:"fixed64,3,rep,packed,name=GoroutineUsed,proto3" json:"GoroutineUsed,omitempty"`     //携程数
	PreWeight       []float64 `protobuf:"fixed64,4,rep,packed,name=PreWeight,proto3" json:"PreWeight,omitempty"`             //权重数
}

func (x *ClusterMonitor) Reset() {
	*x = ClusterMonitor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterMonitor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterMonitor) ProtoMessage() {}

func (x *ClusterMonitor) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterMonitor.ProtoReflect.Descriptor instead.
func (*ClusterMonitor) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{6}
}

func (x *ClusterMonitor) GetCpuUsageRate() []float64 {
	if x != nil {
		return x.CpuUsageRate
	}
	return nil
}

func (x *ClusterMonitor) GetMemoryUsageRate() []float64 {
	if x != nil {
		return x.MemoryUsageRate
	}
	return nil
}

func (x *ClusterMonitor) GetGoroutineUsed() []float64 {
	if x != nil {
		return x.GoroutineUsed
	}
	return nil
}

func (x *ClusterMonitor) GetPreWeight() []float64 {
	if x != nil {
		return x.PreWeight
	}
	return nil
}

var File_db_proto protoreflect.FileDescriptor

var file_db_proto_rawDesc = []byte{
	0x0a, 0x08, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xce, 0x01, 0x0a, 0x0b, 0x44,
	0x42, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x68,
	0x6f, 0x6e, 0x4f, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x50, 0x68, 0x6f, 0x6e, 0x4f, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x4e, 0x69, 0x63, 0x6b,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4e, 0x69, 0x63, 0x6b,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x48, 0x65, 0x61, 0x64, 0x55, 0x72, 0x6c, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x48, 0x65, 0x61, 0x64, 0x55, 0x72, 0x6c, 0x12, 0x25,
	0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x09, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x08, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x5b, 0x0a, 0x0e, 0x43,
	0x61, 0x63, 0x68, 0x65, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x2d, 0x0a,
	0x0b, 0x44, 0x62, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x44, 0x42, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x0a, 0x44, 0x62, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08,
	0x49, 0x73, 0x4f, 0x6e, 0x4c, 0x69, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x49, 0x73, 0x4f, 0x6e, 0x4c, 0x69, 0x6e, 0x65, 0x22, 0xd1, 0x02, 0x0a, 0x07, 0x43, 0x70, 0x75,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x43, 0x50, 0x55, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x03, 0x43, 0x50, 0x55, 0x12, 0x1a, 0x0a, 0x08, 0x56, 0x65, 0x6e, 0x64, 0x6f, 0x72,
	0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x56, 0x65, 0x6e, 0x64, 0x6f, 0x72,
	0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x4d, 0x6f,
	0x64, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x12, 0x1a, 0x0a, 0x08, 0x53, 0x74, 0x65, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x53, 0x74, 0x65, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x12, 0x1e, 0x0a, 0x0a,
	0x50, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x49, 0x44, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x50, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06,
	0x43, 0x6f, 0x72, 0x65, 0x49, 0x44, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x43, 0x6f,
	0x72, 0x65, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x6f, 0x72, 0x65, 0x73, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x43, 0x6f, 0x72, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x4d, 0x6f,
	0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x68, 0x7a, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x4d, 0x68, 0x7a, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x61,
	0x63, 0x68, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x43,
	0x61, 0x63, 0x68, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x46, 0x6c, 0x61, 0x67,
	0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x1c,
	0x0a, 0x09, 0x4d, 0x69, 0x63, 0x72, 0x6f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x4d, 0x69, 0x63, 0x72, 0x6f, 0x63, 0x6f, 0x64, 0x65, 0x22, 0xee, 0x01, 0x0a,
	0x0a, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x54,
	0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x55,
	0x73, 0x65, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x64, 0x50, 0x65, 0x72, 0x63, 0x65,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x55, 0x73, 0x65, 0x64, 0x50, 0x65,
	0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x65, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x46, 0x72, 0x65, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x08, 0x49, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x57, 0x69, 0x72, 0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x57, 0x69,
	0x72, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x4c, 0x61, 0x75, 0x6e, 0x64, 0x72, 0x79, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x4c, 0x61, 0x75, 0x6e, 0x64, 0x72, 0x79, 0x22, 0x8a, 0x03,
	0x0a, 0x08, 0x48, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x48, 0x6f,
	0x73, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x48, 0x6f, 0x73, 0x74,
	0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x48, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x48, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x55, 0x70, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06,
	0x55, 0x70, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x42, 0x6f, 0x6f, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x42, 0x6f, 0x6f, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x63, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x50, 0x72, 0x6f, 0x63, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x4f, 0x53, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x4f, 0x53, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x12, 0x26, 0x0a, 0x0e, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x50, 0x6c,
	0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x28, 0x0a, 0x0f,
	0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x4b, 0x65, 0x72, 0x6e, 0x65, 0x6c,
	0x41, 0x72, 0x63, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x4b, 0x65, 0x72, 0x6e,
	0x65, 0x6c, 0x41, 0x72, 0x63, 0x68, 0x12, 0x32, 0x0a, 0x14, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61,
	0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x2e, 0x0a, 0x12, 0x56, 0x69,
	0x72, 0x74, 0x75, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x6c, 0x65,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x5b, 0x0a, 0x0b, 0x48, 0x6f,
	0x73, 0x74, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x12, 0x22, 0x0a, 0x0c, 0x43, 0x70, 0x75,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x01, 0x52,
	0x0c, 0x43, 0x70, 0x75, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x28, 0x0a,
	0x0f, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x65,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x01, 0x52, 0x0f, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x55, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x65, 0x22, 0xa2, 0x01, 0x0a, 0x0e, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x12, 0x22, 0x0a, 0x0c, 0x43, 0x70,
	0x75, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x01,
	0x52, 0x0c, 0x43, 0x70, 0x75, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x28,
	0x0a, 0x0f, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74,
	0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x01, 0x52, 0x0f, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x55,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x47, 0x6f, 0x72, 0x6f,
	0x75, 0x74, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x64, 0x18, 0x03, 0x20, 0x03, 0x28, 0x01, 0x52,
	0x0d, 0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x50, 0x72, 0x65, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x01, 0x52, 0x09, 0x50, 0x72, 0x65, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x2a, 0x30, 0x0a, 0x0b,
	0x43, 0x61, 0x70, 0x74, 0x63, 0x68, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x0c, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x43, 0x61, 0x70, 0x74, 0x63, 0x68, 0x61, 0x10, 0x00, 0x12, 0x0f, 0x0a,
	0x0b, 0x42, 0x69, 0x6e, 0x64, 0x43, 0x61, 0x70, 0x74, 0x63, 0x68, 0x61, 0x10, 0x01, 0x2a, 0x40,
	0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x47, 0x75,
	0x65, 0x73, 0x74, 0x65, 0x72, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x65, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x72, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x10, 0x03,
	0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x3b, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_db_proto_rawDescOnce sync.Once
	file_db_proto_rawDescData = file_db_proto_rawDesc
)

func file_db_proto_rawDescGZIP() []byte {
	file_db_proto_rawDescOnce.Do(func() {
		file_db_proto_rawDescData = protoimpl.X.CompressGZIP(file_db_proto_rawDescData)
	})
	return file_db_proto_rawDescData
}

var file_db_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_db_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_db_proto_goTypes = []interface{}{
	(CaptchaType)(0),       // 0: CaptchaType
	(UserRole)(0),          // 1: UserRole
	(*DB_UserData)(nil),    // 2: DB_UserData
	(*Cache_UserData)(nil), // 3: Cache_UserData
	(*CpuInfo)(nil),        // 4: CpuInfo
	(*MemoryInfo)(nil),     // 5: MemoryInfo
	(*HostInfo)(nil),       // 6: HostInfo
	(*HostMonitor)(nil),    // 7: HostMonitor
	(*ClusterMonitor)(nil), // 8: ClusterMonitor
}
var file_db_proto_depIdxs = []int32{
	1, // 0: DB_UserData.UserRole:type_name -> UserRole
	2, // 1: Cache_UserData.Db_UserData:type_name -> DB_UserData
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_db_proto_init() }
func file_db_proto_init() {
	if File_db_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_db_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DB_UserData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cache_UserData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CpuInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemoryInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HostInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HostMonitor); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterMonitor); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_db_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_db_proto_goTypes,
		DependencyIndexes: file_db_proto_depIdxs,
		EnumInfos:         file_db_proto_enumTypes,
		MessageInfos:      file_db_proto_msgTypes,
	}.Build()
	File_db_proto = out.File
	file_db_proto_rawDesc = nil
	file_db_proto_goTypes = nil
	file_db_proto_depIdxs = nil
}
