package dohttp

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestDoClient_GetText(t *testing.T) {
	type fields struct {
		Client *http.Client
	}
	type args struct {
		url     string
		headers map[string]string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantText string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			"Get Text",
			fields{&http.Client{}},
			args{"http://home.baidu.com/home/index/contact_us", nil},
			"联系我们",
			false,
		},
		{
			"TimeOut",
			fields{},
			args{"https://google.com", nil},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(30*time.Second, true, false)
			gotText, err := client.GetText(tt.args.url, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoClient.GetText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(gotText, tt.wantText) {
				t.Errorf("DoClient.GetText() = %v, want %v", gotText, tt.wantText)
			}
		})
	}
}

func TestDoClient_GetFile(t *testing.T) {
	type fields struct {
		Client *http.Client
	}
	type args struct {
		url      string
		headers  map[string]string
		savePath string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantSize int64
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			"Get Text",
			fields{&http.Client{}},
			args{"https://code.jquery.com/jquery-3.3.1.slim.min.js",
				nil,
				"/home/doneth/Temp/get_file.txt",
			},
			69917,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &DoClient{
				Client: tt.fields.Client,
			}
			gotSize, err := client.GetFile(tt.args.url, tt.args.headers, tt.args.savePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoClient.GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSize != tt.wantSize {
				t.Errorf("DoClient.GetFile() = %v, want %v", gotSize, tt.wantSize)
			}
		})
	}
}

func TestDoClient_PostForm(t *testing.T) {
	form := url.Values{}
	form.Add("type", "1")
	form.Add("name", "肖申克")
	form.Add("pass", "1234567890")
	form.Add("go", "登录")
	type fields struct {
		Client *http.Client
	}
	type args struct {
		url     string
		form    url.Values
		headers map[string]string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []byte
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name:   "Post Form",
			fields: fields{},
			args: args{"https://hu60.net/q.php/user.login.html?u=index.index.html",
				form,
				nil,
			},
			wantData: []byte("抱歉，用户未激活"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(30*time.Second, true, false)
			gotData, err := client.PostForm(tt.args.url, tt.args.form, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoClient.PostForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(string(gotData), string(tt.wantData)) {
				t.Errorf("DoClient.PostForm() = %v, want %v", string(gotData), string(tt.wantData))
			}
		})
	}
}

func TestDoClient_Form_JSON(t *testing.T) {
	form := url.Values{}
	form.Add("k1", "v1")
	form.Add("k2", "v2")
	form.Add("k3", "v3")
	t.Log("form:", form.Encode())

	jsonMap := map[string]string{}
	jsonMap["k1"] = "k1"
	jsonMap["k2"] = "k2"
	jsonMap["k2"] = "k3"
	bs, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("jsonMap:", string(bs))
}
