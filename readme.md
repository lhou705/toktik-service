# toktik
## 介绍
第五届字节跳动青训营汪汪大队作品

## 软件架构
本项目采用服务架构，使用[consul](https://www.consul.io/)作为注册中心，[kitex](https://www.cloudwego.io/zh/docs/kitex/)框架搭建。分为用户、通信、对象存储、网关、视频五个部分。其中网关部分使用[hertz](https://www.cloudwego.io/zh/docs/hertz/)搭建。
文件存储使用腾讯云对象存储。视频封面在上传视频后自动生成。为了加快访问速度，还使用了cdn技术加速访问。

已部署在服务器上。地址`htps://toktik.lhou.ltd/`。客户端在[这里](https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7)获取。

## 程序目录
```
|
|-- bin             // 编译后程序的可运行文件    
|-- cmd             // 一些启动、关闭的命令
|-- config          // 程序的配置文件
|-- log             // 程序的运行日志
|-- service         // 服务主体文件夹
|      |-- cos      // 文件上传微服务，主要负责将文件上传至对象存储
|      |-- message  // 消息微服务，主要负责好友之间的消息传递
|      |-- user     // 用户微服务，主要负责用户相关的功能，包括注册登陆、关注等
|      |-- video    // 视频微服务，主要负责视频相关的功能，包括投递视频、获取视频列表、视频的点赞评论等
|      |-- web      // 网关微服务，主要负责接收网络请求并进行数据处理
|-- sql             // 建表的SQL语句
|-- thrift          // 定义的接口文件
```

## 安装教程
### 必要使用环境安装
1. 安装mysql数据库（[安装教程](https://dev.mysql.com/doc/mysql-installation-excerpt/8.0/en/)），或者使用兼容mysql的云数据库。
2. 安装consul（[安装教程](https://developer.hashicorp.com/consul/downloads)）
### 安装本项目
#### 克隆本项目
- gitee
```shell
git clone https://gitee.com/lhou/toktik.git
```
- github 
```shell
git clone https://github.com/lhou705/toktik-service.git
```
#### 编译
```shell
bash cmd/build.sh all
```
注意，此脚本支持编译单个服务。目前支持的服务有`cos`、`message`、`video`、 `user`、`web`服务。使用下面的命令可以只编译`cos`服务。
```shell
bash cmd/build.sh cos
```
## 使用
### 设置配置文件
请参考[文档](docs/config.md)。准备好这些文件后，将其放入`config`文件夹中
### 启动
```shell
bash cmd/start.sh all
```
注意，此脚本支持启动单个服务。目前支持的服务有cos、message、video、 user、web服务。使用下面的命令可以只启动cos服务。
```shell
bash cmd/start.sh cos
```
### 停止
```shell
bash cmd/stop.sh all
```
注意，此脚本支持停止单个服务。目前支持的服务有cos、message、video、 user、web服务。使用下面的命令可以只停止cos服务。
```shell
bash cmd/stop.sh cos
```