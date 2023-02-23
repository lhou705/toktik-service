# video模块
## 文件结构
```plantuml
|-- kitex_gen  // 程序的执行代码，由框架自动的生成    
|-- script     // 程序的启动命令，由框架自动的生成
|-- build.sh   // 编译命令，由框架自动生成
|-- config.go  // 配置文件类
|-- go.mod     // 模块文件，初始化包时生成     
|-- go.sum     // 模块文件，初始化包时生成     
|-- handler.go // 逻辑处理函数，框架自动生成，需要自行编写具体的处理逻辑     
|-- kitex.yaml // 框架的版本和服务名，框架自动生成     
|-- main.go    // 程序入口函数
|-- model.go   // 模型文件     
|-- readme.md 
```