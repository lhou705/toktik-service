package main

import (
	"encoding/json"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"os"
	"time"
)

type Mysql struct {
	Addr     string `json:"addr,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
}

func (m *Mysql) GetDsnStr() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Addr, m.Database)
	return dsn
}

type Server struct {
	Addr             string
	ReadWriteTimeOut time.Duration
	ReusePort        bool
	RegisterAddr     string
	Name             string
	Token            string
}

type Cos struct {
	CdnAddr string
}

type Config struct {
	Mysql  Mysql
	Server Server
	Cos    Cos
}

func GetConfigFromFile(filename string) *Config {
	var config *Config
	bytes, err := os.ReadFile(filename)
	if err != nil {
		klog.Fatalf("读取配置文件失败，原因%v", err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		klog.Fatalf("读取配置文件失败，原因%v", err)
	}
	return config
}
