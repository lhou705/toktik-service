package main

import (
	"flag"
	"fmt"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/hashicorp/consul/api"
	consul "github.com/kitex-contrib/registry-consul"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"toktik/service/video/kitex_gen/video/video"
)

var Db *gorm.DB
var config *Config

func main() {
	configFilePath := flag.String("config", "../config/toktik_video.config.json", "配置文件路径")
	flag.Parse()
	fmt.Println("使用配置文件：" + *configFilePath)
	_, err := os.Stat(*configFilePath)
	if err != nil {
		klog.Fatalf("获取配置文件%s失败。错误原因：%v", *configFilePath, err)
	}
	// 初始化注册中心
	conf := GetConfigFromFile(*configFilePath)
	r, err := consul.NewConsulRegisterWithConfig(&api.Config{
		Address:    conf.Server.RegisterAddr,
		Scheme:     "http",
		HttpClient: &http.Client{Timeout: 3 * time.Second},
		Token:      conf.Server.Token,
	})
	if err != nil {
		klog.Fatalf("初始化注册中心失败。错误原因：%v", err)
	}
	// 连接数据库
	addr, err := net.ResolveTCPAddr("tcp", conf.Server.Addr)
	if err != nil {
		klog.Fatalf("创建链接监听失败。错误原因：%v", err)
	}
	initDatabase(conf.Mysql.GetDsnStr())
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
