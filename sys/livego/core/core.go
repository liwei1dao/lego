package core

import (
	"fmt"
	"net"
	"sync"

	"github.com/liwei1dao/lego/sys/log"
)

var (
	EmptyID = ""
)

var (
	PUBLISH = "publish"
	PLAY    = "play"
)

const (
	maxQueueNum           = 1024
	SAVE_STATICS_INTERVAL = 5000
)
const (
	TAG_AUDIO          = 8   //音频
	TAG_VIDEO          = 9   //视屏
	TAG_SCRIPTDATAAMF0 = 18  //AMF0
	TAG_SCRIPTDATAAMF3 = 0xf //AMF3
)

const (
	SOUND_MP3                   = 2
	SOUND_NELLYMOSER_16KHZ_MONO = 4
	SOUND_NELLYMOSER_8KHZ_MONO  = 5
	SOUND_NELLYMOSER            = 6
	SOUND_ALAW                  = 7
	SOUND_MULAW                 = 8
	SOUND_AAC                   = 10
	SOUND_SPEEX                 = 11

	SOUND_5_5Khz = 0
	SOUND_11Khz  = 1
	SOUND_22Khz  = 2
	SOUND_44Khz  = 3

	SOUND_8BIT  = 0
	SOUND_16BIT = 1

	SOUND_MONO   = 0
	SOUND_STEREO = 1

	AAC_SEQHDR = 0
	AAC_RAW    = 1
)

//Tag 类型
const (
	AVC_SEQHDR = 0
	AVC_NALU   = 1
	AVC_EOS    = 2

	FRAME_KEY   = 1
	FRAME_INTER = 2

	VIDEO_H264 = 7
)

type PacketHeader interface {
}

type AudioPacketHeader interface {
	PacketHeader
	SoundFormat() uint8
	AACPacketType() uint8
}
type VideoPacketHeader interface {
	PacketHeader
	IsKeyFrame() bool
	IsSeq() bool
	CodecID() uint8
	CompositionTime() int32
}

type Alive interface {
	Alive() bool
}

type Closer interface {
	Info() Info
	Close(error)
}

type CalcTime interface {
	CalcBaseTimestamp()
}

type ReadCloser interface {
	Closer
	Alive
	Read(*Packet) error
}

type WriteCloser interface {
	Closer
	Alive
	CalcTime
	Write(*Packet) error
}

type Handler interface {
	GetStreams() *sync.Map
	CheckAlive()
	HandleReader(ReadCloser)
	HandleWriter(WriteCloser)
}

type GetSys interface {
	Sys() ISys
	Log() log.ILogger
}

type GetInFo interface {
	GetInfo() (string, string, string)
}

type GetWriter interface {
	GetWriter(Info) WriteCloser
}

type StreamReadWriteCloser interface {
	GetSys
	GetInFo
	Close(error)
	Write(ChunkStream) error
	Read(c *ChunkStream) error
}

type ISys interface {
	ISysOptions
	Handler
	IStaticPushMnager
	IRoomsManager
	GetRtmpServer() IRtmpServer
}

type ISysOptions interface {
	GetAppname() string
	GetHls() bool
	GetHLSAddr() string
	GetHLSKeepAfterEnd() bool
	GetUseHlsHttps() bool
	GetHlsServerCrt() string
	GetHlsServerKey() string
	GetFlv() bool
	GetHTTPFLVAddr() string
	GetApi() bool
	GetApiAddr() string
	GetJWTSecret() string
	GetJWTAlgorithm() string
	GetStaticPush() []string
	GetRTMPAddr() string
	GetRTMPNoAuth() bool
	GetFLVDir() string
	GetFLVArchive() bool
	GetEnableTLSVerify() bool //是否开启TLS验证
	GetTimeout() int          //单位秒
	GetConnBuffSzie() int     //连接对象读写缓存
	GetDebug() bool           //日志是否开启
}
type IRoomsManager interface {
	SetKey(channel string) (key string, err error)
	GetKey(channel string) (newKey string, err error)
	GetChannel(key string) (channel string, err error)
	DeleteChannel(channel string) (ok bool)
	DeleteKey(key string) (ok bool)
}

type IStaticPushMnager interface {
	GetAndCreateStaticPushObject(rtmpurl string) *StaticPush
	GetStaticPushObject(rtmpurl string) (*StaticPush, error)
	ReleaseStaticPushObject(rtmpurl string)
}
type IRtmpServer interface {
	Serve(listener net.Listener) (err error)
}
type IApiServer interface {
}

type IHlsServer interface {
	GetWriter
}
type IHttpFlvServer interface {
}

type Packet struct {
	IsAudio    bool
	IsVideo    bool
	IsMetadata bool
	TimeStamp  uint32 // dts
	StreamID   uint32
	Header     PacketHeader
	Data       []byte
}

type Info struct {
	Key   string
	URL   string
	UID   string
	Inter bool
}

func (info Info) IsInterval() bool {
	return info.Inter
}

func (info Info) String() string {
	return fmt.Sprintf("<key: %s, URL: %s, UID: %s, Inter: %v>",
		info.Key, info.URL, info.UID, info.Inter)
}
