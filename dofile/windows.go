//go:build windows

package dofile

import (
	"golang.org/x/sys/windows"
	"os"
)

// GetDriverInfo 获取磁盘的可用空间。依次为 可用空间、总空间、剩余空间（包含系统保留）
//
// 如果 path 为空""，则返回当前路径所在的分区的信息
func GetDriverInfo(path string) (free uint64, total uint64, avail uint64, err error) {
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return
		}
	}

	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return
	}

	err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &avail)
	if err != nil {
		return
	}

	return
}
