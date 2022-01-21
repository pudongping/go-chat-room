package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	// 在 listen 时没有指定 ip，表示绑定到当前机器的所有 IP 上，如果要指定 ip 的话
	// 则：net.Listen("tcp", "192.168.10.25:2020")
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	// 用于广播消息
	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}

}

type User struct {
	ID             int         // 用户的唯一标识
	Addr           string      // 用户的 IP 地址和端口
	EnterAt        time.Time   // 用户进入的时间
	MessageChannel chan string // 给当前用户发送消息的通道
}

func (u *User) String() string {
	return u.Addr + ", UID:" + strconv.Itoa(u.ID) + ", Enter At:" +
		u.EnterAt.Format("2006-01-02 15:04:05+8000")
}

// 给用户发送的消息
type Message struct {
	OwnerID int    // 发送消息者的 id
	Content string // 消息内容
}

var (
	// 新用户到来，通过该 channel 进行登记
	enteringChannel = make(chan *User)
	// 用户离开，通过该 channel 进行登记
	leavingChannel = make(chan *User)
	// 广播专用的用户普通消息 channel，缓冲是尽可能避免出现异常情况堵塞，这里简单给了 8，具体值根据情况调整
	messageChannel = make(chan Message, 8)
)

// broadcaster 用于记录聊天室用户，并进行消息广播：
// 1. 新用户进来；2. 用户普通消息；3. 用户离开
func broadcaster() {
	// 存储在线的用户
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			// 新用户进入
			users[user] = struct{}{}
		case user := <-leavingChannel:
			// 用户离开
			delete(users, user) // 从 map 中删除掉用户
			// 避免 goroutine 泄漏
			close(user.MessageChannel)
		case msg := <-messageChannel: // 全局的 messageChannel 用来给聊天室所有用户广播消息
			// 给所有在线用户发送消息
			for user := range users {
				if user.ID == msg.OwnerID { // 当发送者为自己时，则不推送消息
					continue
				}
				user.MessageChannel <- msg.Content
			}
		}
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// 1. 新用户进来，构建该用户的实例
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(), // ip 地址
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	// 2. 当前在一个新的 goroutine 中，用来进行读操作，因此需要开一个 goroutine 用于写操作
	// 读写 goroutine 之间可以通过 channel 进行通信
	// 这里的 sendMessage 在一个新的 goroutine 中，如果函数里的 ch 不关闭，该 goroutine 是不会退出的
	// 因此需要注意不关闭 ch 导致的 goroutine 泄漏的问题
	go sendMessage(conn, user.MessageChannel)

	// 3. 给当前用户发送欢迎信息
	// 同时给聊天室所有用户发送有新用户到来的提醒
	user.MessageChannel <- "Welcome, " + user.String()
	msg := Message{
		OwnerID: user.ID,
		Content: "user:`" + strconv.Itoa(user.ID) + "` has enter",
	}
	messageChannel <- msg

	// 4. 将该记录到全局的用户列表中，避免用锁
	enteringChannel <- user

	// 控制超时用户踢出
	var userActive = make(chan struct{})
	go func() {
		d := 1 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	// 5. 循环读取用户的输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		messageChannel <- msg

		// 用户活跃
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}

	// 6. 用户离开，需要做登记，并给聊天室其他用户发通知
	leavingChannel <- user
	msg.Content = "user:`" + strconv.Itoa(user.ID) + "` has left"
	messageChannel <- msg

}

// 只从 channel 中读数据
func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

// 生成用户 ID
var (
	globalID int
	idLocker sync.Mutex
)

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalID++
	return globalID
}
