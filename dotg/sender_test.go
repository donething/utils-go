package dotg

import (
	"bytes"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dohttp"
	"os"
	"testing"
)

var (
	tg     = NewTGBot(os.Getenv("MY_TG_TOKEN"))
	chatID = os.Getenv("MY_TG_CHAT_LIVE")
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
	f1, err := dofile.Read("D:/Tmp/VpsGo/uploads/abc.jpg")
	if err != nil {
		t.Fatal(err)
	}

	medias := []*InputMedia{
		{
			MediaData: &MediaData{
				Type:    TypePhoto,
				Caption: "图片：[测试](https://www.google.com/)",
			},
			Media: bytes.NewReader(f1),
		},
		{
			MediaData: &MediaData{
				Type: TypePhoto,
			},
			Media: bytes.NewReader(f1),
		},
	}

	msg, err := tg.SendMediaGroup(chatID, medias)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("发送本地文件的结果：%+v\n", *msg)
}

func TestSendVideo(t *testing.T) {
	bs, err := os.ReadFile("D:/Tmp/VpsGo/output1.mp4")
	if err != nil {
		t.Fatal(err)
	}

	cbs, err := os.ReadFile("D:/Tmp/VpsGo/output.jpg")
	if err != nil {
		t.Fatal(err)
	}

	m := &InputMedia{
		MediaData: &MediaData{
			Type:      TypeVideo,
			Caption:   "测试流媒体，可播放",
			ParseMode: "",

			SupportsStreaming: true,
		},
		Media:     bytes.NewReader(bs),
		Thumbnail: bytes.NewReader(cbs),
	}

	msg, err := tg.SendMediaGroup(chatID, []*InputMedia{m})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", msg)
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
