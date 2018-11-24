/**
文件操作：使用Get()函数得到一个文件
*/
package file

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DFile struct {
	Path string // Get()函数中传递的路径
	os.FileInfo
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
	ErrIsFile = errors.New("path to file")
	ErrIsDir  = errors.New("path to directory")
)

// 根据指定路径创建文件对象，如果路径不存在，则返回error
func Get(path string) (*DFile, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(path)
	return &DFile{Path: path, FileInfo: stat}, err
}

// 读取文件内容为字节
func (f *DFile) Read() (bytes []byte, err error) {
	if f.IsDir() {
		return nil, ErrIsDir
	}
	in, err := os.OpenFile(f.Path, os.O_RDONLY, PERM)
	if err != nil {
		return
	}
	defer in.Close()
	return ioutil.ReadAll(in)
}

// 将字节写入文件
func (f *DFile) Write(bytes []byte, append bool) (int, error) {
	if f.IsDir() {
		return 0, ErrIsDir
	}
	flag := 0
	if append {
		flag = os.O_APPEND
	} else {
		flag = os.O_TRUNC
	}

	out, err := os.OpenFile(f.Path, flag|os.O_CREATE, PERM)
	if err != nil {
		return 0, err
	}
	defer out.Close()
	return out.Write(bytes)
}

// 获取除去后缀后的文件名(filename)
func (f *DFile) BaseName() (name string) {
	name = f.Name()
	if dotIndex := strings.LastIndex(name, "."); dotIndex >= 0 {
		name = name[:dotIndex]
	}
	return
}

// 返回父文件夹，如果指定路径已为根路径("C:"、"/")，则仍然返回根路径

func (f *DFile) Parent() *DFile {
	p, _ := Get(f.ParentPath())
	return p
}

// 返回父文件夹的路径，如果指定路径已为根路径("C:"、"/")，则仍然返回根路径
func (f *DFile) ParentPath() (path string) {
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
func (f *DFile) Rename(newPath string) error {
	return os.Rename(f.Path, newPath)
}

// 删除文件
func (f *DFile) Del() error {
	return os.RemoveAll(f.Path)
}

// 列出目录
func (f *DFile) List(filter string) ([]DFile, error) {
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
		stat, _ := os.Stat(fpath)
		switch filter {
		case FILE: // 只获取文件
			if !tmp.IsDir() {
				filesList = append(filesList, DFile{Path: fpath, FileInfo: stat})
			}
		case DIR: // 只获取目录
			if tmp.IsDir() {
				filesList = append(filesList, DFile{Path: fpath, FileInfo: stat})
			}
		case ALL: // 获取文件和目录
			filesList = append(filesList, DFile{Path: fpath, FileInfo: stat})
		default: // 过滤条件错误
			return nil, fmt.Errorf("过滤条件错误")
		}
	}
	return filesList, nil
}

// 列出目录下文件的路径
func (f *DFile) ListPaths(filter string) ([]string, error) {
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

// 返回文件MD5值
func (f *DFile) Md5() (md5Str string, err error) {
	if f.IsDir() {
		return "", ErrIsDir
	}
	out, err := os.Open(f.Path)
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
