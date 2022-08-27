package paypal

import (
	"context"
	"errors"
	"fmt"

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

	this.client.SetLog(this.options.Log)
	this.accessToken, err = this.client.GetAccessToken(context.Background())
	return
}

/*
创建收款订单
(*order).Links[1].Href就是支付的链接
*/
func (this *PayPal) CreateOrder(id string, amount float64) (order *paypal.Order, err error) {
	purchaseUnits := make([]paypal.PurchaseUnitRequest, 1)
	purchaseUnits[0] = paypal.PurchaseUnitRequest{
		Amount: &paypal.PurchaseUnitAmount{
			Currency: this.options.Currency,     //收款类型
			Value:    fmt.Sprintf("%v", amount), //收款数量
		},
	}
	payer := &paypal.CreateOrderPayer{
		Name: &paypal.CreateOrderPayerName{
			GivenName: id,
			Surname:   this.options.AppName,
		},
	}
	appContext := &paypal.ApplicationContext{
		ReturnURL: this.options.ReturnURL, //回调链接
	}
	order, err = this.client.CreateOrder(context.Background(), "CAPTURE", purchaseUnits, payer, appContext)
	if err != nil {
		this.options.Log.Errorf("create order errors:%v", err)
		return
	}
	return
}

//回调(可以利用上面的回调链接实现)
func (this *PayPal) PaypalCallback(orderId string) error {
	ctor := paypal.CaptureOrderRequest{}
	order, err := this.client.CaptureOrder(context.Background(), orderId, ctor)
	if err != nil {
		return err
	}
	//查看回调完成后订单状态是否支付完成。
	if (*order).Status != "COMPLETED" {
		return errors.New("pay fail")
	}
	return nil
}
