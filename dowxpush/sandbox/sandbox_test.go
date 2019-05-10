package sandbox

import (
	"github.com/donething/utils-go/dostr"
	"github.com/donething/utils-go/dowxpush/sandbox/entity"
	"log"
	"testing"
	"time"
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

	var msgEntity = entity.PushMsg{
		Touser:     "okbRj1sBVC_dgzjUleRBZmDIcih4",
		TemplateID: "M6yJitQ_5XPYM6yAEI7xHE-LriN7CQuW_IMe0vKtJzM",
		URL:        "",
		Data:       nil,
	}
	msgEntity.AddData("time", dostr.FormatDate(time.Now(), dostr.TimeFormatDefault), "").
		AddData("title", "标题", "").
		AddData("msg", "内容", "")

	resp, err := wx.PushTpl(msgEntity)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", resp)
}
