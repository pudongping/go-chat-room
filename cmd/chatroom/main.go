package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/pudongping/go-chat-room/global"
	"github.com/pudongping/go-chat-room/server"
)

var (
	addr   = ":2022"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |

	Go-Chat-Room，start on %s

`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner+"\n", addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
