package main

import (
	"io"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建用户api
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听当前user channel消息的 goroutine
	go user.ListenMessage()

	return user
}

// 上线通知
func (this *User) Online() {
	// 用户上线
	this.server.mapLock.Lock()
	this.server.OnLineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "上线")
}

// 下线通知
func (this *User) Offline() {
	// 用户下线
	this.server.mapLock.Lock()
	delete(this.server.OnLineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "下线")

}

// 返回给当前客户端消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 发消息通知
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnLineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + " 在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {

		//this.SendMsg("请输入您要修改后的名字：\n")
		//name := readMsg(this.conn)

		// 字符串截取
		newName := strings.Split(msg, "|")[1]

		// 判断newname 是否已经被占用
		_, ok := this.server.OnLineMap[newName]

		if ok {
			this.SendMsg("改名字已经被占用")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnLineMap, this.Name)
			this.server.OnLineMap[newName] = this
			this.server.mapLock.Unlock()
			//this.server.BroadCast(this, "修改名字成功，新名字为："+name)
			this.SendMsg("修改名字成功")
		}
	} else {
		this.server.BroadCast(this, msg)
	}

}

type NullString struct {
	value string
	len   int
	valid bool
}

func readMsg(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		return ""
	}
	return string(buffer[:n-1])
}

// 监听当前User channel 的方法，有消息就发送给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		// 把当前读到的消息发到客户端
		this.conn.Write([]byte(msg + "\n"))
	}
}
