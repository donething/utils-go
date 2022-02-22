// Package doint 接收中断信号，处理关闭程序前的任务
package doint

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Init 初始化中断处理程序
func Init(callback func()) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		callback()
		os.Exit(0)
	}()
}
