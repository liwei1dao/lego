package pay

type (
	ISys interface {
		WebPageAlipay(subject string, orderid string, amount string) (payurl string, err error)
		WapAlipay(subject string, orderid string, amount string) (payurl string, err error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

//网站扫码支付
func WebPageAlipay(subject string, orderid string, amount string) (payurl string, err error) {
	return defsys.WebPageAlipay(subject, orderid, amount)
}

//手机客户端支付
func WapAlipay(subject string, orderid string, amount string) (payurl string, err error) {
	return defsys.WapAlipay(subject, orderid, amount)
}
