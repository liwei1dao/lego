package pay

import (
	"net/url"

	"github.com/smartwalle/alipay/v3"
)

func newSys(options *Options) (sys *Pay, err error) {
	sys = &Pay{options: options}
	if sys.client, err = alipay.New(options.AppId, options.PrivateKey, true); err != nil {
		return
	}
	sys.client.LoadAppPublicCert("应用公钥证书")
	sys.client.LoadAliPayPublicCert("支付宝公钥证书")
	sys.client.LoadAliPayRootCert("支付宝根证书")
	return
}

type Pay struct {
	options *Options
	client  *alipay.Client
}

//网站扫码支付
func (this *Pay) WebPageAlipay(subject string, orderid string, amount string) (payurl string, err error) {
	var (
		pay alipay.TradePagePay
		url *url.URL
	)

	pay = alipay.TradePagePay{}
	// 支付宝回调地址（需要在支付宝后台配置）
	// 支付成功后，支付宝会发送一个POST消息到该地址 "http://www.pibigstar/alipay"
	pay.NotifyURL = this.options.NotifyURL
	// 支付成功之后，浏览器将会重定向到该 URL"http://localhost:8088/return"
	pay.ReturnURL = this.options.ReturnURL
	//支付标题"支付宝支付测试"
	pay.Subject = subject
	//订单号，一个订单号只能支付一次
	pay.OutTradeNo = orderid
	//销售产品码，与支付宝签约的产品码名称,目前仅支持FAST_INSTANT_TRADE_PAY
	pay.ProductCode = "FAST_INSTANT_TRADE_PAY"
	//金额 "0.01"
	pay.TotalAmount = amount
	url, err = this.client.TradePagePay(pay)
	if err != nil {
		return
	}
	payurl = url.String()
	//这个 payURL 即是用于支付的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	return
}

//手机客户端支付
func (this *Pay) WapAlipay(subject string, orderid string, amount string) (payurl string, err error) {
	var (
		pay alipay.TradeWapPay
		url *url.URL
	)
	pay = alipay.TradeWapPay{}
	// 支付成功之后，支付宝将会重定向到该 URL "http://localhost:8088/return"
	pay.ReturnURL = this.options.ReturnURL
	//支付标题"支付宝支付测试"
	pay.Subject = subject
	//订单号，一个订单号只能支付一次
	pay.OutTradeNo = orderid
	//商品code
	pay.ProductCode = "QUICK_MSECURITY_PAY"
	//金额
	pay.TotalAmount = amount

	url, err = this.client.TradeWapPay(pay)
	if err != nil {
		return
	}
	payurl = url.String()
	//这个 payURL 即是用于支付的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	// fmt.Println(payURL)
	// //打开默认浏览器
	// payURL = strings.Replace(payURL, "&", "^&", -1)
	// exec.Command("cmd", "/c", "start", payURL).Start()
	return
}
