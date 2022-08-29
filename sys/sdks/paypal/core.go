package paypal

import "github.com/plutov/paypal/v4"

/*
系统描述:paypal 支付服务
*/

type (
	ISys interface {
		///创建付款订单
		CreateOrder(orderId, uid string, amount float64) (order *paypal.Order, err error)
		//回调订单
		GetOrder(orderId string) (order *paypal.CaptureOrderResponse, err error)
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

func CreateOrder(orderId, uid string, amount float64) (order *paypal.Order, err error) {
	return defsys.CreateOrder(orderId, uid, amount)
}

func GetOrder(orderId string) (order *paypal.CaptureOrderResponse, err error) {
	return defsys.GetOrder(orderId)
}
