package dobdpan

import (
	"os"
	"testing"
)

func TestPrecreate(t *testing.T) {
	bs, err := os.ReadFile("C:/Users/Do/Downloads/金-01.jpg")
	if err != nil {
		t.Fatal(err)
	}
	yk := New(bs, "/4637251374683426/tttt.jpg", 1621090346, nil)
	resp, err := yk.precreate()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", resp)
}

func TestUploadYike(t *testing.T) {
	bs, err := os.ReadFile("C:/Users/Do/Downloads/logo.cdcfac33.png")
	if err != nil {
		t.Fatal(err)
	}

	req := GetYikeReq("aaaa", "bbb")

	f := New(bs, "/Test/logo.cdcfac33.png", 0, req)

	err = f.UploadFile()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("上传成功\n")
}

func TestUploadTerabox(t *testing.T) {
	bs, err := os.ReadFile("C:/Users/Do/Downloads/logo.cdcfac33.png")
	if err != nil {
		t.Fatal(err)
	}

	req1 := GetTeraboxReq("aaa")

	f := New(bs, "/Test/logo.cdcfac33.png", 0, req1)

	err = f.UploadFile()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("上传成功\n")
}

func TestDelAll(t *testing.T) {
	err := DelAll(nil)
	if err != nil {
		t.Fatal(err)
	}
}
