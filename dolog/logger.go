// 初始化日志记录器
// Info, Warn, Error := doLog.InitLog(path, doLog.DefaultFormat)
package dolog

import (
	"io"
	"log"
	"os"
)

// 默认的日志的格式
// INFO: 2021/04/03 01:51:03 main.go:140: 本次图片保存完毕
const DefaultFormat = log.Ldate | log.Ltime | log.Lshortfile

var (
	i *log.Logger
	w *log.Logger
	e *log.Logger
)

// 初始化日志记录器
// @param path string 日志的路径
func InitLog(path string, format int) (*log.Logger, *log.Logger, *log.Logger) {
	//日志输出文件
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		e.Fatalln(err)
	}
	//自定义日志格式
	i = log.New(io.MultiWriter(file, os.Stdout), "INFO: ", format)
	w = log.New(io.MultiWriter(file, os.Stdout), "WARN: ", format)
	e = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", format)
	return i, w, e
}
