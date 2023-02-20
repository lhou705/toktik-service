package config

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	"os"
	"time"
)

type Config struct {
	Server Server
	Consul Consul
	JWT    JWT
	Client Client
}

type Server struct {
	Addr            string        `json:"addr,omitempty"`
	ReadTimeOut     time.Duration `json:"readTimeOut,omitempty"`
	WriteTimeOut    time.Duration `json:"writeTimeOut,omitempty"`
	RequestBodySize int           `json:"requestBodySize"`
	Name            string        `json:"name"`
	RegisterAddr    string        `json:"registerAddr"`
	Weight          int           `json:"weight"`
}

type JWT struct {
	Issuer              string        `json:"issuer,omitempty"`
	TokenExpireDuration time.Duration `json:"tokenExpireDuration,omitempty"`
	Secrete             string        `json:"secrete,omitempty"`
	IdentityKey         string        `json:"identityKey"`
}

type Consul struct {
	Addr string
}

type Client struct {
	User    BaseClient
	Message BaseClient
	Video   BaseClient
}

type BaseClient struct {
	Name string
}

func GetConfigFromFile(filename string) *Config {
	var config *Config
	bytes, err := os.ReadFile(filename)
	if err != nil {
		hlog.Fatalf("读取配置文件错误，错误原因：%v", err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		hlog.Fatalf("读取配置文件错误，错误原因：%v", err)
	}
	return config
}
