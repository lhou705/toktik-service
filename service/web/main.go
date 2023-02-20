package main

import (
	"flag"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"os"
	"toktik/service/web/server"
)

func main() {
	config := flag.String("config", "./config/gateway.json", "网关的配置文件")
	flag.Parse()
	_, err := os.Stat(*config)
	if err != nil {
		hlog.Fatalf("读取配置文件错误，原因%v", err)
	}
	server.LoadConfig(*config)
	server.Init()
	server.Start()
}
