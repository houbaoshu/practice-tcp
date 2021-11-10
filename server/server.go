package main

import (
	//"bofio"//需要吗？
	"fmt"
	"log"
	"net"
	"time"
)

// 处理连接
func handleConnection(conn net.Conn) {
	// 处理连接后关闭连接
	defer conn.Close()
	// 用一个长度为10的类型为byte的slice来接收消息
	msg := make([]byte, 10)
	// 从连接获得消息
	_, err := conn.Read(msg)
	// 处理.Read()错误
	if err != nil {
		fmt.Printf("Read() failed, err: %v\n", err)
	}
	//打印消息
	log.Printf("recv %v msg, msgid: %v msg content: %v\n", conn.RemoteAddr(), time.Now().UnixNano(), string(msg))

}

func main() {
	// 通过创建监听器来创建服务器，端口号为8081
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// 处理.Listen()错误
		fmt.Printf("Listen() failed, err: %v\n", err)
	}

	for {
		// 通过监听器获得一个连接
		conn, err := ln.Accept()
		if err != nil {
			// 处理.Accept()错误
			fmt.Printf("Accept() failed, err: %v\n", err)
		}
		// 开一个协程处理连接
		go handleConnection(conn)
	}

}
