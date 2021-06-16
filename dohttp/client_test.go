package dohttp

import (
	"log"
	"reflect"
	"strings"
	"testing"
	"time"
)

// 妖火cookie中的sid，如下格式
var yaohuoSid = ""
var client = New(10*time.Second, true, false)

func TestProxy(t *testing.T) {
	err := client.SetProxy("socks5://127.0.0.1:1080")
	if err != nil {
		t.Fatal(err)
	}
	text, err := client.GetText("https://google.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(text)
}

func TestPostForm(t *testing.T) {
	type args struct {
		url     string
		form    string
		headers map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "带header的post form",
			args: args{
				url: "https://yaohuo.me/bbs/book_re.aspx",
				form: "face=&sendmsg=0&content=%E7%9C%8B%E7%9C%8B%E6%8A%8A&action=add&id=426174" +
					"&siteid=1000&lpage=1&classid=213&g=%E5%BF%AB%E9%80%9F%E5%9B%9E%E5%A4%8D&sid=" + yaohuoSid,
				headers: map[string]string{"Cookie": "sidyaohuo=" + yaohuoSid},
			},
			want:    []byte("回复成功"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.PostForm(tt.args.url, tt.args.form, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(string(got), string(tt.want)) {
				t.Errorf("PostForm() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostJSONObj(t *testing.T) {
	type args struct {
		url     string
		jsonObj interface{}
		headers map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "",
			args: args{
				url:     "",
				jsonObj: "",
				headers: nil,
			},
			want:    []byte(""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.PostJSONObj(tt.args.url, tt.args.jsonObj, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostJSONObj() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostJSONObj() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostFile(t *testing.T) {
	// 发送本地文件
	client.PostFiles("http://127.0.0.1:10000/test/upload",
		map[string]interface{}{"bili.txt": "D:/哔哩哔哩收藏的视频.txt", "guan.txt": "D:/官解.txt",
			"pp.mp3": "D:/Timon & Pumbaa_ Stand By Me (Song) (320  kbps).mp3",
		},
		map[string]string{"p1": "p111", "p2": "p2222"},
		nil,
	)

	// 发送远程文件
	bs1, _ := client.Get("https://cdn.v2ex.com/gravatar/4598134cecabd98904511e065adca226?s=48&d=retro",
		nil)
	bs2, _ := client.Get("https://cdn.v2ex.com/gravatar/ff349a0ec97ea9e36b5aab456a38dbf2?s=48&d=retro",
		nil)
	client.PostFiles("http://127.0.0.1:10000/test/upload", map[string]interface{}{
		"11": bs1,
		"22": bs2,
	}, map[string]string{
		"test1": "test111",
		"test2": "test222",
	}, nil)
}

func Test_statuscode(t *testing.T) {
	text, err := client.GetText("https://abc.xyz/", nil)
	if err != nil {
		t.Fatal(err.Error(), text)
	}
	log.Println("文本", text)
}

func TestDoClient_Download(t *testing.T) {
	type args struct {
		url      string
		savePath string
		override bool
		headers  map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "正确下载文件",
			args: args{
				url:      "https://code.jquery.com/jquery-1.12.4.min.js",
				savePath: "D:/Temp/jquery-1.12.4.min.js",
				override: false,
				headers:  nil,
			},
			want:    97163,
			wantErr: false,
		},
		{
			name: "正确下载文件",
			args: args{
				url:      "https://code.jquery.com/t.txt",
				savePath: "D:/Temp/t.txt",
				override: false,
				headers:  nil,
			},
			want:    162,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client
			got, err := c.Download(tt.args.url, tt.args.savePath, tt.args.override, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Download() got = %v, want %v", got, tt.want)
			}
		})
	}
}
