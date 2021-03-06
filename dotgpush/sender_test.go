package dotgpush

import (
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dohttp"
	"testing"
)

var (
	tg     = NewTGBot("")
	chatID = ""
)

func init() {
	tg.SetProxy(dohttp.ProxySocks5)
}

func TestTGBot_SendMessage(t *testing.T) {
	msg, err := tg.SendMessage(chatID, "测试普通文本消息")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", *msg)
}

func TestSendMediaGroup(t *testing.T) {
	// 发送远程文件
	/*
		msg, err := tg.SendMediaGroup(chatID, []Media{
			{
				Type:    Photo,
				Media:   "https://cdn.v2ex.com/avatar/b600/b4a3/49950_large.png?m=1456725848",
				Caption: "头像1",
			},
			{
				Type:    Photo,
				Media:   "https://cdn.v2ex.com/gravatar/4598134cecabd98904511e065adca226?s=48",
				Caption: "头像2",
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("发送本地文件的结果：%+v\n", *msg)
	*/

	// 发送本地文件

	f1, err := dofile.Read("C:/Users/Do/Downloads/金-01.jpg")
	if err != nil {
		t.Fatal(err)
	}
	f2, err := dofile.Read("C:/Users/Do/Downloads/糧-01.jpg")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := tg.SendMediaGroup(chatID, []Media{
		{
			Type:    Photo,
			Media:   f1,
			Caption: "图1",
		},
		{
			Type:  Photo,
			Media: f2,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("发送本地文件的结果：%+v\n", *msg)
}
