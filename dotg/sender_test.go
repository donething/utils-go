package dotg

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dohttp"
	"testing"
)

var (
	tg     = NewTGBot("xxx")
	chatID = "yyy"
)

func init() {
	err := tg.SetProxy(dohttp.ProxySocks5)
	if err != nil {
		fmt.Printf("设置代理出错：%s\n", err)
	}
}

func TestTGBot_SendMessage(t *testing.T) {
	txt := EscapeMk("测#试Markdown文本*消息*：") + "[搜索](https://www.google.com/)"
	msg, err := tg.SendMessage(chatID, txt)
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

	f1, err := dofile.Read("D:/Tmp/VpsGo/uploads/abc.jpg")
	if err != nil {
		t.Fatal(err)
	}

	medias := []Media{
		{
			Type:    TypePhoto,
			Media:   f1,
			Caption: "图片：[搜索](https://www.google.com/)",
		},
		{
			Type:  TypePhoto,
			Media: f1,
		},
	}

	msg, err := tg.SendMediaGroup(chatID, medias)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("发送本地文件的结果：%+v\n", *msg)
}

func TestReplaceMk(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "测试",
			args: args{text: "测试_不错*继续[中括号](小括号)大于>井号#感叹号!结尾。"},
			want: "测试\\_不错\\*继续\\[中括号\\]\\(小括号\\)大于\\>井号\\#感叹号\\!结尾。",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeMk(tt.args.text); got != tt.want {
				t.Errorf("EscapeMk() = %v, want %v", got, tt.want)
			}
		})
	}
}
