package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	Conn       net.Conn
	flag       int
}

func NewClient(serverIP string, serverPort int) *Client {

	// 创建客户度
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       999,
	}

	// 连接服务端
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))

	if err != nil {
		fmt.Println("net Dial error : ", err)
		return nil
	}
	client.Conn = conn
	client.Name = conn.RemoteAddr().String()
	// 返回客户端
	return client
}

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "请输入IP地址， 例如127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "请输入端口号，例如 8888")

}

func (this *Client) menu() bool {
	var flag int
	fmt.Println("1、单聊")
	fmt.Println("2、公聊")
	fmt.Println("3、更新名字")
	fmt.Println("0、退出")

	// 接收客户端的控制台输入
	fmt.Scan(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("请输入合法数值")
		return false
	}
}

func (this *Client) Run() {
	for this.flag != 0 {
		// 输入不合法就永远循环
		for this.menu() != true {
		}
		// 根据不同模式处理不同业务
		switch this.flag {
		case 1:
			fmt.Println("开启单聊模式")
			break
		case 2:
			fmt.Println("开启群聊模式")
			break
		case 3:
			fmt.Println("修改名字")
			break
		case 0:
			fmt.Println("退出")

		}
	}
}

func main() {

	// 命令行解析
	flag.Parse()

	client := NewClient(serverIP, serverPort)

	if client == nil {
		fmt.Println(">>>> 连接服务器失败...")
		return
	}
	fmt.Println(">>> 连接服务器成功...")

	// 启动客户端业务
	client.Run()

}
