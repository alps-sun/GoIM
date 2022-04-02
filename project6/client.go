package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
	fmt.Scanln(&flag)

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
			this.Group()
			break
		case 3:
			this.UpdataName()
			break
		case 0:
			fmt.Println("退出")

		}
	}
}

func (this *Client) UpdataName() bool {
	fmt.Println("请输入您好修改的名字：")

	// 控制台的输入用client.Name 接收
	fmt.Scanln(&this.Name)

	sendMsg := "rename|" + this.Name + "\n"

	_, err := this.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return false
	}
	return true
}

func (this *Client) Group() {
	var sendMsg string
	fmt.Println("开始群聊,请输入内容， exit 退出 ：")
	fmt.Scanln(&sendMsg)

	for sendMsg != "exit" {
		if len(sendMsg) != 0 {
			_, err := this.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write error:", err)
				break
			}
		}
		sendMsg = ""
		fmt.Println("开始群聊,请输入内容， exit 退出 ：")
		fmt.Scanln(&sendMsg)

	}

}

// 处理server 回应的消息， 直接显示到标准的输出即可
func (this *Client) DealResposne() {
	// 一旦client.conn 有数据，就直接打印到stdout 标准输出上，永久阻塞监听
	io.Copy(os.Stdout, this.Conn)
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

	// 单独开启一个goroutine 处理server 的回执消息
	go client.DealResposne()
	// 启动客户端业务
	client.Run()

}
