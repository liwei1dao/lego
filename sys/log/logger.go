package log

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

func newSys(options *Options) (sys *Logger, err error) {
	sys = &Logger{
		Out:          os.Stderr,
		Level:        options.Loglevel,
		Skip:         options.CallerSkip,
		ExitFunc:     os.Exit,
		ReportCaller: options.ReportCaller,
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	fileSuffix := time.Now().In(cstSh).Format("2006-01-02") + ".log"

	writer, _ := rotatelogs.New(
		options.FileName+"-"+fileSuffix,
		rotatelogs.WithLinkName(options.FileName+".log"),
		rotatelogs.WithMaxAge(time.Duration(options.MaxAgeTime)*time.Hour*24),      //设置保存时间
		rotatelogs.WithRotationTime(time.Duration(options.RotationTime)*time.Hour), //设置日志分割的时间，隔多久分割一次
	)
	sys.Out = writer
	if options.Encoder == TextEncoder {
		sys.Formatter = &TextFormatter{
			TimestampFormat: "2006-01-02 15:03:04",
		}
	} else {
		sys.Formatter = &JSONFormatter{
			TimestampFormat: "2006-01-02 15:03:04",
		}
	}
	return
}

type Logger struct {
	Out          io.Writer
	Formatter    Formatter
	Level        Loglevel
	Skip         int
	mu           MutexWrap
	entryPool    sync.Pool
	ReportCaller bool
	ExitFunc     exitFunc
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger, logger.Skip)
}
func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
	logger.entryPool.Put(entry)
}

func (this *Logger) Clone(skip int) ILog {
	return NewEntry(this, skip)
}

func (logger *Logger) WithField(key string, value interface{}) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithField(key, value)
}
func (logger *Logger) WithFields(fields ...Field) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields...)
}
func (logger *Logger) WithError(err error) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}
func (logger *Logger) WithContext(ctx context.Context) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithContext(ctx)
}
func (logger *Logger) WithTime(t time.Time) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithTime(t)
}
func (logger *Logger) Exit(code int) {
	runHandlers()
	if logger.ExitFunc == nil {
		logger.ExitFunc = os.Exit
	}
	logger.ExitFunc(code)
}
func (logger *Logger) IsLevelEnabled(level Loglevel) bool {
	return logger.Level >= level
}

func (this *Logger) Debug(msg string, args ...Field) {
	this.Log(DebugLevel, msg, args...)
}
func (this *Logger) Info(msg string, args ...Field) {
	this.Log(InfoLevel, msg, args...)
}
func (this *Logger) Warn(msg string, args ...Field) {
	this.Log(WarnLevel, msg, args...)
}
func (this *Logger) Error(msg string, args ...Field) {
	this.Log(ErrorLevel, msg, args...)
}
func (this *Logger) Panic(msg string, args ...Field) {
	this.Log(PanicLevel, msg, args...)
}
func (this *Logger) Fatal(msg string, args ...Field) {
	this.Log(FatalLevel, msg, args...)
}
func (this *Logger) Log(level Loglevel, msg string, args ...Field) {
	this.WithFields(args...).Log(level, msg)
}
func (this *Logger) Debugf(format string, args ...interface{}) {
	this.Logf(DebugLevel, format, args...)
}
func (this *Logger) Infof(format string, args ...interface{}) {
	this.Logf(InfoLevel, format, args...)
}
func (this *Logger) Warnf(format string, args ...interface{}) {
	this.Logf(WarnLevel, format, args...)
}
func (this *Logger) Errorf(format string, args ...interface{}) {
	this.Logf(ErrorLevel, format, args...)
}
func (this *Logger) Fatalf(format string, args ...interface{}) {
	this.Logf(FatalLevel, format, args...)
}
func (this *Logger) Panicf(format string, args ...interface{}) {
	this.Logf(PanicLevel, format, args...)
}
func (logger *Logger) Logf(level Loglevel, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Logf(level, format, args...)
		logger.releaseEntry(entry)
	}
}
func (this *Logger) Debugln(args ...interface{}) {
	this.Logln(DebugLevel, args...)
}
func (this *Logger) Infoln(args ...interface{}) {
	this.Logln(InfoLevel, args...)
}
func (this *Logger) Println(args ...interface{}) {
	entry := this.newEntry()
	entry.Println(args...)
	this.releaseEntry(entry)
}
func (this *Logger) Warnln(args ...interface{}) {
	this.Logln(WarnLevel, args...)
}
func (this *Logger) Errorln(args ...interface{}) {
	this.Logln(ErrorLevel, args...)
}
func (this *Logger) Fatalln(args ...interface{}) {
	this.Logln(FatalLevel, args...)
}
func (this *Logger) Panicln(args ...interface{}) {
	this.Logln(PanicLevel, args...)
}
func (this *Logger) Logln(level Loglevel, args ...interface{}) {
	if this.IsLevelEnabled(level) {
		entry := this.newEntry()
		entry.Logln(level, args...)
		this.releaseEntry(entry)
	}
}
