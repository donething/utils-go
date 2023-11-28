package dovideo

import (
	"testing"
)

func TestGetVideoResolution(t *testing.T) {
	w, h, err := GetResolution("D:/Tmp/VpsGo/video.mp4")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(w, h)
}

func TestGetVideoDuration(t *testing.T) {
	seconds, err := GetDuration("D:/Tmp/VpsGo/video.mp4")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(seconds)
}

func TestGetFrame(t *testing.T) {
	err := GetFrame("D:/Tmp/VpsGo/video.mp4", "D:/Tmp/VpsGo/video.jpg",
		"00:00:03", "200:200")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCut(t *testing.T) {
	dst, err := Cut("D:/Tmp/VpsGo/Tmp/jux-222-C.mp4", 5*1024*1024, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", dst)
}

func TestConvt(t *testing.T) {
	err := Convt("D:/Tmp/VpsGo/video.mp4", "D:/Tmp/VpsGo/video.mkv")
	if err != nil {
		t.Fatal(err)
	}
}

func TestConcat(t *testing.T) {
	outputPath := "D:/Tmp/live/zuji_15722883_1684242152.mp4"
	err := Concat("D:/Tmp/live/zuji_15722883_1684242159", ".ts", outputPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("合并完成。输出的视频路径：%s\n", outputPath)
}

func TestCutMp4(t *testing.T) {
	dst, err := CutMp4("D:/Temp/VpsGo/TYOD-263.1080p.mp4", 2*1024*1024*1024, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", dst)
}
