package main

import (
	"math/rand"
	"net"
	"time"
)

type Work struct{ count int }

func main() {
	work := make(chan *Work)
	for i := 0; i < 10; i++ {
		go worker(i, work)
	}

	work <- new(Work)
	for {
		select {
		case s:= <- work:
			if s.count >= 100 {
				return
			}
			work <- s
		}
	}
	<-work
}

func worker(i int, work chan *Work) {
	for {
		get := <-work //接收工作
		// 拨号
		dail()
		get.count++
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		work <- get //交接工作
	}
}

func dail() {
	conn, err := net.Dial("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("msg"))
	if err != nil {
		panic(err)
	}
}
