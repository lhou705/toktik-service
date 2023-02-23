package client

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/hashicorp/consul/api"
	consul "github.com/kitex-contrib/registry-consul"
	"net/http"
	"time"
	"toktik/service/web/config"
	"toktik/service/web/kitex_gen/message/message"
	"toktik/service/web/kitex_gen/user/user"
	"toktik/service/web/kitex_gen/video/video"
)

var UserClient user.Client
var MessageClient message.Client
var VideoClient video.Client

func InitClient(clientConf config.Client, consulConf config.Consul) {
	r, err := consul.NewConsulResolverWithConfig(&api.Config{
		Address:    consulConf.Addr,
		Scheme:     "http",
		HttpClient: &http.Client{Timeout: 3 * time.Second},
		Token:      consulConf.Token,
	})
	//r, err := consul.NewConsulResolver(consulConf.Addr)
	if err != nil {
		hlog.Errorf("初始化注册中心失败，原因：%v", err)
	}
	initUserClient(clientConf.User.Name, r)
	initMessageClient(clientConf.Message.Name, r)
	initVideoClient(clientConf.Video.Name, r)

}

func initUserClient(name string, r discovery.Resolver) {
	userClient, err := user.NewClient(name, client.WithResolver(r))
	if err != nil {
		hlog.Errorf("初始化用户服务失败，原因：%v", err)
	}
	UserClient = userClient
}

func initMessageClient(name string, r discovery.Resolver) {
	messageClient, err := message.NewClient(name, client.WithResolver(r))
	if err != nil {
		hlog.Errorf("初始化信息服务失败，原因：%v", err)
	}
	MessageClient = messageClient
}

func initVideoClient(name string, r discovery.Resolver) {
	videoClient, err := video.NewClient(name, client.WithResolver(r))
	if err != nil {
		hlog.Errorf("初始化视频服务失败，原因：%v", err)
	}
	VideoClient = videoClient
}
