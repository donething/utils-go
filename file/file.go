/**
文件操作：使用Get()函数得到一个文件
*/
package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DFile struct {
	Path string // Get()函数中传递的路径
}

const (
	// 新建文件的权限
	PERM = 0644
	// 当前系统的路径分隔符
	SEP = string(os.PathSeparator)

	// 文件和目录的tag
	FILE = "FILE"
	DIR  = "DIR"
	ALL  = FILE + " " + DIR
)

var (
	ErrFileNotExists = errors.New("no such file or directory")
	ErrIsFile        = errors.New("path to file")
	ErrIsDir         = errors.New("path to directory")
)

// 根据指定路径创建文件对象，如果路径不存在，则返回error
func Get(path string) (DFile, error) {
	_, err := os.Stat(path)
	if err != nil {
		// 文件不存在错误
		if strings.Contains(err.Error(), "no such file or directory") {
			return DFile{}, ErrFileNotExists
		}
		// 其它错误
		return DFile{}, err
	}
	if os.IsNotExist(err) {
		return DFile{}, ErrFileNotExists
	}
	return DFile{Path: path}, nil
}

// 读取文件内容为字节
func (f DFile) Read() (bytes []byte, err error) {
	if f.IsDir() {
		return nil, ErrIsDir
	}
	file, err := os.OpenFile(f.Path, os.O_RDONLY, PERM)
	if err != nil {
		return
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

// 将字节写入文件
func (f DFile) Write(bytes []byte, append bool) (int, error) {
	if f.IsDir() {
		return 0, ErrIsDir
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
	defer file.Close()
	return file.Write(bytes)
}

// 文件是否为目录
func (f DFile) IsDir() bool {
	fileinfo, err := os.Stat(f.Path)
	if err != nil {
		return false
	}
	return fileinfo.IsDir()
}

// 获取文件名(filename.suffix)
func (f DFile) Name() (name string) {
	spIndex := strings.LastIndex(f.Path, SEP)
	name = f.Path[spIndex+1:]
	return
}

// 获取除去后缀后的文件名(filename)
func (f DFile) BaseName() (name string) {
	name = f.Name()
	if dotIndex := strings.LastIndex(name, "."); dotIndex >= 0 {
		name = name[:dotIndex]
	}
	return
}

// 返回父文件夹，如果指定路径已为根路径("C:"、"/")，则仍然返回根路径
func (f DFile) Parent() (path string) {
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
func (f DFile) Rename(newPath string) error {
	return os.Rename(f.Path, newPath)
}

// 删除文件
func (f DFile) Del() error {
	return os.RemoveAll(f.Path)
}

// 列出目录
func (f DFile) List(filter string) ([]DFile, error) {
	// 指定的对象为文件，无法列出目录
	if !f.IsDir() {
		return nil, ErrIsFile
	}

	// 将返回的文件列表
	var filesList = make([]DFile, 0, 0)

	// 获取目录下的文件
	files, err := ioutil.ReadDir(f.Path)
	if err != nil {
		return nil, err
	}

	// 根据过滤条件选择追加文件
	for _, tmp := range files {
		fpath := filepath.Clean(f.Path + SEP + tmp.Name()) // 文件路径
		switch filter {
		case FILE: // 只获取文件
			if !tmp.IsDir() {
				filesList = append(filesList, DFile{Path: fpath})
			}
		case DIR: // 只获取目录
			if tmp.IsDir() {
				filesList = append(filesList, DFile{Path: fpath})
			}
		case ALL: // 获取文件和目录
			filesList = append(filesList, DFile{Path: fpath})
		default: // 过滤条件错误
			return nil, fmt.Errorf("过滤条件错误")
		}
	}
	return filesList, nil
}

// 列出目录下文件的路径
func (f DFile) ListPaths(filter string) ([]string, error) {
	dfiles, err := f.List(filter)
	if err != nil {
		return nil, err
	}
	// 讲DFile切片转为路径字符串切片
	var paths = make([]string, 0, 0)
	for _, f := range dfiles {
		paths = append(paths, f.Path)
	}
	return paths, nil
}
