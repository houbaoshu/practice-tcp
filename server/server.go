package main

import (
	"fmt"
	"net"
	"time"
)

type Log struct {str string}

func main() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	// 处理连接
	c := handleConnection(ln)

	for {
		s := <-c
		fmt.Println(s.str)
	}

}

func handleConnection(ln net.Listener) <-chan Log{
	c := make(chan Log)

	go func() {
		for{
			// 获取连接
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}
			// 关闭连接
			defer conn.Close()
			// 获取连接时间
			t := time.Now()

			msg := make([]byte, 10)
			_, err = conn.Read(msg)
			if err != nil {
				panic(err)
			}
			c <- Log{fmt.Sprintf("[%d:%d:%d %d:%d:%d.%d] recv %v msg, msgid: %v, msg content: %v",
				t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Minute(),t.Nanosecond() / 1e6, conn.RemoteAddr(),t.UnixNano(), string(msg))}
		}
	}()


	return c
}