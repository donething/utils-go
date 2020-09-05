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
var yaohuoSid = "31A*****0-0-0"
var client = New(60*time.Second, true, false)

func TestGetText(t *testing.T) {
	type args struct {
		url     string
		headers map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "正确获取内容",
			args: args{
				url:     "https://www.v2ex.com/api/members/show.json?id=1",
				headers: nil,
			},
			want:    "{\"username\": \"Livid\", \"website\": \"https://livid.v2ex.com/\", \"github\": \"\", \"psn\": \"\", \"avatar_normal\": \"https://cdn.v2ex.com/avatar/c4ca/4238/1_mini.png?m=1583753654\", \"bio\": \"Remember the bigger green\", \"url\": \"https://www.v2ex.com/u/Livid\", \"tagline\": \"Gravitated and spellbound\", \"twitter\": \"\", \"created\": 1272203146, \"status\": \"found\", \"avatar_large\": \"https://cdn.v2ex.com/avatar/c4ca/4238/1_mini.png?m=1583753654\", \"avatar_mini\": \"https://cdn.v2ex.com/avatar/c4ca/4238/1_mini.png?m=1583753654\", \"location\": \"\", \"btc\": \"\", \"id\": 1}",
			wantErr: false,
		},
		{
			name: "401错误",
			args: args{
				url:     "https://gg.doio.xyz/",
				headers: nil,
			},
			want:    "401",
			wantErr: true,
		},
		{
			name: "访问被ban网站",
			args: args{
				url:     "https://google.com",
				headers: nil,
			},
			want:    "google",
			wantErr: false,
		},
		{
			name: "带header请求",
			args: args{
				url:     "https://yaohuo.me/bbs-426174.html",
				headers: map[string]string{"Cookie": "sidyaohuo=" + yaohuoSid},
			},
			want:    "我也在找呜呜",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetText(tt.args.url, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("GetText() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	n, err := client.Download("https://code.jquery.com/jquery-1.12.4.min.js",
		"D:/Temp/jquery-1.12.4.min.js", false, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Printf("下载完成：%d\n", n)
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
				url:     "https://yaohuo.me/bbs/book_re.aspx",
				form:    "face=&sendmsg=0&content=%E7%9C%8B%E7%9C%8B%E6%8A%8A&action=add&id=426174&siteid=1000&lpage=1&classid=213&g=%E5%BF%AB%E9%80%9F%E5%9B%9E%E5%A4%8D&sid=" + yaohuoSid,
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
			got, err := client.PostFile(tt.args.url, tt.args.path, tt.args.fieldname, tt.args.otherForm, tt.args.headers)
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
