package main

import (
	"io/ioutil"

	"github.com/qf-tech/qftech-base/pkg/log"
	"github.com/qf-tech/qftech-base/pkg/qfsign"
)

func main() {
	log.Init(nil)
	log.Sugare.Info("hello world")

	// 1.new the hadler to sign and verify
	cfg := qfsign.ConfigParams{
		Algorithm:          qfsign.RSASign,
		PublicKeyCertPath:  "./pem/rsa_public_key.pem",
		PrivateKeyCertPath: "./pem/rsa_private_key.pem",
	}
	handler, err := qfsign.NewHandler(cfg)
	if err != nil {
		log.Sugare.Errorf("new sign handler err: %v", err)
		return
	}

	// 2.generate sign
	data, _ := ioutil.ReadFile("./test.sh")
	sign, err := handler.Sign(data)
	if err != nil {
		log.Sugare.Errorf("gen sign err: %v", err)
		return
	}
	log.Sugare.Infof("gen sign: %s", sign)

	// 3.verify sign
	signTmp := "YBPQE1reP29Ebp4n+8EVJavtW9oXf7qqRhIr2/QYVMqoo8D3tDCbW04j69bXtz6/rlxjl25YQri0gKk+JpuW1yE2NMY+h3rIV3JgBxm5YORkZCOd99eYOJ6sVRAMBUlaWhMxc91O3kLq0K9CSAL35VKdWUHDaZDiTKRxnq+HryM="
	err = handler.Verify(data, signTmp)
	if err != nil {
		log.Sugare.Errorf("verify sign err: %v", err)
		return
	}
	log.Sugare.Info("verify sign succeed")
}
