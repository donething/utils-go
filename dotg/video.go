package dotg

import (
	"fmt"
	"github.com/donething/utils-go/dovideo"
	"os"
	"path/filepath"
	"strings"
)

// GenVideoMedia 生成上传视频到TG的 InputMedia 实例
//
// 仅适合 TG Local Server 模式，因为媒体数据是通过"file://"协议传的，而不是建立pipe管道读写数据，
// 这样避免 write: connection reset by peer，也更稳定
//
// 会转码、生成封面，不会切割长视频
//
// 注意：转码成功时会删除原视频
//
// 会返回新视频、封面的路径，以便上传后删除
func GenVideoMedia(path string, title string) (*InputMedia, error) {
	tag := "GenVideoMedia"

	// 获取视频封面
	thumbPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".jpg"
	err := dovideo.GetFrame(path, thumbPath, "00:00:03", "")
	if err != nil {
		return nil, fmt.Errorf("[%s]获取封面出错：%w", tag, err)
	}

	// 准备媒体数据
	cbs, err := os.Open(thumbPath)
	if err != nil {
		return nil, fmt.Errorf("[%s]打开缩略图出错：%w", tag, err)
	}

	w, h, err := dovideo.GetResolution(path)
	if err != nil {
		return nil, fmt.Errorf("[%s]获取视频分辨率出错：%w", tag, err)
	}

	media := &InputMedia{
		Type:              TypeVideo,
		Media:             fmt.Sprintf("file://%s", path),
		Thumbnail:         cbs,
		Caption:           title,
		Width:             w,
		Height:            h,
		SupportsStreaming: true,
	}

	return media, nil
}
