package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type File struct {
	Path string
	This *os.File
}

const (
	// 新建文件的权限
	PERM = 0644
	// 当前系统的路径分隔符
	SEP = string(os.PathSeparator)
)

// 根据指定路径创建文件对象
func Get(path string) (file File, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	return File{Path: path, This: f}, nil
}

// 读取文件内容为字节
func (f File) Read() (bytes []byte, err error) {
	if f.IsDir() {
		err = fmt.Errorf("指定路径(%s)为目录，无法读取字节", f.Path)
		return
	}
	file, err := os.OpenFile(f.Path, os.O_RDONLY, PERM)
	if err != nil {
		return
	}
	return ioutil.ReadAll(file)
}

// 将字节写入文件
func (f File) Write(bytes []byte, append bool) (int, error) {
	if f.IsDir() {
		return 0, fmt.Errorf("指定路径(%s)为目录，无法写入字节", f.Path)
	}
	flag := 0
	if append {
		flag = os.O_APPEND
	} else {
		flag = os.O_TRUNC
	}

	file, err := os.OpenFile(f.Path, flag|os.O_CREATE, PERM)
	if err != nil {
		return 0, err
	}
	return file.Write(bytes)
}

// 文件是否存在
func (f File) Exist() bool {
	_, err := os.Stat(f.Path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// 文件是否为目录
func (f File) IsDir() bool {
	file, err := os.Stat(f.Path)
	if err != nil {
		return false
	}
	return file.IsDir()
}

// 获取文件名(filename.suffix)
func (f File) Name() (name string) {
	spIndex := strings.LastIndex(f.Path, SEP)
	name = f.Path[spIndex+1:]
	return
}

// 获取除去后缀后的文件名(filename)
func (f File) BaseName() (name string) {
	name = f.Name()
	if dotIndex := strings.LastIndex(name, "."); dotIndex >= 0 {
		name = name[:dotIndex]
	}
	return
}

// 返回父文件夹，如果指定路径已为根路径("C:"、"/")，则仍然返回根路径
func (f File) Parent() (path string) {
	path = f.Path
	if index := strings.LastIndex(path, SEP); index > 0 {
		path = path[0:index]
	}

	// 为类Unix系统专门处理("/home")
	if path == "" {
		path = SEP
	}
	return
}

// 重命名文件
func (f File) Rename(newPath string) error {
	if !f.Exist() {
		return fmt.Errorf("源路径(%s)不存在", f.Path)
	}
	return os.Rename(f.Path, newPath)
}

// 删除文件
func (f File) Del() error {
	return os.RemoveAll(f.Path)
}
