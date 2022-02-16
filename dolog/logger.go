// 初始化日志记录器
// 参考：https://www.jianshu.com/p/a9427a4e2ada
// 用法，初始化：创建 logger go文件后：
// var (
//  Info  *log.Logger
//  Warn  *log.Logger
//  Error *log.Logger
// )
// const LogName = "run.log"
// func init() {
//  Info, Warn, Error = dolog.InitLog(LogName, dolog.DefaultFormat)
// }
// 使用：logger.Info.Printf("创建配置文件：%s\n", confPath)

package dolog

import (
	"io"
	"log"
	"os"
)

// DefaultFormat 默认的日志的格式
//
// INFO: 2021/04/03 01:51:03 main.go:140: 本次图片保存完毕
const DefaultFormat = log.Ldate | log.Ltime | log.Lshortfile

var (
	i *log.Logger
	w *log.Logger
	e *log.Logger
)

// InitLog 初始化日志记录器
//
// 参数 path string 日志的路径
func InitLog(path string, format int) (*log.Logger, *log.Logger, *log.Logger) {
	// 日志输出文件
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		e.Fatalln(err)
	}
	// 自定义日志格式
	i = log.New(io.Writer(os.Stdout), "INFO: ", format)
	w = log.New(io.MultiWriter(file, os.Stdout), "WARN: ", format)
	e = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", format)
	return i, w, e
}
