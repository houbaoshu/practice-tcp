package main

// #go run go client.go run -c 10 -n 10000 -s "hello, tcp!"
import (
	"fmt"
	"net"
	"os" //接收终端传来的参数
)

func handleConnection(conn net.Conn) {
	// 处理连接后关闭
	defer conn.Close()

	// 发送消息
	_, err := conn.Write([]byte("hello tcp!"))
	// 发送消息失败处理错误
	if err != nil {
		panic(err)
	}
	// 打印发送消息成功提示
	fmt.Println("Send message success!")
}
func run() {
	for i, v := range os.Args[2:] {
		switch v {
		case "-c":
			// 有-c标识，为设置通道数，读取其后的值
			channel := os.Args[i+1]
			fmt.Println(channel)
		case "-n":
			// 有-n标识，为设置信息数，读取其后的值
			amount := os.Args[i+1]
			fmt.Println(amount)
		case "-s":
			// 有-s标识，为设置信息内容，读取其后的值
			message := os.Args[i+1]
			fmt.Println(message)
		default:
		}
	}
	// 创建拨号器来创建客户端
	conn, err := net.Dial("tcp", ":8081")
	// 处理.Dial错误
	if err != nil {
		panic(err)
	}

	handleConnection(conn)

}
func main() {
	//从终端接收参数，接收到run才执行
	// just print params
	switch os.Args[1] { //第一个参数是文件名称，所以这是第二个参数
	case "run":
		fmt.Printf("Running... %v\n", os.Args[2:])
		run()
	default:
		panic("Bad Commmand!")
	}
}
