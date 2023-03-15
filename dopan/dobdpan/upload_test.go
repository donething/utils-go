package dobdpan

import (
	"testing"
)

func TestUploadYike(t *testing.T) {
	req := GetYikeReq("aaa", "ccc")

	f := New("/Test/logo.cdcfac33.png", 0, req)

	err := f.UploadFile("C:/Users/Do/Downloads/logo.cdcfac33.png")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("上传成功\n")
}

func TestUploadTerabox(t *testing.T) {
	req1 := GetTeraboxReq("aaa")

	f := New("/Test/logo.cdcfac33.png", 0, req1)

	err := f.UploadFile("C:/Users/Do/Downloads/logo.cdcfac33.png")
	if err != nil {
		t.Fatalf("出错：%s\n", err)
	}

	t.Logf("上传成功\n")
}

func TestDelAll(t *testing.T) {
	err := DelAll(nil)
	if err != nil {
		t.Fatal(err)
	}
}
