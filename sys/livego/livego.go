package livego

func newSys(options Options) (sys *Engine, err error) {
	sys = &Engine{
		options: options,
	}
	return
}

type Engine struct {
	options Options
	rooms   *Rooms
}

func (this *Engine) CheckAppName(appname string) bool {
	for _, app := range this.options.Server {
		if app.Appname == appname {
			return app.Live
		}
	}
	return false
}
func (this *Engine) GetStaticPushUrlList(appname string) ([]string, bool) {
	for _, app := range this.options.Server {
		if (app.Appname == appname) && app.Live {
			if len(app.StaticPush) > 0 {
				return app.StaticPush, true
			} else {
				return nil, false
			}
		}
	}
	return nil, false
}
