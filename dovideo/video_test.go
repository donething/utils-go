package dovideo

import "testing"

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
	err := Cut("D:/Tmp/VpsGo/video.mp4", 300*1024*1024, "D:/Tmp/VpsGo/Tmp")
	if err != nil {
		t.Fatal(err)
	}
}

func TestConvt(t *testing.T) {
	err := Convt("D:/Tmp/VpsGo/video.mp4", "D:/Tmp/VpsGo/video.mkv")
	if err != nil {
		t.Fatal(err)
	}
}
