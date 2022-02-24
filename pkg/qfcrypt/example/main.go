package main

import (
	logs "github.com/qf-tech/qftech-base/pkg/log"
	"github.com/qf-tech/qftech-base/pkg/qfcrypt"
)

func initCrypt() {
	curVer := "QFTECH"
	key := "lmN4dtPyeC5r29DYBLl0P0OoA4Afy/2UnCg0zd+hHhg="

	oldKeys := make(map[string][]byte)
	zeroKeytr := "q_1dY=Khec2nMNxV"
	oldKeys[qfcrypt.ZeroVersionFlag] = []byte(zeroKeytr)

	_ = qfcrypt.Init(curVer, []byte(key), oldKeys)

}

func main() {
	logs.Init(nil)

	curVer := "QFTECH"
	initCrypt()

	plainText := "hello world"
	encData, err := qfcrypt.ConfigAes.Encrypt([]byte(plainText), curVer)
	if err != nil {
		logs.Sugare.Errorf("encrypt err: %v", err)
		return
	}

	decData, err := qfcrypt.ConfigAes.Decrypt(encData, curVer)
	if err != nil {
		logs.Sugare.Errorf("decrypt err: %v", err)
		return
	}

	logs.Sugare.Infof("encrypt data: %s", encData)
	logs.Sugare.Infof("decrypt data: %s", decData)
}
