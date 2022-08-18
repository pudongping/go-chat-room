# go-chat-room

使用 golang 写的一个简单的聊天室服务

> 这套代码里面有 3 套小服务

## 目录结构

```shell

├── LICENSE
├── README.md
├── cmd
│   ├── chatroom
│   │   └── main.go
│   ├── tcp
│   │   ├── client.go
│   │   └── server.go
│   └── websocket
│       ├── client.go
│       ├── server_gorilla.go
│       └── server_nhooyr.go
├── config
│   └── chatroom.yaml
├── global
│   ├── config.go
│   └── init.go
├── go.mod
├── go.sum
├── logic
│   ├── broadcast.go
│   ├── message.go
│   ├── offline.go
│   ├── sensitive.go
│   └── user.go
├── server
│   ├── handle.go
│   ├── home.go
│   └── websocket.go
└── template
    └── home.html

```

## 1. 基于 tcp 的简单聊天室

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

## 2. 使用第三方库写的一个简要 demo

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

