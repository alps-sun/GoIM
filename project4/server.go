package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	user := NewUser(conn, this)
	user.Online()

	// 监听用户是否活跃 超时强踢
	isLive := make(chan bool)

	// 接收客户端端请求
	go func() {

		buffer := make([]byte, 4096)

		for {
			n, err := conn.Read(buffer)

			if err != nil && err != io.EOF {
				fmt.Println("Conn read err...")
				return
			}

			if n == 0 {
				user.Offline()
				return
			}

			msg := string(buffer[:n-1])

			// 用户对msg 进行处理
			user.DoMessage(msg)
			isLive <- true
		}

	}()

	// 当前handler阻塞, select 随机执行一个可运行的 case。如果没有 case 可运行，它将阻塞，直到有 case 可运行

	for {
		select {
		case <-isLive:
			// 当前用户是活跃的，应该重置定时器
			// 不做任何操作，为了激活select,更新下面的定时器
			// 第一个case走完之后，还会继续执行地79行的判断，  自然就执行了time.After(time.Second * 10) 这一句
		case <-time.After(time.Second * 1500):
			user.SendMsg("你被踢了\n")

			//销毁资源
			close(user.C)

			//关闭conn
			conn.Close()

			//退出当情连接
			return

		}
	}

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
