package ethpay

import "github.com/liwei1dao/lego/core"

type (
	IETHPay interface {
		LookBalance(addr string) uint64                       //查看余额
		GetUserPayAddr(uhash string) (addr string, err error) //获取用户支付地址
		RecycleUserMoney(uaddr string) (trans string,err error) //回收资金
	}
)

var (
	pay IETHPay
)

func OnInit(s core.IService, opt ...Option) (err error) {
	pay, err = newEthPay(opt...)
	return
}

func LookBalance(addr string) uint64 {
	return pay.LookBalance(addr)
}

func GetUserPayAddr(uhash string) (addr string, err error) {
	return pay.GetUserPayAddr(uhash)
}

func RecycleUserMoney(uaddr string) (trans string,err error){
	return pay.RecycleUserMoney(uaddr)
}


func mustHexDecode(raw string) []byte {
	if raw == "0x" {
		return []byte{}
	}
	if len(raw) > 2 && raw[:2] == "0x" {
		raw = raw[2:]
	}
	data, err := hex.DecodeString(raw)
	if err != nil {
		panic(err)
	}
	return data
}

func keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		_, _ = d.Write(b)
	}
	return d.Sum(nil)
}

