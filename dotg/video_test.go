package dotg

import (
	"testing"
)

func TestGenTgMedia(t *testing.T) {
	media, dst, thumb, err := GenTgMedia("D:/Tmp/VpsGo/video.mp4", "测试视频标题")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v, %s, %s\n", *media, dst, thumb)
}
