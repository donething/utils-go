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
	var data = entity.TplGeneral{
		Time:  entity.PushDataItem{dostr.FormatDate(time.Now(), dostr.TimeFormatDefault), ""},
		Title: entity.PushDataItem{"标题", ""},
		Msg:   entity.PushDataItem{"内容", ""},
	}
	var msgEntity = entity.PushMsg{
		Touser:     "okbRj1sBVC_dgzjUleRBZmDIcih4",
		TemplateID: "M6yJitQ_5XPYM6yAEI7xHE-LriN7CQuW_IMe0vKtJzM",
		URL:        "",
		Data:       data,
	}

	res, err := wx.PushTpl(msgEntity)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(res))
}
