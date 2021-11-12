/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		// fmt.Println("Run: " + strings.Join(args, " "))
		fmt.Printf("channel: %v, number: %v, message: %v", channel, number, message)
		// 开启任务
		run()
	},
}

// Max size of the buffer of result channel.
const MaxResult = 1000000
const maxIdleConn = 500

var (
	// 并发数
	conc int
	// 信息数
	number int
	// 信息内容
	message string
)

type result struct {
	startTime time.Time
	endTime   time.Time
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 添加flag --concurrency/-c type: int
	// Number of workers to run concurrently. Total number of requests cannot be smaller than the concurrency level. Default is 10.
	runCmd.Flags().IntVarP(&conc, "concurrency", "c", 10, "Number of workers to run concurrently. Total number of requests cannot be smaller than the concurrency level. Default is 10.")

	// 添加flag --number/-n type: int
	runCmd.Flags().IntVarP(&number, "number", "n", 200, "The numbur of message to send. Default is 200.")

	// 添加flag --message/-s type: string (信息必须要有)
	runCmd.Flags().StringVarP(&message, "message", "s", "", "The content of message (required)")
	runCmd.MarkFlagRequired("message")
}

func handleConnection(conn net.Conn) {
	// 处理连接后关闭
	defer conn.Close()

	// 发送消息
	_, err := conn.Write([]byte(message))
	// 发送消息失败处理错误
	if err != nil {
		panic(err)
	}
	// 打印发送消息成功提示
	fmt.Println("Send message success!")
}

func run() {
	//创建信道
	results := make(chan *result, min(conc*1000, maxIdleConn))
	stopCh := make(chan struct{}, conc)

	// 创建拨号器来创建客户端
	conn, err := net.Dial("tcp", ":8081")
	// 处理.Dial错误
	if err != nil {
		panic(err)
	}

	handleConnection(conn)

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
