package main

import (
	"flag"
	"fmt"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/hashicorp/consul/api"
	consul "github.com/kitex-contrib/registry-consul"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"time"
	"toktik/service/video/kitex_gen/cos/cos"
	"toktik/service/video/kitex_gen/video/video"
)

var Db *gorm.DB
var cosClient cos.Client
var config *Config

func main() {
	configFilePath := flag.String("config", "../config/video.config.json", "配置文件路径")
	flag.Parse()
	fmt.Println("使用配置文件：" + *configFilePath)
	_, err := os.Stat(*configFilePath)
	if err != nil {
		klog.Fatalf("获取配置文件%s失败。错误原因：%v", *configFilePath, err)
	}
	// 初始化注册中心
	conf := GetConfigFromFile(*configFilePath)
	check := &api.AgentServiceCheck{
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "1m",
	}
	r, err := consul.NewConsulRegister(conf.Server.RegisterAddr, consul.WithCheck(check))
	//r, err := consul.NewConsulRegisterWithConfig(&api.Config{
	//	Address: conf.Server.RegisterAddr,
	//	Scheme:  "http",
	//})
	if err != nil {
		klog.Fatalf("初始化注册中心失败。错误原因：%v", err)
	}
	// 连接数据库
	addr, err := net.ResolveTCPAddr("tcp", conf.Server.Addr)
	if err != nil {
		klog.Fatalf("创建链接监听失败。错误原因：%v", err)
	}
	initDatabase(conf.Mysql.GetDsnStr())
	// 连接cos服务
	initCosClient(conf.Cos.Name, conf.Server.RegisterAddr)
	svr := video.NewServer(new(VideoImpl),
		server.WithServiceAddr(addr),
		server.WithReusePort(conf.Server.ReusePort),
		server.WithReadWriteTimeout(conf.Server.ReadWriteTimeOut*time.Second),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.Server.Name}))
	config = conf
	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

func initDatabase(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		klog.Fatalf("连接数据库失败。错误原因：%v", err)
	}
	err = db.AutoMigrate(&Video{}, &Favorite{}, &Comment{})
	if err != nil {
		klog.Fatalf("迁移数据表失败。错误原因：%v", err)
	}
	gplus.Init(db)
	Db = db
}

func initCosClient(name string, consulAddr string) {
	r, err := consul.NewConsulResolver(consulAddr)
	if err != nil {
		klog.Errorf("初始化注册中心失败，原因：%v", err)
	}
	newClient, err := cos.NewClient(name, client.WithResolver(r))
	if err != nil {
		klog.Errorf("初始化信息服务失败，原因：%v", err)
	}
	cosClient = newClient
}
