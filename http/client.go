// 对http.Client的包装
// 使用dohttp.New()创建新的客户端
package dohttp

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

// http.Client的包装
type DoClient struct {
	*http.Client
}

// client参数
var tr = &http.Transport{}

// 创建新的DoClient
func New(timeout time.Duration, needCookieJar bool, checkSSL bool) *DoClient {
	// 根据参数，创建http.Client
	c := &http.Client{}
	// 超时时间
	c.Timeout = timeout
	// 需要管理Cookie
	if needCookieJar {
		cookieJar, _ := cookiejar.New(nil)
		c.Jar = cookieJar
	}
	// 不需要检查SSL
	if !checkSSL {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		c.Transport = tr
	}

	return &DoClient{c}
}

// 设置代理
// proxy格式如："http://127.0.0.1:1080"
// 若为空字符串""，则清除之前设置的代理
func (client *DoClient) SetProxy(proxy string) {
	if proxy == "" {
		tr.Proxy = nil
	} else {
		tr.Proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}
	}
	client.Transport = tr
}

// 执行请求
// 此函数执行完毕后不会关闭response.Body，且其它调用链中，后续需要读取Response的也不能关闭
// 需要在调用链中读取完Response后的函数（GetText()、GetFile()、Post()等）中关闭
func (client *DoClient) request(req *http.Request, headers map[string]string) (res *http.Response, err error) {
	// 填充请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	// 执行请求
	res, err = client.Do(req)
	// 此时还不能关闭response，否则无法读取响应的内容
	// defer res.Body.Close()

	// 因为没有后续操作，所以此处不需判断err==nil
	return
}

// 执行Get请求，返回http.Response的指针
// 该函数没有关闭response.Body，需读取响应后自行关闭
func (client *DoClient) Get(url string, headers map[string]string) (res *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	res, err = client.request(req, headers)
	// 因为没有后续操作，所以此处不需判断err==nil
	return
}

// 读取文本类型
func (client *DoClient) GetText(url string, headers map[string]string) (text string, err error) {
	res, err := client.Get(url, headers)
	if err != nil {
		return
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return string(data), nil
}

// 下载文件到本地
func (client *DoClient) GetFile(url string, headers map[string]string, savePath string) (size int64, err error) {
	res, err := client.Get(url, headers)
	if err != nil {
		return
	}
	defer res.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer out.Close()
	size, err = io.Copy(out, res.Body)
	return
}

// Post请求
// 次函数关闭了response：res.Body.Close()
func (client *DoClient) Post(req *http.Request, headers map[string]string) (data []byte, err error) {
	res, err := client.request(req, headers)
	if err != nil {
		return
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	return
}

// Post表单
func (client *DoClient) PostForm(url string, form url.Values, headers map[string]string) (data []byte, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return client.Post(req, headers)
}

// Post JSON字符串
func (client *DoClient) PostJSONString(url string, jsonStr string, headers map[string]string) (data []byte, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonStr))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Post(req, headers)
}

// POST map、struct等数据结构
func (client *DoClient) PostJSONObj(url string, jsonObj interface{}, headers map[string]string) (data []byte, err error) {
	jsonBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return
	}
	return client.PostJSONString(url, string(jsonBytes), headers)
}
