# go-chat-room

使用 golang 写的一个简单的聊天室服务

> 这套代码里面有 3 套小服务

## 目录结构

```shell

├── LICENSE
├── README.md  -- 说明文档
├── cmd  -- 命令行目录
│   ├── chatroom  -- 基于浏览器作为 websocket 客户端的聊天室目录
│   │   └── main.go  -- 入口文件
│   ├── tcp  -- tcp 聊天室目录
│   │   ├── client.go  -- tcp 聊天室 websocket 客户端
│   │   └── server.go  -- tcp 聊天室 websocket 服务端
│   └── websocket  -- websocket 第三方插件包 demo 代码目录
│       ├── client.go  -- nhooyr.io/websocket 库 websocket 客户端
│       ├── server_gorilla.go  -- gorilla/websocket 库 websocket 服务端
│       └── server_nhooyr.go  -- nhooyr.io/websocket 库 websocket 服务端
├── config  -- 配置文件目录
│   └── chatroom.yaml  -- 配置文件
├── global  -- 全局变量目录
│   ├── config.go  -- 读取配置信息
│   └── init.go  -- 初始化变量
├── go.mod
├── go.sum
├── logic  -- 逻辑层
│   ├── broadcast.go  -- 广播器
│   ├── message.go  -- 各种消息类型
│   ├── offline.go  -- 发送离线消息
│   ├── sensitive.go  -- 过滤敏感词（直接暴力的替换）
│   └── user.go  -- 用户生成 token、解析 token、接收消息
├── server  -- 服务层
│   ├── handle.go  -- 路由处理
│   ├── home.go  -- http 路由对应的处理方法
│   └── websocket.go  -- websocket 路由对应的处理方法
└── template  -- 模版层
    └── home.html  -- 聊天室前端页面

```

## 1. 基于 tcp 协议的简单聊天室

> 所有涉及到的代码都在 `cmd/tcp` 目录下

### 代码测试方式

```shell

# 先启动 websocket 服务端
go run cmd/tcp/server.go 


# 以下分别开启 3 个终端进行客户端连接
# 终端 1 中启动 websocket 客户端
go run cmd/tcp/client.go
# output is：
# Welcome, 127.0.0.1:62841, UID:1, Enter At:2022-08-18 17:30:24+8000
# user:`2` has enter
# user:`3` has enter

# 终端 2 中启动 websocket 客户端
go run cmd/tcp/client.go
# output is：
# Welcome, 127.0.0.1:62846, UID:2, Enter At:2022-08-18 17:30:47+8000
# user:`3` has enter

# 终端 3 中启动 websocket 客户端
go run cmd/tcp/client.go
# output is：
# Welcome, 127.0.0.1:62854, UID:3, Enter At:2022-08-18 17:31:09+8000

```

## 2. 使用第三方库写的一个简要 demo（基于 websocket 协议）

> 所有涉及到的代码都在 `cmd/websocket` 目录下

这里主要是采用 `nhooyr.io/websocket` 库写的一个 `websocket` 服务端和客户端，同时也提供了 `gorilla/websocket` 库的 `websocket` 服务端代码写法。
在实际项目中，还是推荐使用 `nhooyr.io/websocket` 库。

### 代码测试方式

```shell

# 先启动 websocket 服务端
go run cmd/websocket/server_nhooyr.go 
# output is：
# 2022/08/18 17:01:20 接收到客户端：Hello WebSocket Server

# 启动 websocket 客户端
go run cmd/websocket/client.go
# output is：
# 接收到服务端响应：Hello WebSocket Client

```

## 3. 基于浏览器作为客户端的聊天室（基于 websocket 协议）

> 项目入口文件在 `cmd/chatroom` 目录中 `main.go` 文件

### 代码测试方式

```shell

# 先启动 websocket 服务端
go run cmd/chatroom/main.go
# output is：
# 
#    ____              _____
#   |    |    |   /\     |
#   |    |____|  /  \    | 
#   |    |    | /----\   |
#   |____|    |/      \  |
#
#        Go-Chat-Room，start on :2022
#
#

```

### 多使用几个浏览器访问 `127.0.0.1:2022`

> 需要先输入一个昵称，然后加入聊天室，才能够发送消息

#### 浏览器 1

![image.png](https://upload-images.jianshu.io/upload_images/14623749-d728e7741cf05df5.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

#### 浏览器 2

![image.png](https://upload-images.jianshu.io/upload_images/14623749-56ed8b7131eb5e2a.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

#### 终端显示为

![image.png](https://upload-images.jianshu.io/upload_images/14623749-7233fbf12bc9363f.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

### 聊天室性能测试

```shell

# 编译成二进制文件
go build -o chatroom cmd/chatroom/main.go

# 启动聊天室
./chatroom
 
# 性能测试
# 尝试 10 个用户同时进入聊天室，并每 20s 各发送一条消息
go run cmd/benchmark/main.go -u 10 -m 20s -l 0  
# 尝试 200 个用户每个用户间隔 10ms 进入聊天室，并每隔 20s 各发送一条消息
go run cmd/benchmark/main.go -u 200 -m 20s -l 10ms  

```