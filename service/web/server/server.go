package server

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/utils"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
	"time"
	"toktik/service/web/router"
)

func initServer() {
	serverConf := conf.Server
	config := &consulApi.Config{
		Address: conf.Consul.Addr,
		Scheme:  "http",
	}
	consulClient, err := consulApi.NewClient(config)
	if err != nil {
		hlog.Fatalf("初始化网关微服务错误，原因:%v", err)
	}
	r := consul.NewConsulRegister(consulClient)
	h = server.Default(
		server.WithHostPorts(serverConf.Addr),
		server.WithReadTimeout(serverConf.ReadTimeOut*time.Second),
		server.WithWriteTimeout(serverConf.WriteTimeOut*time.Second),
		server.WithMaxRequestBodySize(serverConf.RequestBodySize),
		server.WithRegistry(r, &registry.Info{
			ServiceName: serverConf.Name,
			Addr:        utils.NewNetAddr("tcp", serverConf.RegisterAddr),
			Weight:      10,
			Tags:        nil,
		}),
	)
	router.RegisterRouter(h)
}
