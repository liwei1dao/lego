package log

func NewTurnlog(isturnon bool, log ILogger) ILogger {
	return &Turnlog{
		isturnon: isturnon,
		log:      log,
	}
}

type Turnlog struct {
	isturnon bool
	log      ILogger
}

func (this *Turnlog) SetName(name string) {
	if this.log != nil {
		this.log.SetName(name)
	}
}

func (this *Turnlog) Enabled(lvl Loglevel) bool {
	if this.isturnon && this.log != nil {
		return this.log.Enabled(lvl)
	} else {
		return false
	}
}
func (this *Turnlog) Debug(msg string, args ...Field) {
	if this.isturnon && this.log != nil {
		this.log.Debug(msg, args...)
	}
}
func (this *Turnlog) Info(msg string, args ...Field) {
	if this.isturnon && this.log != nil {
		this.log.Info(msg, args...)
	}
}
func (this *Turnlog) Print(msg string, args ...Field) {
	if this.isturnon && this.log != nil {
		this.log.Print(msg, args...)
	}
}
func (this *Turnlog) Warn(msg string, args ...Field) {
	if this.isturnon && this.log != nil {
		this.log.Warn(msg, args...)
	}
}
func (this *Turnlog) Error(msg string, args ...Field) {
	if this.log != nil {
		this.log.Error(msg, args...)
	}
}
func (this *Turnlog) Panic(msg string, args ...Field) {
	if this.log != nil {
		this.log.Panic(msg, args...)
	}
}
func (this *Turnlog) Fatal(msg string, args ...Field) {
	if this.log != nil {
		this.log.Fatal(msg, args...)
	}
}
func (this *Turnlog) Debugf(format string, args ...interface{}) {
	if this.isturnon && this.log != nil {
		this.log.Debugf(format, args...)
	}
}
func (this *Turnlog) Infof(format string, args ...interface{}) {
	if this.isturnon && this.log != nil {
		this.log.Infof(format, args...)
	}
}
func (this *Turnlog) Printf(format string, args ...interface{}) {
	if this.isturnon && this.log != nil {
		this.log.Printf(format, args...)
	}
}
func (this *Turnlog) Warnf(format string, args ...interface{}) {
	if this.isturnon && this.log != nil {
		this.log.Warnf(format, args...)
	}
}
func (this *Turnlog) Errorf(format string, args ...interface{}) {
	if this.log != nil {
		this.log.Errorf(format, args...)
	}
}
func (this *Turnlog) Fatalf(format string, args ...interface{}) {
	if this.log != nil {
		this.log.Fatalf(format, args...)
	}
}
func (this *Turnlog) Panicf(format string, args ...interface{}) {
	if this.log != nil {
		this.log.Panicf(format, args...)
	}
}
func (this *Turnlog) Debugln(args ...interface{}) {
	if this.isturnon && this.log != nil {
		this.log.Debugln(args...)
	}
}
func (this *Turnlog) Infoln(args ...interface{}) {
	if this.isturnon && this.log != nil {
		this.log.Infoln(args...)
	}
}
func (this *Turnlog) Println(args ...interface{}) {
	if this.isturnon && this.log != nil {
		this.log.Println(args...)
	}
}
func (this *Turnlog) Warnln(args ...interface{}) {
	if this.isturnon && this.log != nil {
		this.log.Warnln(args...)
	}
}
func (this *Turnlog) Errorln(args ...interface{}) {
	if this.log != nil {
		this.log.Errorln(args...)
	}
}
func (this *Turnlog) Fatalln(args ...interface{}) {
	if this.log != nil {
		this.log.Fatalln(args...)
	}
}
func (this *Turnlog) Panicln(args ...interface{}) {
	if this.log != nil {
		this.log.Panicln(args...)
	}
}
