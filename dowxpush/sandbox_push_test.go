package dowxpush

import (
	"github.com/donething/utils-go/dowxpush/entity"
	"log"
	"testing"
)

func TestNewSandbox(t *testing.T) {
	wx, err := NewSandbox("", "")
	if err != nil {
		t.Log(err)
	}
	log.Printf("%v", wx)

	// 第二个测试
	wx, err = NewSandbox("wxfbc379f966f4a234", "59be39dfcf4e95be0989748dad3e0ff6")
	if err != nil {
		t.Log(err)
	}
	log.Printf("%v", wx)
}

func TestWXSendbox_Push(t *testing.T) {
	wx, err := NewSandbox("wxfbc379f966f4a234", "59be39dfcf4e95be0989748dad3e0ff6")
	if err != nil {
		t.Fatal(err)
	}
	success, res, err := wx.PushTpl(entity.PushMsg{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(success, res)
}
