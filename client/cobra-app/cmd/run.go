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
	"math/rand"
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
		run()
	},
}
var (
	// 并发数
	conc int
	// 信息量
	N int
	// 信息内容
	msg string
)


type Report struct {
	//测试开始时间
	testStartTime time.Time
	//测试结束时间
	testEndTime time.Time
	//持续时间
	duration time.Duration
	//总请求数
	N int
	//请求成功数
	n int
	QPS int64
	P50 time.Duration
	P60 time.Duration
	P70 time.Duration
	P80 time.Duration
	P90 time.Duration
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
	runCmd.Flags().IntVarP(&N, "number", "n", 200, "The numbur of message to send. Default is 200.")

	// 添加flag --message/-s type: string (信息必须要有)
	runCmd.Flags().StringVarP(&msg, "message", "s", "", "The content of message (required)")
	runCmd.MarkFlagRequired("message")
}


func (r *Report) work(worker int, quit chan bool){
	r.N++
	timeOut := time.After(80 * time.Millisecond)//80ms连接不上就退出
	go func() {
		for i := 0; ; i++ {
			select {
			case <-quit:
				return
			case <- timeOut:
				return
			default:
				// 创建拨号器来创建连接
				conn, err := net.Dial("tcp", ":8081")
				// 处理.Dial错误
				if err != nil {
					panic(err)
				}
				// 处理连接后关闭
				defer conn.Close()
				// 处理连接
				_, err = conn.Write([]byte(msg))
				// 发送消息失败处理错误
				if err != nil {
					panic(err)
				}
				r.n++
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) //休眠100ms
			}
		}
	}()
}





func  run() {
	var r *Report
	// 创建连接交给conc个工人做， 谁做的快就提交谁的报告
	r.testStartTime = time.Now()
	quit := make(chan bool)
	// 创建10个并发, 10个工人
	for i := 0; i < conc; i++ {
		r.work(i, quit)
	}

	if r.N == N {
		quit <- true
	}

	r.testEndTime = time.Now()

}




