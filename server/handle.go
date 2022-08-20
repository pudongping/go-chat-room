package server

import (
	"net/http"

	"github.com/pudongping/go-chat-room/logic"
)

func RegisterHandle() {
	// 广播消息处理
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)              // 首页
	http.HandleFunc("/user_list", userListHandleFunc) // 在线用户列表
	http.HandleFunc("/ws", WebSocketHandleFunc)       // 处理 ws 服务
}
