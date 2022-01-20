package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})

	// 新开了一个 goroutine 用于接收消息
	go func() {
		// 通过 io.Copy 来操作 IO，包括从标准输入读取数据写入 TCP 连接中，以及从 TCP 连接中读取数据写入标准输出
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		// 新开的 goroutine 通过一个 channel 来和 main goroutine 通讯
		done <- struct{}{} // signal the main goroutine
	}()

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done

}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
