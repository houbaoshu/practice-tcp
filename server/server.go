package main

import (
	"fmt"
	"net"
)

func main() {
	// create a Listener
	ln, err := net.Listen("tcp", ":8081")
	// handle err
	if err != nil {
		panic(err)
	}

	for i:=0;;i++{
		// get a connection
		conn, err := ln.Accept()
		// handle err
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		msg := make([]byte, 10)
		_, err = conn.Read(msg)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(msg), i)

	}

}
