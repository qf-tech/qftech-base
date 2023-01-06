package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"

	"github.com/qf-tech/qftech-base/pkg/qfsign"
)

const Version = "v1.0.0"

type Config struct {
	Mode           string
	PublicKeyCert  string
	PrivateKeyCert string
	FilePath       string
	Data           string
	Sign           string
	verify         bool
	PrintVersion   bool
}

var config Config

func init() {
	flag.StringVar(&config.Mode, "mode", "", "sign algorithm, such as: rsa")
	flag.StringVar(&config.PublicKeyCert, "public-key-cert", "", "public key cert path")
	flag.StringVar(&config.PrivateKeyCert, "private-key-cert", "", "private key cert path")
	flag.StringVar(&config.FilePath, "file-path", "", "the path of file that need to be signed or verify")
	flag.StringVar(&config.Data, "data", "", "the data need to be signed or verify, if you use file-path, the data will disable")
	flag.StringVar(&config.Sign, "sign", "", "the sign of file data, when use -verify should input -sign")
	flag.BoolVar(&config.verify, "verify", false, "verify the sign, otherwise generate the sign of file data")
	flag.BoolVar(&config.PrintVersion, "version", false, "print version and exit")
}

func assertConfigNull() bool {
	return reflect.DeepEqual(config, Config{})
}

func getMacAddrs() (macAddrs string) {
	var macs []string
	var macs23 []string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v\n", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 17 {
			macs = append(macs, macAddr)
		} else {
			macs23 = append(macs23, macAddr)
		}
	}
	sort.Strings(macs)
	sort.Strings(macs23)
	if len(macs) >= 1 {
		macAddrs = macs[0]
		return
	} else if len(macs) == 0 && len(macs23) >= 1 {
		macAddrs = macs23[0]
		return
	}
	return
}

func getMacsAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v\n", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		macAddrs = append(macAddrs, macAddr)
	}

	return
}

func main() {
	path, _ := os.Executable()
	_, svrName := filepath.Split(path)
	flag.Parse()
	if assertConfigNull() {
		fmt.Printf("%s -h or %s -h to get it usage\n", svrName, svrName)
		os.Exit(0)
	}

	if config.PrintVersion {
		fmt.Printf("%s %s (Go Version: %s)\n", svrName, Version, runtime.Version())
		os.Exit(0)
	}

	// 1.new the hadler to sign and verify
	cfg := qfsign.ConfigParams{
		Algorithm:          qfsign.AlgorithmType(config.Mode),
		PublicKeyCertPath:  config.PublicKeyCert,
		PrivateKeyCertPath: config.PrivateKeyCert,
	}
	handler, err := qfsign.NewHandler(cfg)
	if err != nil {
		fmt.Printf("new sign handler err: %v", err)
		return
	}

	var data []byte
	if config.FilePath == "" {
		data = []byte(config.Data)
		if config.Data == "" {
			data = []byte(getMacAddrs())
		}
	} else {
		if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
			fmt.Printf("file path not exist, err: %v\n", err)
			return
		}
		data, err = ioutil.ReadFile(config.FilePath)
		if err != nil {
			fmt.Printf("read file data err: %v\n", err)
			return
		}
	}

	// 2.generate sign or verify
	if !config.verify {
		sign, err := handler.Sign(data)
		if err != nil {
			fmt.Printf("gen sign err: %v\n", err)
			return
		}

		if config.FilePath == "" && config.Data == "" {
			fmt.Printf("mac data: %s\n", string(data))
		}
		fmt.Printf("generate sign: %s\n", sign)
	} else {
		if config.Sign == "" {
			pth := "/opt/openresty/nginx/.license/casb.license"
			if _, err := os.Stat(pth); os.IsNotExist(err) {
				fmt.Printf("license path not exist, err: %v\n", err)
				return
			}
			tmpData, err := ioutil.ReadFile(pth)
			if err != nil {
				fmt.Printf("read license file err: %v\n", err)
				return
			}
			config.Sign = string(tmpData)
		}

		if config.FilePath != "" || config.Data != "" {
			if err := handler.Verify(data, config.Sign); err != nil {
				fmt.Printf("verify sign err: %v\n", err)
				return
			}
			fmt.Println("verify sign succeed")
		} else {
			macs := getMacsAddrs()
			var err error
			for _, mac := range macs {
				if err = handler.Verify([]byte(mac), config.Sign); err != nil {
					fmt.Printf("verify sign err: %v\n", err)
					continue
				}
				fmt.Printf("verify sign succeed, mac: %s\n", mac)
				break
			}
		}
	}
}
