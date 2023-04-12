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
//  var err error
//  Info, Warn, Error, err = InitLog("./logger_test.log", 0, false)
// }
// 使用：logger.Info.Printf("创建配置文件：%s\n", confPath)

package dolog

import (
	"io"
	"log"
	"os"
)

// DefaultFlag 默认的日志属性
//
// INFO: 2021/04/03 01:51:03 main.go:140: 本次图片保存完毕
const DefaultFlag = log.Ldate | log.Ltime | log.Lshortfile

// InitLog 初始化日志记录器
//
// path 日志文件的路径
//
// flag 日志的属性，是否包含日期、时间、源代码的文件名等。当传值 0 时，使用 DefaultFlag
//
// i2File 是否将 info 日志写入到日志文件。推荐 false
func InitLog(path string, flag int, i2File bool) (*log.Logger, *log.Logger, *log.Logger, error) {
	// 日志输出到文件
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, nil, err
	}

	// 默认
	if flag == 0 {
		flag = DefaultFlag
	}

	// 是否将 Info 日志输出到文件
	var iw io.Writer
	if i2File {
		iw = io.MultiWriter(file, os.Stdout)
	} else {
		iw = io.Writer(os.Stdout)
	}

	// 日志类型
	i := log.New(iw, "INFO: ", flag)
	w := log.New(io.MultiWriter(file, os.Stdout), "WARN: ", flag)
	e := log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", flag)

	return i, w, e, nil
}

// CkPanic 出错时，强制关闭程序
func CkPanic(err error) {
	if err != nil {
		panic(err)
	}
}
