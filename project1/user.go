package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// 创建用户api
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}
	// 启动监听当前user channel消息的 goroutine
	go user.ListenMessage()

	return user
}

// 监听当前User channel 的方法，有消息就发送给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		// 把当前读到的消息发到客户端
		this.conn.Write([]byte(msg + "\n"))
	}
}
