package log

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newSys(options Options) (log Logger, err error) {
	createlogfile(options.Filename)
	var allCore []zapcore.Core
	hook := lumberjack.Logger{
		Filename:   options.Filename, //日志文件路径
		MaxSize:    2,                //每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,               //最多保留备份个数
		MaxAge:     7,                //文件最多保存多少天
		Compress:   false,            //是否压缩 disabled by default
		LocalTime:  true,             //使用本地时间
	}
	var level zapcore.Level
	switch options.Loglevel {
	case DebugLevel:
		level = zap.DebugLevel
	case InfoLevel:
		level = zap.InfoLevel
	case WarnLevel:
		level = zap.WarnLevel
	case ErrorLevel:
		level = zap.ErrorLevel
	case PanicLevel:
		level = zap.PanicLevel
	case FatalLevel:
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	fileWriter := zapcore.AddSync(&hook)
	consoleDebugging := zapcore.Lock(os.Stdout)
	var encoderConfig zapcore.EncoderConfig
	timeFormat := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
	}
	_, err = os.OpenFile(options.FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	if options.Debugmode {
		//重新生成文件
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = timeFormat
		allCore = append(allCore, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleDebugging, level))
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = timeFormat
	}
	allCore = append(allCore, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), fileWriter, level))
	core := zapcore.NewTee(allCore...)
	tlog := zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(options.Loglayer))
	log = &Logger{
		tlog: tlog,
		log:  tlog.Sugar(),
	}
	return
}

type Logger struct {
	tlog *zap.Logger
	log  *zap.SugaredLogger
}

func FieldTozapField(fields ...Field) (fds []zap.Field) {
	fds = make([]zap.Field, 0)
	for _, v := range fields {
		field := zap.Field{Key: v.Key}

		switch v.Value.(type) {
		case []byte:
			field.Type = zapcore.BinaryType
			field.Interface = v.Value
		case bool:
			field.Type = zapcore.BoolType
			if v.Value.(bool) {
				field.Integer = 1
			} else {
				field.Integer = 0
			}
		case byte:
			field.Type = zapcore.Uint8Type
			field.Integer = int64(v.Value.(byte))
		case time.Duration:
			field.Type = zapcore.DurationType
			field.Integer = int64(v.Value.(time.Duration))
		case float64:
			field.Type = zapcore.Float64Type
			field.Integer = v.Value.(int64)
		case float32:
			field.Type = zapcore.Float32Type
			field.Integer = v.Value.(int64)
		case int64:
			field.Type = zapcore.Int64Type
			field.Integer = v.Value.(int64)
		case uint64:
			field.Type = zapcore.Uint64Type
			field.Integer = int64(v.Value.(uint64))
		case int32:
			field.Type = zapcore.Int32Type
			field.Integer = int64(v.Value.(int32))
		case uint32:
			field.Type = zapcore.Uint64Type
			field.Integer = int64(v.Value.(uint32))
		case int16:
			field.Type = zapcore.Int16Type
			field.Integer = int64(v.Value.(int16))
		case uint16:
			field.Type = zapcore.Uint64Type
			field.Integer = int64(v.Value.(uint16))
		case int8:
			field.Type = zapcore.Int8Type
			field.Integer = int64(v.Value.(int8))
		case string:
			field.Type = zapcore.StringType
			field.String = v.Value.(string)
		case LogStrut:
			field.Type = zapcore.StringType
			field.String = v.Value.(LogStrut).ToString()
		default:
			field.Type = zapcore.UnknownType
			field.Interface = v.Value
		}
		fds = append(fds, field)
	}
	return
}

func (this *Logger) Debug(msg string, fields ...Field) {
	this.tlog.Debug(msg, FieldTozapField(fields...)...)
}
func (this *Logger) Info(msg string, fields ...Field) {
	this.tlog.Info(msg, FieldTozapField(fields...)...)
}
func (this *Logger) Warn(msg string, fields ...Field) {
	this.tlog.Warn(msg, FieldTozapField(fields...)...)
}
func (this *Logger) Error(msg string, fields ...Field) {
	this.tlog.Error(msg, FieldTozapField(fields...)...)
}
func (this *Logger) Panic(msg string, fields ...Field) {
	this.tlog.Panic(msg, FieldTozapField(fields...)...)
}
func (this *Logger) Fatal(msg string, fields ...Field) {
	this.tlog.Fatal(msg, FieldTozapField(fields...)...)
}
func (this *Logger) Debugf(format string, a ...interface{}) {
	this.log.Debugf(format, a...)
}
func (this *Logger) Infof(format string, a ...interface{}) {
	this.log.Infof(format, a...)
}
func (this *Logger) Warnf(format string, a ...interface{}) {
	this.log.Warnf(format, a...)
}
func (this *Logger) Errorf(format string, a ...interface{}) {
	this.log.Errorf(format, a...)
}
func (this *Logger) Panicf(format string, a ...interface{}) {
	this.log.Panicf(format, a...)
}
func (this *Logger) Fatalf(format string, a ...interface{}) {
	this.log.Fatalf(format, a...)
}
