package dotg

import (
	"bytes"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dohttp"
	"github.com/donething/utils-go/dovideo"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	tg     = NewTGBot(os.Getenv("MY_TG_TOKEN"))
	chatID = os.Getenv("MY_TG_CHAT_LIVE")
)

func init() {
	// tg.SetAddr("http://xxx:yyy")

	if strings.Contains(tg.addr, "telegram") {
		err := tg.SetProxy(dohttp.ProxySocks5)
		if err != nil {
			fmt.Printf("设置代理出错：%s\n", err)
		}
	}
}

func TestTGBot_SendMessage(t *testing.T) {
	txt := EscapeMk("测#试Markdown文本*消息*：") + "[搜索](https://www.google.com/)"
	for {
		msg, err := tg.SendMessage(chatID, txt)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v\n", *msg)
	}
}

func TestSendMediaGroupPic(t *testing.T) {
	path := "C:/Users/Do/Downloads/ipz275pl.jpg"
	f1, err := dofile.Read(path)
	if err != nil {
		t.Fatal(err)
	}

	medias := []*InputMedia{
		{
			Type:    TypePhoto,
			Caption: "图片：[测试](https://www.google.com/)",
			Name:    "p1 Reader",
			Media:   bytes.NewReader(f1),
		},
		{
			Type: TypePhoto,
			Name: "p2 路径",
			// Media: fmt.Sprintf("file:///%s", path),
			Media: bytes.NewReader(f1),
		},
	}

	msg, err := tg.SendMediaGroup(chatID, medias)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("发送本地文件的结果：%+v\n", *msg)
}

func TestTGBot_SendMediaGroupVideo(t *testing.T) {
	path := "D:/Tmp/VpsGo/Tmp/out 空格.mp4"
	bs, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	cbs, err := os.ReadFile("D:/Tmp/VpsGo/Tmp/output.jpg")
	if err != nil {
		t.Fatal(err)
	}

	w, h, err := dovideo.GetResolution(path)
	if err != nil {
		t.Fatal(err)
	}

	m := &InputMedia{
		Type:      TypeVideo,
		Media:     bytes.NewReader(bs),
		Thumbnail: bytes.NewReader(cbs),
		Caption:   "测试流媒体，可播放",
		ParseMode: "",

		Width:             w,
		Height:            h,
		SupportsStreaming: true,
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

func TestTGBot_SendVideo(t *testing.T) {
	msg, err := tg.SendVideo(os.Getenv("MY_TG_CHAT_LIVE"), "测试标题", "D:/Tmp/VpsGo/Tmp/jux-222-C.mp4",
		5*1024*1024, "", true)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", msg)
	time.Sleep(100 * time.Second)
}
