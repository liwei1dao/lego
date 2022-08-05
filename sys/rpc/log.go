package rpc

///日志***********************************************************************
func (this *rpc) Debug() bool {
	return this.options.Debug
}
func (this *rpc) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Debugf("[SYS rpc] "+format, a)
	}
}
func (this *rpc) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Infof("[SYS rpc] "+format, a)
	}
}
func (this *rpc) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Warnf("[SYS rpc] "+format, a)
	}
}
func (this *rpc) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Errorf("[SYS rpc] "+format, a)
	}
}
func (this *rpc) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Panicf("[SYS rpc] "+format, a)
	}
}
func (this *rpc) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Fatalf("[SYS rpc] "+format, a)
	}
}
