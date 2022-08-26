package paypal

import "github.com/plutov/paypal/v4"

/*
系统描述:paypal 支付服务
*/

type (
	ISys interface {
		///创建付款订单
		CreateOrder(amount string) (order *paypal.Order, err error)
		//回调订单 查看是否完成
		PaypalCallback(orderId string) error
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

func CreateOrder(amount string) (order *paypal.Order, err error) {
	return defsys.CreateOrder(amount)
}

func PaypalCallback(orderId string) error {
	return defsys.PaypalCallback(orderId)
}
