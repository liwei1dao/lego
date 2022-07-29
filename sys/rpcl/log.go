package rpcl

///日志***********************************************************************
func (this *RPCL) Debug() bool {
	return this.options.Debug
}
func (this *RPCL) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Debugf("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Infof("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Warnf("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Errorf("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Panicf("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Fatalf("[SYS RPCL] "+format, a)
	}
}
