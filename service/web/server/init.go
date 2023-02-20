package server

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"toktik/service/web/client"
	"toktik/service/web/config"
	"toktik/service/web/utils"
)

var h *server.Hertz
var conf *config.Config

func LoadConfig(filename string) {
	conf = config.GetConfigFromFile(filename)
}

func Init() {
	initServer()
	utils.InitConfig(conf.JWT)
	client.InitClient(conf.Client, conf.Consul)
}

func Start() {
	h.Spin()
}
