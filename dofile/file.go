package dofile

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
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

const (
	WriteCreate = os.O_CREATE
	WriteAppend = os.O_CREATE | os.O_APPEND
	WriteTrunc  = os.O_CREATE | os.O_TRUNC
)

// 读取文件
func Read(path string) ([]byte, error) {
	fi, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}

// 向文件写内容
func Write(bs []byte, path string, mode int, perm os.FileMode) (int, error) {
	fi, err := os.OpenFile(path, mode, perm)
	if err != nil {
		return 0, err
	}
	defer fi.Close()
	return fi.Write(bs)
}

// 复制文件
// 返回值n表示复制的字节数
// 参考：https://opensource.com/article/18/6/copying-files-go
// 更详细的的实现，可参考：https://stackoverflow.com/a/21067803
func CopyFile(src string, dst string, override bool) (int64, error) {
	// 目标文件是否存在
	exist, err := Exists(dst)
	if err != nil {
		return 0, err
	}
	if exist && !override {
		return 0, fmt.Errorf("the dst file is already exists")
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	// non-regular files (e.g., directories, symlinks, devices, etc.
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// 判断路径是否存在
// 参考：https://blog.csdn.net/xielingyun/article/details/49992455
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	// 如果返回的错误为nil,说明文件或文件夹存在
	// 如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
	// 如果返回的错误为其它类型,则不确定是否在存在
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
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
func Md5(path string) (string, error) {
	isdir, err := isDir(path)
	if err != nil {
		return "", err
	}
	if isdir {
		return "", fmt.Errorf("目标为目录")
	}

	out, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	md5hash := md5.New()
	if _, err = io.Copy(md5hash, out); err != nil {
		return "", err
	}

	m5 := md5hash.Sum(nil)
	md5Str := fmt.Sprintf("%x", m5)
	return md5Str, nil
}

// 文件完整性检测
func CheckIntegrity(path string) (bool, error) {
	// 判断目标是否存在，不存在则返回错误
	exist, err := Exists(path)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, fmt.Errorf("目标不存在")
	}
	// 判断目标是否为目录，为目录则返回错误
	isdir, err := isDir(path)
	if err != nil {
		return false, err
	}
	if isdir {
		return false, fmt.Errorf("目标为目录")
	}

	// 开始验证文件
	suffix := path[strings.LastIndex(path, "."):] // 文件后缀（包括点号"."）

	magicNum, err := getMagicNum(suffix)
	if err != nil {
		return false, err
	}

	headBytes, tailBytes, err := ReadHeadTailBytes(path, magicNum.nHead, magicNum.nTail)
	head := hex.EncodeToString(headBytes)
	tail := hex.EncodeToString(tailBytes)

	// 文件类型的魔术数字可以为空字符串""，此时判为符合
	if (magicNum.headMust == "" || strings.ToUpper(head) == magicNum.headMust) &&
		(magicNum.tailMust == "" || strings.ToUpper(tail) == magicNum.tailMust) {
		return true, nil
	}
	return false, nil
}

// 读取文件头尾的几个字节
func ReadHeadTailBytes(path string, n1 int64, n2 int64) ([]byte, []byte, error) {
	// 打开目标文件
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// 读取文件头n字节内容
	startBytes := make([]byte, n1)
	_, err = file.Read(startBytes)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	// 读取文件尾n字节内容
	_, err = file.Seek(-1*n2, io.SeekEnd)
	if err != nil {
		return nil, nil, err
	}
	endBytes := make([]byte, n2)
	_, err = file.Read(endBytes)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	return startBytes, endBytes, nil
}

// 获取魔术数字
// 文件头标志：https://www.cnblogs.com/WangAoBo/p/6366211.html
func getMagicNum(suffix string) (fileMagicNum, error) {
	switch strings.ToLower(suffix) {
	case ".jpg", ".jpeg":
		return fileMagicNum{4, 2, "FFD8FFE0", "FFD9"}, nil
	case ".png":
		return fileMagicNum{4, 4, "89504E47", "AE426082"}, nil
	case ".gif":
		return fileMagicNum{4, 4, "47494638", ""}, nil
	default:
		return fileMagicNum{}, fmt.Errorf("未知的文件格式：%s", suffix)
	}
}

// 根据文件类型，选择合适的程序打开文件
// 来源：https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenAs(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	err := cmd.Start()
	return err
}

// 在资源管理器中显示文件
func ShowInExplorer(path string) error {
	var cmd *exec.Cmd
	// 使用explorer显示文件
	switch runtime.GOOS {
	case "windows":
		// explorer /select,path...，只能使用反斜杠"\"，不能使用斜杠"/"
		path = strings.Replace(path, `/`, `\`, -1)
		cmd = exec.Command("explorer", "/select,", path)
	default:
		return fmt.Errorf("还未适配平台：%s", runtime.GOOS)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	// 不知何故，即使执行成功，依然会返回`exit status 1`的error。所以此处手动排除该错误
	if err != nil && err.Error() != `exit status 1` {
		return fmt.Errorf("err: %s,stderr: %s,strout: %s", err.Error(), stderr.String(), out.String())
	}
	return nil
}
