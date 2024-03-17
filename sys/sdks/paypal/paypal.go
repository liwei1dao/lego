package paypal

import (
	"context"

	"github.com/plutov/paypal/v4"
)

func newSys(options *Options) (sys *PayPal, err error) {
	sys = &PayPal{options: options}
	err = sys.init()
	return
}

type PayPal struct {
	options     *Options
	client      *paypal.Client
	accessToken *paypal.TokenResponse
}

func (this *PayPal) init() (err error) {
	if this.options.IsSandBox {
		if this.client, err = paypal.NewClient(this.options.ClientID, this.options.SecretID, paypal.APIBaseSandBox); err != nil {
			this.options.Log.Errorln(err)
			return
		}
	} else {
		if this.client, err = paypal.NewClient(this.options.ClientID, this.options.SecretID, paypal.APIBaseLive); err != nil {
			this.options.Log.Errorln(err)
			return
		}
	}

	this.client.SetLog(this)
	this.accessToken, err = this.client.GetAccessToken(context.Background())
	return
}

func (this *PayPal) Write(p []byte) (n int, err error) {
	n = len(p)
	this.options.Log.Debugf(string(p))
	return
}

/*
创建收款订单
(*order).Links[1].Href就是支付的链接
*/
func (this *PayPal) CreateOrder(orderId, uid string, amount string) (order *paypal.Order, err error) {
	purchaseUnits := make([]paypal.PurchaseUnitRequest, 1)
	purchaseUnits[0] = paypal.PurchaseUnitRequest{
		Amount: &paypal.PurchaseUnitAmount{
			Currency: this.options.Currency, //收款类型
			Value:    amount,                //收款数量
		},
		InvoiceID:   orderId,
		ReferenceID: uid,
		Description: this.options.AppName,
	}
	payer := &paypal.CreateOrderPayer{
		Name: &paypal.CreateOrderPayerName{
			GivenName: uid,
			Surname:   uid,
		},
	}
	appContext := &paypal.ApplicationContext{
		ReturnURL: this.options.ReturnURL + orderId, //回调链接
	}
	order, err = this.client.CreateOrder(context.Background(), "CAPTURE", purchaseUnits, payer, appContext)
	if err != nil {
		this.options.Log.Errorf("create order errors:%v", err)
		return
	}
	return
}

//回调(可以利用上面的回调链接实现)
func (this *PayPal) GetOrder(orderId string) (order *paypal.CaptureOrderResponse, err error) {
	ctor := paypal.CaptureOrderRequest{}
	order, err = this.client.CaptureOrder(context.Background(), orderId, ctor)
	return
}

//回调(可以利用上面的回调链接实现) orderId 就是返回的token
func (this *PayPal) PaypalCallback(orderId string) (isucc bool, err error) {
	var (
		order *paypal.CaptureOrderResponse
	)
	_, err = this.client.GetAccessToken(context.TODO())
	if err != nil {
		return
	}
	//log.Info(accessToken.Token,orderId)
	ctor := paypal.CaptureOrderRequest{}
	order, err = this.client.CaptureOrder(context.TODO(), orderId, ctor)
	if err != nil {
		return
	}
	//查看回调完成后订单状态是否支付完成。
	if (*order).Status != "COMPLETED" {
		return
	}
	isucc = true
	this.options.Log.Debugf("支付成功 Status:%v Address:%v", order.Status, order.Address)
	return
}
