//go:build linux

package dofile

import (
	"syscall"
)

// GetDriverInfo 获取磁盘的可用空间。依次为 可用空间、总空间、剩余空间（包含系统保留）
//
// 如果 path 为空""，则返回当前路径所在的分区的信息
//
// 参考 https://stackoverflow.com/a/60724929
func GetDriverInfo(path string) (free uint64, total uint64, avail uint64, err error) {
	if path == "" {
		path, err = syscall.Getwd()
		if err != nil {
			return
		}
	}

	var stat syscall.Statfs_t
	err = syscall.Statfs(wd, &stat)
	if err != nil {
		return
	}

	// 可用空间
	free = stat.Bavail * uint64(stat.Bsize)
	// 总空间
	total = stat.Blocks * uint64(stat.Bsize)
	// 剩余空间（包含系统保留）
	avail = stat.Bfree * uint64(stat.Bsize)

	return
}
