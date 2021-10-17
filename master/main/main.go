package main

import "runtime"

func initEnv() {
	// 设置线程数等于核心数
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// 初始化线程
	initEnv()

	// 启动 Http 服务
}
