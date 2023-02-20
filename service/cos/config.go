package main

import (
	"encoding/json"
	"github.com/cloudwego/kitex/pkg/klog"
	"os"
	"time"
)

type Config struct {
	Cos    Cos
	Server Server
}

type Server struct {
	Addr             string
	ReadWriteTimeOut time.Duration
	ReusePort        bool
	RegisterAddr     string
	Name             string
}

type Cos struct {
	Addr      string
	SecretID  string
	SecretKey string
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
