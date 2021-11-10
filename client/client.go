package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	// 处理连接后关闭
	defer conn.Close()
	// 发送消息
	_, err := conn.Write([]byte("hello tcp!"))
	// 发送消息失败处理错误
	if err != nil {
		fmt.Printf("Send message error: %v\n", err)
	}
	// 打印发送消息成功提示
	fmt.Println("Send message success!!!")
}

func main() {
	// 创建拨号器来创建客户端
	conn, err := net.Dial("tcp", ":8081")
	// 处理.Dial错误
	if err != nil {
		fmt.Printf(".Dial error: %v", err)
	}

	handleConnection(conn)
}
