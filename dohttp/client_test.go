package dohttp

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"
)

// 妖火cookie中的sid，如下格式
var yaohuoSid = ""
var client = New(60*time.Second, true, false)

func TestProxy(t *testing.T) {
	err := client.SetProxy("socks5://127.0.0.1:1081")
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
	type args struct {
		url       string
		path      string
		fieldname string
		otherForm map[string]string
		headers   map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "上传文件",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.PostFile(tt.args.url, tt.args.path, tt.args.fieldname,
				tt.args.otherForm, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statuscode(t *testing.T) {
	text, err := client.GetText("https://gg.doio.xyz/", nil)
	if err != nil {
		if errors.Is(err, ErrStatusCode) {
			log.Println("已匹配到状态码错误", err)
			return
		}
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
