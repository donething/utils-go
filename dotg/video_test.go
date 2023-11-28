package dotg

import (
	"testing"
)

func TestGenTgMedia(t *testing.T) {
	media, err := GenVideoMedia("D:/Downloads/PT/无码破解中字/爱田奈奈/JUC-620-UC.mp4", "测试视频标题")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", *media)
}
