package dolog

import (
	"io"
	"log"
	"os"
)

const (
	// 可选的log文件名
	LOG_NAME = "run.log"
	// 可选的log记录格式
	LOG_FLAGS = log.LstdFlags | log.Lshortfile
)

// 设置将log保存到文件
// openFlags为文件读写模式
func Log2File(logName string, openFlags int, logFormat int) (err error) {
	// 打印log时显示时间戳
	log.SetFlags(logFormat)

	// 将日志输出到屏幕和日志文件
	lf, err := os.OpenFile(logName, openFlags, 0644)
	if err != nil {
		return
	}

	// 此句不能有，否则日志不能保存到文件中
	// defer lf.Close()
	// MultiWriter()的参数顺序也重要，如果使用"-H windowsgui"参数build，并且需要将日志保存到文件，
	// 则需要将日志文件的指针（lf）放到os.Stdout之前，否则log不会产生输出
	log.SetOutput(io.MultiWriter(lf, os.Stdout))
	return nil
}