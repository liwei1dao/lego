package paypal_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/sdks/paypal"
)

func Test_sys(t *testing.T) {
	if sys, err := paypal.NewSys(
		paypal.SetClientID("AepPyo1YtNF7SjKqdueN3oboCzs57qpUYU4lFdlbAnqo3emDicTKdbOZbvkacEuKH9L0r6ZNSws8-pJW"),
		paypal.SetSecretID("EDD2ofRG43vXgQhyShhdt3mC8zV7NPq_PJPQCr04Db2mr35uiB1FcL5LMzdp87-dHfPNlnHNciGujjI8"),
		paypal.SetIsSandBox(true),
		paypal.SetReturnURL("http://www.seawhale.store:8000/payback"),
	); err != nil {
		fmt.Printf("livego init err:%v \n", err)
		return
	} else {
		order, err := sys.CreateOrder("orderid", "liwei1dao", 1)
		fmt.Printf("order:%v err:%v \n", order, err)
	}

}
