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
	"net"
	"net/http"
	"os"
	"time"
	"toktik/service/user/kitex_gen/user/user"
)

var Db *gorm.DB

func main() {
	configFilePath := flag.String("config", "../config/toktik_user.config.json", "配置文件路径")
	flag.Parse()
	fmt.Println("使用配置文件：" + *configFilePath)
	_, err := os.Stat(*configFilePath)
	if err != nil {
		klog.Fatalf("获取配置文件%s失败。错误原因：%v", *configFilePath, err)
	}
	// 初始化注册中心
	config := GetConfigFromFile(*configFilePath)
	r, err := consul.NewConsulRegisterWithConfig(&api.Config{
		Address:    config.Server.RegisterAddr,
		Scheme:     "http",
		HttpClient: &http.Client{Timeout: 3 * time.Second},
		Token:      config.Server.Token,
	})
	if err != nil {
		klog.Fatalf("初始化注册中心失败。错误原因：%v", err)
	}
	// 连接数据库
	initDatabase(config.Mysql.GetDsnStr())
	addr, err := net.ResolveTCPAddr("tcp", config.Server.Addr)
	if err != nil {
		klog.Fatalf("创建链接监听失败。错误原因：%v", err)
	}
	svr := user.NewServer(new(UserImpl),
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

func initDatabase(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		klog.Fatalf("连接数据库失败。错误原因：%v", err)
	}
	err = db.AutoMigrate(&User{}, &Follow{})
	if err != nil {
		klog.Fatalf("迁移数据表失败。错误原因：%v", err)
	}
	gplus.Init(db)
	Db = db
}
