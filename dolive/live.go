// Package dolive 捕获视频流

package dolive

import (
	"bytes"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dohttp"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Live 用 New() 创建实例后使用
type Live struct {
	// 直播流的地址
	URL string

	// 下面几项会在创建该实例时，根据参数生成
	// 视频保存的目录
	FDir string
	// 视频文件名
	FName string
	// 视频格式。如 "mp4"
	FFormat string
	// 单个文件的最大字节数，为 0 表示无限制。建议 1GB: 1024*1024*1024
	FLSize int

	// 实际保存视频流时，记录数据
	// 已保存的字节数
	Total int
	// 当前写入的文件实例
	Cur *os.File
	// 直播流被保存到的文件的路径列表
	Paths []string
}

var client = dohttp.New(false, false)

// New 创建 Live 实例
//
// 参数 fLsize 单个文件的最大字节数，为 0 表示无限制。建议 SizeOneGB (1024*1024*1024)
func New(url string, path string, fLSize int) *Live {
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	name := strings.TrimRight(filepath.Base(path), ext)
	format := strings.Replace(ext, ".", "", 1)

	live := Live{
		URL:     url,
		FDir:    dir,
		FName:   name,
		FFormat: format,
		FLSize:  fLSize,
	}

	return &live
}

// Capture 捕获直播流到视频文件
//
// 参数 deal 处理捕获的视频文件的回调函数，参数为文件的路径
func (l *Live) Capture(headers map[string]string, deal func(path string) error) error {
	// 打开直播流
	req, err := http.NewRequest("GET", l.URL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Exec(req, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	defer l.closeCurFile()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("直播流响应：%s", resp.Status)
	}

	// 保存直播流
	// 缓存
	var buf = make([]byte, 1<<20)
	// 读写
	for {
		n, err := resp.Body.Read(buf)

		// 读取出错
		if n < 0 {
			return err
		}
		// 已读完
		if n == 0 {
			break
		}

		// 打开文件流
		if len(l.Paths) == 0 || (l.FLSize != 0 && l.Total > l.FLSize*len(l.Paths)) {
			fmt.Printf("写入第 %d 个文件\n", len(l.Paths)+1)
			err = l.createFileStream()
			if err != nil {
				return err
			}

			// 递归实现转到新文件存储视频
			// 重读直播流是因为新文件头需要视频信息，不然通过引入 ffmpeg 等直播捕获完毕后再切割，限制太大
			if len(l.Paths) >= 2 && l.FLSize != 0 {
				// 此处处理上一个下载完毕的视频
				err := deal(l.Paths[l.Total/(l.FLSize)-1])
				if err != nil {
					return err
				}

				return l.Capture(headers, deal)
			}
		}

		// 复制流到文件
		// 注意不能直接用 buf，必须 buf[:n]
		// 避免读取字节少于上次时，写入上次的多余数据
		_, err = io.Copy(l.Cur, bytes.NewReader(buf[:n]))
		if err != nil {
			return err
		}
		l.Total += n
	}

	// 此处处理最后一个视频（不分段时的整个视频，分段时的最后一个视频）
	return deal(l.Paths[len(l.Paths)-1])
}

// 首次或更换文件写入数据前，需要创建新的文件流
func (l *Live) createFileStream() error {
	// 首先，关闭上个文件流
	l.closeCurFile()

	// 生成路径
	path := filepath.Join(l.FDir, l.FName) + fmt.Sprintf(".%s", l.FFormat)
	if len(l.Paths) >= 1 {
		// 当发现需要分段保存视频，要先按约定重命名第一个文件，并将新路径保存到 Live 实例
		if len(l.Paths) == 1 {
			first := filepath.Join(l.FDir, l.FName) + fmt.Sprintf("_%02d.%s", 1, l.FFormat)
			if err := os.Rename(path, first); err != nil {
				return err
			}
			l.Paths[0] = first
		}

		path = filepath.Join(l.FDir, l.FName) + fmt.Sprintf("_%02d.%s", len(l.Paths)+1, l.FFormat)
	}

	// 打开文件流
	file, err := os.OpenFile(path, dofile.OTrunc, 0644)
	if err != nil {
		return err
	}

	// 保存到 Live 实例
	l.Cur = file
	l.Paths = append(l.Paths, path)
	return nil
}

// 关闭当前文件流
func (l *Live) closeCurFile() {
	if l.Cur != nil {
		l.Cur.Close()
	}
}
