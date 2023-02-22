package main

import (
	"flag"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/hashicorp/consul/api"
	consul "github.com/kitex-contrib/registry-consul"
	cosSdk "github.com/tencentyun/cos-go-sdk-v5"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
	"toktik/service/cos/kitex_gen/cos/cos"
)

var client *cosSdk.Client

func main() {
	// 读取配置文件
	configFilePath := flag.String("config", "../config/toktik_cos.config.json", "配置文件路径")
	flag.Parse()
	fmt.Println("使用配置文件：" + *configFilePath)
	_, err := os.Stat(*configFilePath)
	if err != nil {
		klog.Fatalf("获取配置文件%s失败。错误原因：%v", *configFilePath, err)
	}
	// 加载cos服务
	config := GetConfigFromFile(*configFilePath)
	NewClient(config.Cos.Addr, config.Cos.SecretID, config.Cos.SecretKey)
	// 初始化注册中心
	//r, err := consul.NewConsulRegister(config.Server.RegisterAddr,consul.WithCheck(&api.AgentServiceCheck{}))
	check := &api.AgentServiceCheck{
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "1m",
	}
	r, err := consul.NewConsulRegister(config.Server.RegisterAddr, consul.WithCheck(check))
	if err != nil {
		klog.Fatalf("初始化注册中心失败。错误原因：%v", err)
	}
	// 启动微服务
	addr, err := net.ResolveTCPAddr("tcp", config.Server.Addr)
	if err != nil {
		klog.Fatalf("创建链接监听失败。错误原因：%v", err)
	}
	svr := cos.NewServer(new(CosImpl),
		server.WithServiceAddr(addr),
		server.WithReusePort(config.Server.ReusePort),
		server.WithReadWriteTimeout(config.Server.ReadWriteTimeOut*time.Second),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.Server.Name}))
	err = svr.Run()
	if err != nil {
		klog.Fatalf("微服务启动失败。错误原因：%v", err)
	}
}

func NewClient(urlStr, SecretID, SecretKey string) {
	u, err := url.Parse(urlStr)
	if err != nil {
		klog.Fatalf("对象存储服务连接匹配失败。错误原因：%v", err)
	}
	b := &cosSdk.BaseURL{
		BucketURL: u,
	}
	client = cosSdk.NewClient(b, &http.Client{
		Transport: &cosSdk.AuthorizationTransport{
			SecretID:  SecretID,
			SecretKey: SecretKey,
		}})
}
