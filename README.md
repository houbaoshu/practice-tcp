# 使用指南
### 服务器启动
```shell
./server
```

### 客户端启动

```shell
./client run -c 10 -n 1000 -s "hello, tcp!"
```


#### Flags:

* -c, --concurrency int   Number of workers to run concurrently. Total number of requests cannot be smaller than the concurrency level. . (default 10)
*   -h, --help              help for run
*   -s, --message string    The content of message (required)
*   -n, --number int        The number of message to send. Default is 200. (default 200)
