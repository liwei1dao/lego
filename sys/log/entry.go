package log

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"time"
)

var ErrorKey = "error"

func NewEntry(logger *Logger, skip int) *Entry {
	return &Entry{
		Logger: logger,
		Skip:   skip,
		Data:   make(Fields, 6),
	}
}

type Entry struct {
	Logger  *Logger
	Data    Fields
	Time    time.Time
	Level   Loglevel
	Skip    int
	Caller  *runtime.Frame
	Message string
	Buffer  *bytes.Buffer
	Context context.Context
	err     string
}

func (entry *Entry) Dup() *Entry {
	data := make(Fields, len(entry.Data))
	for k, v := range entry.Data {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time, Context: entry.Context, err: entry.err}
}
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Field{key, value})
}
func (entry *Entry) WithFields(fields ...Field) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	fieldErr := entry.err
	for _, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v.Value); t != nil {
			switch {
			case t.Kind() == reflect.Func, t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Func:
				isErrField = true
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", v.Key)
			if fieldErr != "" {
				fieldErr = entry.err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			data[v.Key] = v.Value
		}
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time, err: fieldErr, Context: entry.Context}
}
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrorKey, err)
}
func (entry *Entry) WithContext(ctx context.Context) *Entry {
	dataCopy := make(Fields, len(entry.Data))
	for k, v := range entry.Data {
		dataCopy[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: dataCopy, Time: entry.Time, err: entry.err, Context: ctx}
}
func (entry *Entry) WithTime(t time.Time) *Entry {
	dataCopy := make(Fields, len(entry.Data))
	for k, v := range entry.Data {
		dataCopy[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: dataCopy, Time: t, err: entry.err, Context: entry.Context}
}
func (entry Entry) HasCaller() (has bool) {
	return entry.Logger != nil &&
		entry.Logger.ReportCaller &&
		entry.Caller != nil
}
func (entry *Entry) Log(level Loglevel, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprint(args...))
	}
}
func (entry *Entry) Print(msg string, args ...Field) {
	entry.Info(msg, args...)
}
func (entry *Entry) Debug(msg string, args ...Field) {
	entry.WithFields(args...).Log(DebugLevel, msg)
}
func (entry *Entry) Info(msg string, args ...Field) {
	entry.WithFields(args...).Log(InfoLevel, msg)
}
func (entry *Entry) Warn(msg string, args ...Field) {
	entry.WithFields(args...).Log(WarnLevel, msg)
}
func (entry *Entry) Error(msg string, args ...Field) {
	entry.WithFields(args...).Log(ErrorLevel, msg)
}
func (entry *Entry) Fatal(msg string, args ...Field) {
	entry.WithFields(args...).Log(FatalLevel, msg)
	entry.Logger.Exit(1)
}
func (entry *Entry) Panic(msg string, args ...Field) {
	entry.WithFields(args...).Log(PanicLevel, msg)
}

func (entry *Entry) Logf(level Loglevel, format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.Log(level, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Logf(DebugLevel, format, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Logf(InfoLevel, format, args...)
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Logf(WarnLevel, format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Logf(ErrorLevel, format, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Logf(FatalLevel, format, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.Logf(PanicLevel, format, args...)
}

func (entry *Entry) Logln(level Loglevel, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.Log(level, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Debugln(args ...interface{}) {
	entry.Logln(DebugLevel, args...)
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.Logln(InfoLevel, args...)
}

func (entry *Entry) Println(args ...interface{}) {
	entry.Infoln(args...)
}

func (entry *Entry) Warnln(args ...interface{}) {
	entry.Logln(WarnLevel, args...)
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warnln(args...)
}

func (entry *Entry) Errorln(args ...interface{}) {
	entry.Logln(ErrorLevel, args...)
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.Logln(FatalLevel, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	entry.Logln(PanicLevel, args...)
}

func (entry *Entry) log(level Loglevel, msg string) {
	var buffer *bytes.Buffer

	newEntry := entry.Dup()

	if newEntry.Time.IsZero() {
		newEntry.Time = time.Now()
	}

	newEntry.Level = level
	newEntry.Message = msg

	if newEntry.Logger.ReportCaller {
		newEntry.Caller = getCaller(entry.Skip)
	}

	buffer = bufferPool.Get()
	defer func() {
		newEntry.Buffer = nil
		buffer.Reset()
		bufferPool.Put(buffer)
	}()
	buffer.Reset()
	newEntry.Buffer = buffer

	newEntry.write()

	newEntry.Buffer = nil

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(newEntry)
	}
}

func (entry *Entry) write() {
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	if _, err := entry.Logger.Out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}
func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
