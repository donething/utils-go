package dofile

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

const (
	// OCreate 若不存在则创建，只写模式
	OCreate = os.O_CREATE | os.O_WRONLY
	// OAppend 追加，不存在则先创建文件，只写
	OAppend = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	// OTrunc 覆盖或新建，只写模式
	OTrunc = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
)

// Read 读取文件
func Read(path string) ([]byte, error) {
	fi, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}

// Write 向文件写内容
func Write(bs []byte, path string, mode int, perm os.FileMode) (int, error) {
	fi, err := os.OpenFile(path, mode, perm)
	if err != nil {
		return 0, err
	}
	defer fi.Close()
	return fi.Write(bs)
}

// CopyFile 复制文件
//
// 返回值 n 表示复制的字节数
//
// 参考 https://opensource.com/article/18/6/copying-files-go
//
// 更详细的的实现，可参考 https://stackoverflow.com/a/21067803
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

// Exists 判断路径是否存在
//
// 参考 https://blog.csdn.net/xielingyun/article/details/49992455
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

// IsDir 判断是否为目录
//
// 参考 https://www.reddit.com/r/golang/comments/2fjwyk/isdir_in_go
func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// Md5 返回文件MD5值
func Md5(path string) (string, error) {
	isdir, err := IsDir(path)
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

// OpenAs 根据文件类型，选择合适的程序打开文件
//
// 来源 https://gist.github.com/hyg/9c4afcd91fe24316cbf0
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

// ShowInExplorer 在资源管理器中显示文件
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

// ValidFileName 合法化文件名
//
// 去除Windows下不能作为文件名的字符：<>:"/\|?*
func ValidFileName(src string, repl string) string {
	var reg = regexp.MustCompile(`[<>:"/\\|?*]`)
	return reg.ReplaceAllString(src, repl)
}
