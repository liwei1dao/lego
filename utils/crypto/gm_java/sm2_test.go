package gm_java

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ZZMarquis/gm/sm2"
)

func Test_RSA_Decrypt(t *testing.T) {
	data, _ := base64.StdEncoding.DecodeString("FXlrn8jX61JDcBtOOh59yy/sM2r1hBT5XODayZKRDVE=")
	privKey, err := sm2.RawBytesToPrivateKey(data)
	fmt.Printf("privKey:%v err:%v\n", privKey, err)
	content, err := ioutil.ReadFile("encode.txt")
	fmt.Printf("content:%s err:%v\n", string(content), err)
	mdata, _ := base64.StdEncoding.DecodeString(string(content))
	orgdata, err := sm2.Decrypt(privKey, mdata, sm2.C1C3C2)
	fmt.Printf("orgdata:%s err:%v\n", string(orgdata), err)
}
