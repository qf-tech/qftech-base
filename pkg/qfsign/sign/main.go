package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/qf-tech/qftech-base/pkg/qfsign"
)

const Version = "v1.0.0"

type Config struct {
	Mode           string
	PublicKeyCert  string
	PrivateKeyCert string
	FilePath       string
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
	flag.StringVar(&config.Sign, "sign", "", "the sign of file data, when use -verify should input -sign")
	flag.BoolVar(&config.verify, "verify", false, "verify the sign, otherwise generate the sign of file data")
	flag.BoolVar(&config.PrintVersion, "version", false, "print version and exit")
}

func assertConfigNull() bool {
	return reflect.DeepEqual(config, Config{})
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

	if config.FilePath == "" {
		fmt.Println("err: file path is empty, can't generate sign")
		return
	}

	if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
		fmt.Printf("file path not exist, err: %v\n", err)
		return
	}
	data, err := ioutil.ReadFile(config.FilePath)
	if err != nil {
		fmt.Printf("read file data err: %v\n", err)
		return
	}

	// 2.generate sign or verify
	if !config.verify {
		sign, err := handler.Sign(data)
		if err != nil {
			fmt.Printf("gen sign err: %v\n", err)
			return
		}
		fmt.Printf("generate sign: %s\n", sign)
	} else {
		if config.Sign == "" {
			fmt.Println("err: sign is empty, cann't to verify")
			return
		}

		if err := handler.Verify(data, config.Sign); err != nil {
			fmt.Printf("verify sign err: %v\n", err)
			return
		}
		fmt.Println("verify sign succeed")
	}
}
