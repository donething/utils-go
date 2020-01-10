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

func TestDoClient_Get(t *testing.T) {
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
		wantData []byte
		wantErr  bool
	}{
		{
			name:     "Test1",
			fields:   fields{},
			args:     args{url: "https://www.v2ex.com", headers: nil},
			wantData: []byte("DOCTYPE"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(10*time.Second, true, false)
			gotData, err := client.Get(tt.args.url, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoClient.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(string(gotData), "DOCTYPE") {
				t.Errorf("DoClient.Get() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

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
		{
			"Get Text",
			fields{},
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
			client := New(10*time.Second, true, false)
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
		{
			"Get File",
			fields{},
			args{"https://code.jquery.com/jquery-3.3.1.slim.min.js",
				nil,
				`E:/Temp/get_file.txt`,
			},
			69917,
			false,
		},
		{
			"Get File Not Found",
			fields{},
			args{"https://hu60.cn/1.txt",
				nil,
				`E:/Temp/temp.jpg`},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(10*time.Second, true, false)
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
	form := "type=1&name=肖申克&pass=1234567890&go=登录"
	type fields struct {
		Client *http.Client
	}
	type args struct {
		url     string
		form    string
		headers map[string]string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []byte
		wantErr  bool
	}{
		{
			name:   "post Form",
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

func TestDoClient_ReadTwiceResponse(t *testing.T) {
	client := New(30*time.Second, false, false)
	text, err := client.GetText("https://cililianbt.com/search/搜索/0/0/1.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "磁力链") {
		t.Error("出错")
	}
}

func TestDoClient_SetProxy(t *testing.T) {
	client := New(30*time.Second, false, false)

	client.SetProxy("http://127.0.0.1:1080")

	log.Printf("client信息：%+v\n", client.Transport)
	text, err := client.GetText("https://api.ipify.org", nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(text)
}

func TestDoClient_Post(t *testing.T) {
	client := New(10*time.Second, false, false)

	form := "reginvcode=cb1e6c4be12e1364&action=reginvcodeck"

	str, err := client.PostForm("http://fdfds1223fd.com", form, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
}

func TestDoClient_PostTest(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://fdfds1223fd.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	t.Log(res.Body)
}

func TestDoClient_PostFile(t *testing.T) {
	client := New(30*time.Second, false, false)
	otherForm := map[string]string{
		"file_id": "0",
	}
	data, err := client.PostFile(
		"https://sm.ms/api/upload?inajax=1&ssl=1",
		"D:/Users/Doneth/Pictures/BaiduShurufa_2018-12-23_19-40-45.png",
		"smfile",
		otherForm,
		nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestCheckNetworkConn(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			"Test Conn",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckNetworkConn(); got != tt.want {
				t.Errorf("CheckNetworkConn() = %v, want %v", got, tt.want)
			}
		})
	}
}
