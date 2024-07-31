package log

import (
	"fmt"
	"os"
	"time"

	"github.com/liwei1dao/lego/utils/pools"
)

func newSys(options *Options) (sys *Logger, err error) {
	hook := LogFileOut{
		Filename:   options.FileName,                               //日志文件路径
		MaxAge:     options.MaxAgeTime,                             //备份日志保存天数
		CupTime:    time.Duration(options.CupTimeTime) * time.Hour, //日志切割间隔时间
		Compress:   options.Compress,                               //是否压缩 disabled by default
		MaxBackups: options.MaxBackups,                             //最大备份数
		LocalTime:  true,                                           //使用本地时间
	}
	if !options.IsDebug {
		if err = hook.openNew(); err != nil {
			return
		}
	}
	out := make(writeTree, 0, 2)
	out = append(out, AddSync(&hook))
	if options.IsDebug {
		out = append(out, Lock(os.Stdout))
	}
	sys = &Logger{
		config:     NewDefEncoderConfig(),
		formatter:  NewConsoleEncoder(),
		out:        out,
		level:      options.Loglevel,
		addCaller:  options.ReportCaller,
		callerSkip: options.CallerSkip,
		addStack:   FatalLevel,
	}
	return
}

type Logger struct {
	config     *EncoderConfig //编码配置
	level      LevelEnabler   //日志输出级别
	formatter  Formatter      //日志格式化
	name       string         //日志标签
	out        IWrite         //日志输出
	addCaller  LevelEnabler   //是否打印堆栈信息
	addStack   LevelEnabler   //堆栈信息输出级别
	callerSkip int            //堆栈输出深度
}

func (this *Logger) Clone(name string, skip int) ILogger {
	return &Logger{
		config:     this.config,
		formatter:  this.formatter,
		name:       name,
		out:        this.out,
		level:      this.level,
		addCaller:  this.addCaller,
		callerSkip: skip,
		addStack:   this.addStack,
	}
}
func (this *Logger) SetName(name string) {
	this.name = name
}
func (this *Logger) Enabled(lvl Loglevel) bool {
	return this.level.Enabled(lvl)
}
func (this *Logger) Debug(msg string, args ...Field) {
	this.Log(DebugLevel, msg, args...)
}
func (this *Logger) Info(msg string, args ...Field) {
	this.Log(InfoLevel, msg, args...)
}
func (this *Logger) Print(msg string, args ...Field) {
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
	os.Exit(1)
}
func (this *Logger) Log(level Loglevel, msg string, args ...Field) {
	if this.level.Enabled(level) {
		this.log(level, msg, args...)
	}
}
func (this *Logger) Debugf(format string, args ...interface{}) {
	this.Logf(DebugLevel, format, args...)
}
func (this *Logger) Infof(format string, args ...interface{}) {
	this.Logf(InfoLevel, format, args...)
}
func (this *Logger) Printf(format string, args ...interface{}) {
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
	os.Exit(1)
}
func (this *Logger) Panicf(format string, args ...interface{}) {
	this.Logf(PanicLevel, format, args...)
}
func (this *Logger) Logf(level Loglevel, format string, args ...interface{}) {
	if this.level.Enabled(level) {
		this.log(level, fmt.Sprintf(format, args...))
	}
}
func (this *Logger) Debugln(args ...interface{}) {
	this.Logln(DebugLevel, args...)
}
func (this *Logger) Infoln(args ...interface{}) {
	this.Logln(InfoLevel, args...)
}
func (this *Logger) Println(args ...interface{}) {
	this.Logln(InfoLevel, args...)
}
func (this *Logger) Warnln(args ...interface{}) {
	this.Logln(WarnLevel, args...)
}
func (this *Logger) Errorln(args ...interface{}) {
	this.Logln(ErrorLevel, args...)
}
func (this *Logger) Fatalln(args ...interface{}) {
	this.Logln(FatalLevel, args...)
	os.Exit(1)
}
func (this *Logger) Panicln(args ...interface{}) {
	this.Logln(PanicLevel, args...)
}
func (this *Logger) Logln(level Loglevel, args ...interface{}) {
	if this.level.Enabled(level) {
		this.log(level, this.sprintlnn(args...))
	}
}

func (this *Logger) log(level Loglevel, msg string, args ...Field) {
	entry := this.check(level, msg, args...)
	this.write(entry)
	if level <= PanicLevel {
		panic(entry)
	}
	putEntry(entry)
}

func (this *Logger) check(level Loglevel, msg string, args ...Field) (entry *Entry) {
	entry = getEntry()
	entry.Name = this.name
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg
	entry.WithFields(args...)
	addStack := this.addStack.Enabled(level)
	addCaller := this.addCaller.Enabled(level)
	if !addCaller && !addStack {
		return
	}
	stackDepth := stacktraceFirst
	if addStack {
		stackDepth = stacktraceFull
	}
	stack := captureStacktrace(this.callerSkip+callerSkipOffset, stackDepth)
	defer stack.Free()
	if stack.Count() == 0 {
		if addCaller {
			if entry.Err != "" {
				entry.Err = entry.Err + ",error: failed to get caller"
			} else {
				entry.Err = "error:failed to get caller"
			}
		}
		return
	}
	frame, more := stack.Next()
	if addCaller {
		entry.Caller.Defined = frame.PC != 0
		entry.Caller.PC = frame.PC
		entry.Caller.File = frame.File
		entry.Caller.Line = frame.Line
		entry.Caller.Function = frame.Function
		entry.Caller.Stack = ""
	}
	if addStack {
		buffer := pools.BufferPoolGet()
		defer buffer.Free()
		stackfmt := newStackFormatter(buffer)
		stackfmt.FormatFrame(frame)
		if more {
			stackfmt.FormatStack(stack)
		}
		entry.Caller.Stack = buffer.String()
	}
	return
}

func (this *Logger) write(entry *Entry) {
	buf, err := this.formatter.Format(this.config, entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}
	err = this.out.WriteTo(buf.Bytes())
	buf.Free()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain write, %v\n", err)
		return
	}
	if entry.Level < ErrorLevel {
		this.out.Sync()
	}
	return
}

func (this *Logger) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
