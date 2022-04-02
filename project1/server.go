package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      int
	Message   chan string
	OnLineMap map[string]*User

	// 锁
	mapLock sync.RWMutex
}

// 新建一个服务器
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		Message:   make(chan string),
		OnLineMap: make(map[string]*User),
	}

	return server
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("业务处理开始")

	// 用户上线
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnLineMap[user.Name] = user
	this.mapLock.Unlock()

	// 上线用户广播
	this.BroadCast(user, " 用户上线")

	// 接收客户端端请求
	go func() {

		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)

		if err != nil && err != io.EOF {
			return
		}

		if n == 0 {
			this.BroadCast(user, "下线")
		}

		msg := string(buffer[:n-1])
		this.BroadCast(user, msg)

	}()

	// 当前handler阻塞, select 随机执行一个可运行的 case。如果没有 case 可运行，它将阻塞，直到有 case 可运行
	select {}
}

// 监听message 的goroutine,用户通知发送广播消息，通知用户上线
func (this *Server) LitenerMessage() {
	for {
		msg := <-this.Message

		// 将message 发送给在线用户
		this.mapLock.Lock()
		for _, cli := range this.OnLineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 启动服务端接口
func (this *Server) Start() {
	// 监听端口
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println("start listen  error:", err)
		return
	}

	// 关闭连接
	defer listen.Close()

	// 启动监听 Message 的 goroutine
	go this.LitenerMessage()

	for {
		// accept
		conn, err := listen.Accept()

		if err != nil {
			fmt.Println("start accept error : ", err)
			continue
		}

		fmt.Println("开始建立连接....")

		// go程  处理业务，然后主线程继续监听
		go this.Handler(conn)
	}
}
