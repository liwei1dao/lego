package log

import (
	"encoding/json"
	"time"

	"github.com/liwei1dao/lego/utils/pools"
)

func NewConsoleEncoder() Formatter {
	return &ConsoleFormatter{}
}

type ConsoleFormatter struct {
	*EncoderConfig
}

func (this *ConsoleFormatter) Format(config *EncoderConfig, entry *Entry) (*pools.Buffer, error) {
	line := pools.BufferPoolGet()
	isfirst := true
	if config.TimeKey != "" && config.EncodeTime != nil {
		line.AppendString(config.EncodeTime(entry.Time))
		isfirst = false
	}
	if config.LevelKey != "" && config.EncodeLevel != nil {
		if !isfirst {
			line.AppendString(config.ConsoleSeparator)
		}
		isfirst = false
		line.AppendString(config.EncodeLevel(entry.Level))
	}
	if entry.Caller.Defined {
		if config.CallerKey != "" && config.EncodeCaller != nil {
			if !isfirst {
				line.AppendString(config.ConsoleSeparator)
			}
			isfirst = false
			line.AppendString(config.EncodeCaller(entry.Caller))
		}
		if config.FunctionKey != "" {
			if !isfirst {
				line.AppendString(config.ConsoleSeparator)
			}
			isfirst = false
			line.AppendString(entry.Caller.Function)
		}
	}
	if entry.Name != "" {
		if !isfirst {
			line.AppendString(config.ConsoleSeparator)
		}
		isfirst = false
		line.AppendString(entry.Name)
	}

	if config.MessageKey != "" {
		if !isfirst {
			line.AppendString(config.ConsoleSeparator)
		}
		isfirst = false
		line.AppendString(entry.Message)
	}
	for _, v := range entry.Data {
		if !isfirst {
			line.AppendString(config.ConsoleSeparator)
		}
		isfirst = false
		line.AppendString(v.Key)
		line.AppendString(":")
		writetoline(line, v.Value)
	}

	if entry.Caller.Stack != "" && config.StacktraceKey != "" {
		line.AppendString("\n")
		line.AppendString(entry.Caller.Stack)
	}
	line.AppendString("\n")
	return line, nil
}

func writetoline(line *pools.Buffer, v interface{}) {
	switch v := v.(type) {
	case nil:
		line.AppendString("nil")
		break
	case string:
		line.AppendString(v)
		break
	case []byte:
		line.AppendBytes(v)
		break
	case int:
		line.AppendInt(int64(v))
		break
	case int8:
		line.AppendInt(int64(v))
		break
	case int16:
		line.AppendInt(int64(v))
		break
	case int32:
		line.AppendInt(int64(v))
		break
	case int64:
		line.AppendInt(int64(v))
		break
	case uint:
		line.AppendUint(uint64(v))
		break
	case uint8:
		line.AppendUint(uint64(v))
		break
	case uint16:
		line.AppendUint(uint64(v))
		break
	case uint32:
		line.AppendUint(uint64(v))
		break
	case uint64:
		line.AppendUint(uint64(v))
		break
	case float32:
		line.AppendFloat(float64(v), 64)
		break
	case float64:
		line.AppendFloat(v, 64)
		break
	case bool:
		line.AppendBool(v)
		break
	case time.Time:
		line.AppendTime(v, time.RFC3339Nano)
		break
	case time.Duration:
		line.AppendInt(v.Nanoseconds())
		break
	default:
		d, _ := json.Marshal(v)
		line.AppendBytes(d)
		break
	}
}
