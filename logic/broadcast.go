package logic

import (
	"expvar"
	"fmt"
	"log"

	"github.com/pudongping/go-chat-room/global"
)

func init() {
	expvar.Publish("message_queue", expvar.Func(calcMessageQueueLen))
}

func calcMessageQueueLen() interface{} {
	fmt.Println("===len=:", len(Broadcaster.messageChannel))
	return len(Broadcaster.messageChannel)
}

// broadcaster 广播器
type broadcaster struct {
	// 所有聊天室用户
	users map[string]*User

	// 所有 channel 统一管理，可以避免外部乱用

	// 用户进入聊天室时，通过该 channel 告知 Broadcaster，即将该用户加入 Broadcaster 的 users 中
	enteringChannel chan *User // 将该用户加入广播器的用户列表中
	// 用户离开聊天室时，通过该 channel 告知 Broadcaster，
	// 即将该用户从 Broadcaster 的 users 中删除，同时需要关闭该用户对应的 messageChannel，避免 goroutine 泄露
	leavingChannel chan *User // 用户离开
	// 用户发送的消息，通过该 channel 告知 Broadcaster，之后 Broadcaster 将它发送给 users 中的用户
	messageChannel chan *Message

	// 判断该昵称用户是否可进入聊天室（重复与否）：true 能，false 不能
	// 用来接收用户昵称，方便 Broadcaster 所在 goroutine 能够无锁判断昵称是否存在
	checkUserChannel chan string
	// 用来回传该用户昵称是否已经存在
	checkUserCanInChannel chan bool

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*User
}

// Broadcaster 实例化一个广播实例，方便外部调用
var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, global.MessageQueueLen),

	// 虽然没有显示使用锁，但这里要求 checkUserChannel 必须是无缓冲的，否则判断可能会出错
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*User),
}

// Start 启动广播器
// 需要在一个新 goroutine 中运行，因为它不会返回
// - select 关键字和 { 之间不允许存在任何表达式和语句；
//
// - fallthrough 语句不能使用；
//
// - 每个 case 关键字后必须跟随一个 channel 接收数据操作或者一个 channel 发送数据操作，所以叫做专门为 channel 设计的；
//
// - 所有的非阻塞 case 操作中将有一个被随机选择执行（而不是按照从上到下的顺序），然后执行此操作对应的 case 分支代码块；
//
// - 在所有的 case 操作均阻塞的情况下，如果 default 分支存在，则 default 分支代码块将得到执行； 否则，当前 goroutine 进入阻塞状态；
func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			// 新用户进入
			b.users[user.NickName] = user

			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			// 用户离开
			delete(b.users, user.NickName)
			// 避免 goroutine 泄露（用户退出时，一定要关闭 u.MessageChannel 这个 channel）
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:
			// 给所有在线用户发送消息
			for _, user := range b.users {
				// 群发消息的时候，消息内容不发给自己
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
			OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

// UserEntering 将该用户加入广播器的用列表中
func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

// UserLeaving 用户离开
func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

// Broadcast 广播
// 广播器所在 goroutine 叫 broadcaster goroutine
func (b *broadcaster) Broadcast(msg *Message) {
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("broadcast queue 满了")
	}
	b.messageChannel <- msg
}

// CanEnterRoom 判断该昵称用户是否可进入聊天室（重复与否）：true 能，false 不能
func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
