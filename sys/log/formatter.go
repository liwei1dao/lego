package log

import (
	"sync"
	"time"

	"github.com/liwei1dao/lego/utils/pools"
)

var _fieldsPool = sync.Pool{
	New: func() interface{} {
		return make(Fields, 0, 6)
	},
}

func getFields() Fields {
	return _fieldsPool.Get().(Fields)
}

func putFields(fields Fields) {
	fields = fields[:0]
	_fieldsPool.Put(fields)
}

func NewDefEncoderConfig() *EncoderConfig {
	return &EncoderConfig{
		TimeKey:          "ts",
		LevelKey:         "level",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		ConsoleSeparator: "\t",
		EncodeTime:       DefTimeEncoder,
		EncodeLevel:      LowercaseLevelEncoder,
		EncodeCaller:     ShortCallerEncoder,
	}
}

type TimeEncoder func(time.Time) string
type LevelEncoder func(Loglevel) string
type CallerEncoder func(EntryCaller) string

func DefTimeEncoder(t time.Time) string {
	return t.Format("2006/01/02 15:04:05.000")
}

func LowercaseLevelEncoder(l Loglevel) string {
	return l.CapitalString()
}

func ShortCallerEncoder(caller EntryCaller) string {
	return caller.TrimmedPath()
}

type EncoderConfig struct {
	MessageKey       string `json:"messageKey" yaml:"messageKey"`
	LevelKey         string `json:"levelKey" yaml:"levelKey"`
	TimeKey          string `json:"timeKey" yaml:"timeKey"`
	CallerKey        string `json:"callerKey" yaml:"callerKey"`
	FunctionKey      string `json:"functionKey" yaml:"functionKey"`
	StacktraceKey    string `json:"stacktraceKey" yaml:"stacktraceKey"`
	ConsoleSeparator string `json:"consoleSeparator" yaml:"consoleSeparator"`
	EncodeTime       TimeEncoder
	EncodeLevel      LevelEncoder
	EncodeCaller     CallerEncoder
}

type Formatter interface {
	Format(config *EncoderConfig, entry *Entry) (*pools.Buffer, error)
}
