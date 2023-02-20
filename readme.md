# toktik
## 介绍
第五届字节跳动青训营汪汪大队作品

## 软件架构
本项目采用服务架构，使用[consul](https://www.consul.io/)作为注册中心，[kitex](https://www.cloudwego.io/zh/docs/kitex/)框架搭建。分为用户、通信、对象存储、网关、视频五个部分。其中网关部分使用[hertz](https://www.cloudwego.io/zh/docs/hertz/)搭建。
文件存储使用腾讯云对象存储。视频封面在上传视频后自动生成。为了加快访问速度，还使用了cdn技术加速访问。


## 安装教程
### 必要使用环境安装
1. 安装mysql数据库
2. 安装consul
### 安装本项目
1. 克隆本项目
```
git clone https://gitee.com/lhou/toktik.git
```
2. 切换到当前分支
```
git checkout  master
```
2.编译
```
bash cmd/build.sh
```
## 使用
### 设置配置文件
#### 网关的配置文件`web.config.json`
```json lines
{
  "server": {
    "addr": "your server addr",
    "readTimeOut": 10,
    "writeTimeOut": 10,
    "name": "your server name in consul",
    "weight": 10,
    "registerAddr": "your server addr registered in consul",
    "requestBodySize": 100000000000
  },
  "jwt": {
    "issuer": "your issuer",
    "tokenExpireDuration": 1440,
    "secrete": "your secret",
    "identityKey": "your identity key"
  },
  "consul": {
    "addr": "your consul server addr"
  },
  "client": {
    "user": {
      "name": "toktik.user"
    },
    "message": {
      "name": "toktik.message"
    },
    "video": {
      "name": "toktik.video"
    }
  }
}
```
#### 用户、信息模块配置文件`user.config.json`、`message.config.json`
```json lines
{
  "server": {
    "name": "your server name",
    "addr": "your server addr",
    "reusePort": true,
    "readWriteTimeOut": 10,
    "registerAddr": "your consul addr"
  },
  "mysql": {
    "addr": "your mysql server addr",
    "user": "your user",
    "password": "your password",
    "database": "your database"
  }
}
```
#### 视频模块配置文件`video.config.json`
```
{
  "server": "同用户、信息模块配置文件",
  "mysql": "同用户、信息模块配置文件",
  "cos": {
    "name": "对象存储的服务名称",
    "cdnAddr": "cdn加速域名"
  }
}
```
#### 腾讯云对象存储配置文件`cos.config.json`
```json lines
{
  "cos": {
    "addr": "腾讯云对象存储访问域名",
    "secretID": "访问Id",
    "secretKey": "访问密钥"
  },
  "server": "同用户、信息模块配置文件"
}
```
准备好上述文件后，将其放入`config`文件夹中
### 启动
```
bash cmd/start.sh
```
## 本项目参与者

## 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request