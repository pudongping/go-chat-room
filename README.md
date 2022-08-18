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

## 1. 使用第三方库写的一个简要 demo

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

