# 配置设置
### 网关的配置文件`toktik_web.config.json`
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
    "addr": "your consul server addr",
    "token": "your access token for consul. optional"
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
### 用户、信息模块配置文件`toktik_user.config.json`、`toktik_message.config.json`
```json lines
{
  "server": {
    "name": "your server name",
    "addr": "your server addr",
    "reusePort": true,
    "readWriteTimeOut": 10,
    "registerAddr": "your consul addr",
    "token": "your access token for consul. optional"
  },
  "mysql": {
    "addr": "your mysql server addr",
    "user": "your user",
    "password": "your password",
    "database": "your database"
  }
}
```
### 视频模块配置文件`toktik_video.config.json`
```json lines
{
  "server": "同用户、信息模块配置文件",
  "mysql": "同用户、信息模块配置文件",
  "cos": {
    "name": "对象存储的服务名称",
    "cdnAddr": "cdn加速域名"
  }
}
```
### 腾讯云对象存储配置文件`toktik_cos.config.json`
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