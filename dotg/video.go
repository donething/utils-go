package dotg

import (
	"fmt"
	"github.com/donething/utils-go/dovideo"
	"os"
	"path/filepath"
	"strings"
)

// GenTgMedia 生成上传视频到TG的 InputMedia 实例
func GenTgMedia(path string, title string) (media *InputMedia, dstPath string, thumbnail string, err error) {
	tag := "GenTgMedia"
	dstPath = path
	// 不是 mp4 格式的视频，才要转码为 mp4
	if strings.ToLower(filepath.Ext(path)) != ".mp4" {
		dstPath = strings.TrimSuffix(path, filepath.Ext(path)) + ".mp4"
		err = dovideo.Convt(path, dstPath)
		if err != nil {
			return nil, "", "", fmt.Errorf("[%s]转换视频编码出错：%w", tag, err)
		}

		// 删除原视频。本来可以放在末尾的，但是占用磁盘空间，所以在转码成功后删除
		err = os.Remove(path)
		if err != nil {
			return nil, "", "", fmt.Errorf("[%s]删除原视频出错：%w", tag, err)
		}
	}

	// 获取视频封面
	thumbnail = strings.TrimSuffix(dstPath, filepath.Ext(dstPath)) + ".jpg"
	err = dovideo.GetFrame(dstPath, thumbnail, "00:00:03", "320:320")
	if err != nil {
		return nil, "", "", fmt.Errorf("[%s]获取封面出错：%w", tag, err)
	}

	// 准备媒体数据
	vbs, err := os.Open(dstPath)
	if err != nil {
		return nil, "", "", fmt.Errorf("[%s]打开目标视频出错：%w", tag, err)
	}

	cbs, err := os.Open(thumbnail)
	if err != nil {
		return nil, "", "", fmt.Errorf("[%s]打开缩略图出错：%w", tag, err)
	}

	w, h, err := dovideo.GetResolution(dstPath)
	if err != nil {
		return nil, "", "", fmt.Errorf("[%s]获取视频分辨率出错：%w", tag, err)
	}
	media = &InputMedia{
		MediaData: &MediaData{
			Type:              TypeVideo,
			Caption:           title,
			Width:             w,
			Height:            h,
			SupportsStreaming: true,
		},
		Media:     vbs,
		Thumbnail: cbs,
	}

	return
}
