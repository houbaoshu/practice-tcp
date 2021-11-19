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
	"github.com/spf13/cobra"
	"math"
	"net"
	"time"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the application",
	Long: `This is a tcp client and the default port is 8081.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}
var (
	conc int    // 并发数
	N    int    // 信息量
	msg  string // 信息内容
)

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
	runCmd.Flags().IntVarP(&conc, "concurrency", "c", 10, "Number of workers to run concurrently. Total number of requests cannot be smaller than the concurrency level. ")

	// 添加flag --number/-n type: int
	runCmd.Flags().IntVarP(&N, "number", "n", 200, "The number of message to send.")

	// 添加flag --message/-s type: string (信息必须要有)
	runCmd.Flags().StringVarP(&msg, "message", "s", "", "The content of message (required)")
	runCmd.MarkFlagRequired("message")
}

type Work struct{
	count int
	testStartTime time.Time
	testEndTime time.Time
	QPS float64
	P90 time.Duration
}

func (w *Work) print() {
	s := w.testStartTime
	e := w.testEndTime
	fmt.Printf("test start time: %v:%v:%v.%v\n",s.Hour(), s.Minute(), s.Second(), s.Nanosecond()/1e6)
	fmt.Printf("test start time: %v:%v:%v.%v\n",e.Hour(), e.Minute(), e.Second(), e.Nanosecond()/1e6)
	fmt.Printf("QPS : %v [#/sec]\n", int(w.QPS))
	fmt.Println("Percentage the requests served with a certain time(ms)")
	fmt.Printf("P90 : %v ms\n", w.P90.Milliseconds())
}
func run() {
	work := make(chan *Work)
	for i := 0; i < conc; i++ {
		go worker(i, work)
	}

	work <- new(Work)
	first := <-work
	first.testStartTime = time.Now()
	work <- first

	for {
		select {
		case check := <- work:
			if  math.Abs(0.9 - float64(check.count)/float64(N)) < 0.02 { //抽样检查有可能得不到, time.Duration为空的概念没有？
				check.P90  = time.Since(check.testStartTime)
			}
			if check.count >= N {
				check.testEndTime = time.Now()
				check.QPS = float64(check.count) / (check.testEndTime.Sub(check.testStartTime)).Seconds()
				check.print()
				return
			}
			work <- check
		}
	}

}

func worker(i int, work chan *Work) {
	for {
		get := <-work //接收工作
		// 拨号
		dail()
		get.count++
		//time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)//一般为设为100ms 80ms 50 ms
		work <- get //交接工作
	}
}

func dail() {
	conn, err := net.Dial("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
}
