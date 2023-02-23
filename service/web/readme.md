# web模块
## 文件结构
```plantuml
|-- client  // 微服务客户端程序初始化    
|-- common     // 公共模块
|-- config   // 配置文件模型
|-- handler  // 处理函数
|-- kitex_gen     // 客户端文件，kitex框架自动生成     
|-- mw // 中间件
|-- router // 路由
|-- server //服务器相关，用于初始化web服务
|-- utils // 工具类
|-- go.mod     // 模块文件，初始化包时生成     
|-- go.sum     // 模块文件，初始化包时生成   
|-- main.go    // 程序入口函数
```