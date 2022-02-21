package main

import (
	logs "github.com/qf-tech/qftech-base/pkg/log"
)

func main() {
	// use your config
	// config := &logs.LogConfig{
	// 	MaxCount: 30,
	// 	MaxSize:  10,
	// 	Compress: true,
	// 	FilePath: "./log/server.log",
	// 	Level:    logs.InfoLevel,
	// }

	// if config is nil, will use default config
	logs.Init(nil)
	logs.Sugare.Info("hello world")
}
