package dobdpan

import (
	"os"
	"testing"
)

func TestUploadTerabox(t *testing.T) {
	req1 := GetTeraboxReq("aaa")

	// 路径
	pf, err := NewPath("C:/Users/Do/Downloads/logo.cdcfac33.png",
		"/Test/logo.cdcfac33.png", req1)
	if err != nil {
		t.Fatal(err)
	}

	err = pf.UploadFile()
	if err != nil {
		t.Fatalf("出错：%s\n", err)
	}

	t.Logf("路径文件上传成功\n")

	// 二进制
	bs, err := os.ReadFile("C:/Users/Do/Downloads/logo.cdcfac33.png")
	if err != nil {
		t.Fatal(err)
	}
	bf := NewBytes(bs, "/Test/test.png", req1, 0)

	err = bf.UploadFile()
	if err != nil {
		t.Fatalf("出错：%s\n", err)
	}

	t.Logf("二进制文件上传成功\n")
}

func TestDelAll(t *testing.T) {
	err := DelAll(nil)
	if err != nil {
		t.Fatal(err)
	}
}
