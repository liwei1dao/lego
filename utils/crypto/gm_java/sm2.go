package gm_java

/*
国密算法库 github.com/ZZMarquis/gm/sm2 封装
主要处理Java 版本国密算法
*/

import (
	"encoding/base64"
	"fmt"

	"github.com/ZZMarquis/gm/sm2"
)

//样例 base64Key:FXlrn8jX61JDcBtOOh59yy/sM2r1hBT5XODayZKRDVE=
func ReadPrivateKeyFormBase64(base64Key string) (privKey *sm2.PrivateKey, err error) {
	var (
		key []byte
	)
	if key, err = base64.StdEncoding.DecodeString(base64Key); err != nil {
		err = fmt.Errorf("ReadPrivateKeyFormBase64 Base64 DecodeString err:%v", err)
		return
	}
	privKey, err = sm2.RawBytesToPrivateKey(key)
	return
}

//样例 base64Key:BHa3F+W4YhmWoqa7glAURrU7vijUSNtg+9ZnQREuq/8+6MsGAc7
func ReadPublicKeyFormBase64(base64Key string) (pubKey *sm2.PublicKey, err error) {
	var (
		key []byte
	)
	if key, err = base64.StdEncoding.DecodeString(base64Key); err != nil {
		err = fmt.Errorf("ReadPublicKeyFormBase64 Base64 DecodeString err:%v", err)
		return
	}
	pubKey, err = sm2.RawBytesToPublicKey(key)
	return
}

///SM2 加密
func SM2_Encrypt(pubKey *sm2.PublicKey, in []byte) ([]byte, error) {
	return sm2.Encrypt(pubKey, in, sm2.C1C3C2)
}

///SM2 解密
func SM2_Decrypt(privKey *sm2.PrivateKey, in []byte) ([]byte, error) {
	return sm2.Decrypt(privKey, in, sm2.C1C3C2)
}
