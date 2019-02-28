package dofile

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 文件的魔术数字的信息
type fileMagicNum struct {
	// 需要读取文件头尾的字节数
	nHead int64
	nTail int64
	// 文件头尾必须等于该十六进制字符串，该字符串为大写
	headMust string
	tailMust string
}

var (
	WRITE_CREATE = os.O_CREATE
	WRITE_APPEND = os.O_CREATE | os.O_APPEND
	WRITE_TRUNC  = os.O_CREATE | os.O_TRUNC
)

// 读取文件
func Read(path string) (bs []byte, err error) {
	fi, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	bs, err = ioutil.ReadAll(fi)
	if err != nil {
		return
	}
	return
}

// 向文件写内容
func Write(bs []byte, path string, mode int, perm os.FileMode) (n int, err error) {
	fi, err := os.OpenFile(path, mode, perm)
	if err != nil {
		return
	}
	n, err = fi.Write(bs)
	return
}

// 判断路径是否存在
// 参考：https://blog.csdn.net/xielingyun/article/details/49992455
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	// 如果返回的错误为nil,说明文件或文件夹存在
	// 如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
	// 如果返回的错误为其它类型,则不确定是否在存在
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 判断是否为目录
// 参考：https://www.reddit.com/r/golang/comments/2fjwyk/isdir_in_go
func isDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// 返回文件MD5值
func Md5(path string) (md5Str string, err error) {
	isdir, err := isDir(path)
	if err != nil {
		return
	}
	if isdir {
		return "", fmt.Errorf("目标为目录")
	}

	out, err := os.Open(path)
	if err != nil {
		return
	}
	defer out.Close()

	md5hash := md5.New()
	if _, err = io.Copy(md5hash, out); err != nil {
		return
	}

	m5 := md5hash.Sum(nil)
	md5Str = fmt.Sprintf("%x", m5)
	return
}

// 文件完整性检测
func CheckIntegrity(path string) (integrity bool, err error) {
	// 判断目标是否存在，不存在则返回错误
	exist, err := PathExists(path)
	if err != nil {
		return
	}
	if !exist {
		return false, fmt.Errorf("目标不存在")
	}
	// 判断目标是否为目录，为目录则返回错误
	isdir, err := isDir(path)
	if err != nil {
		return
	}
	if isdir {
		return false, fmt.Errorf("目标为目录")
	}

	// 验证文件
	suffix := path[strings.LastIndex(path, "."):] // 文件后缀（包括点号"."）

	magicNum, err := getMagicNum(suffix)
	if err != nil {
		return
	}

	headBytes, tailBytes, err := ReadHeadTailBytes(path, magicNum.nHead, magicNum.nTail)
	head := hex.EncodeToString(headBytes)
	tail := hex.EncodeToString(tailBytes)

	if (magicNum.headMust == "" || strings.ToUpper(head) == magicNum.headMust) &&
		(magicNum.tailMust == "" || strings.ToUpper(tail) == magicNum.tailMust) {
		return true, nil
	}
	return false, nil
}

// 读取文件头尾的几个字节
func ReadHeadTailBytes(path string, n1 int64, n2 int64) (hbs []byte, tbs []byte, err error) {
	// 打开目标文件
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	// 读取文件头n字节内容
	startBytes := make([]byte, n1)
	_, err = file.Read(startBytes)
	if err != nil && err != io.EOF {
		return
	}

	// 读取文件尾n字节内容
	_, err = file.Seek(-1*n2, io.SeekEnd)
	if err != nil {
		return
	}
	endBytes := make([]byte, n2)
	_, err = file.Read(endBytes)
	if err != nil && err != io.EOF {
		return
	}

	return startBytes, endBytes, nil
}

// 获取魔术数字
// 文件头标志：https://www.cnblogs.com/WangAoBo/p/6366211.html
func getMagicNum(suffix string) (magic fileMagicNum, err error) {
	switch suffix {
	case ".jpg", ".jpeg":
		return fileMagicNum{4, 2, "FFD8FFE0", "FFD9"}, nil
	case ".png":
		return fileMagicNum{4, 4, "89504E47", ""}, nil
	case ".gif":
		return fileMagicNum{4, 4, "47494638", ""}, nil
	default:
		return magic, fmt.Errorf("未知的文件格式：%s", suffix)
	}
}
