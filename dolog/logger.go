// 初始化日志记录器
//
// var Info, Warn, Error = dolog.InitLog(0)
//
// 使用：logger.Info.Printf("创建配置文件：%s\n", confPath)
//
// 参考：https://www.jianshu.com/p/a9427a4e2ada

package dolog

import (
	"log"
	"os"
)

// DefaultFlag 默认的日志属性
//
// INFO: 2021/04/03 01:51:03 main.go:140: 本次图片保存完毕
const DefaultFlag = log.Ldate | log.Ltime | log.Lshortfile

// InitLog 初始化日志记录器
//
// flag 日志的属性，是否包含日期、时间、源代码的文件名等。当传值 0 时，使用 DefaultFlag
func InitLog(flag int) (*log.Logger, *log.Logger, *log.Logger) {
	// 默认日志标签
	if flag == 0 {
		flag = DefaultFlag
	}

	// 日志类型
	i := log.New(os.Stdout, "INFO: ", flag)
	w := log.New(os.Stdout, "WARN: ", flag)
	e := log.New(os.Stderr, "ERROR: ", flag)

	return i, w, e
}

// CkPanic 出错时，强制关闭程序
func CkPanic(err error) {
	if err != nil {
		panic(err)
	}
}
