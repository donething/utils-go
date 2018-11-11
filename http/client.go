// 对http.Client的包装
// 使用dohttp.New()创建新的客户端
package dohttp

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

type DoClient struct {
	*http.Client
}

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
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略https证书错误
		}
		c.Transport = tr
	}

	return &DoClient{c}
}

// 执行请求
// 此函数执行完毕后不会关闭response.Body，且其它调用链中，后续需要读取Response的也不能关闭
// 需要在调用链中读取完Response后的函数（GetText()、GetFile()、Post()等）中关闭
func (client *DoClient) Request(req *http.Request, headers map[string]string) (res *http.Response, err error) {
	// 创建请求
	if err != nil {
		return
	}
	// 填充请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	// 执行请求
	res, err = client.Do(req)
	if err != nil {
		return
	}
	// 此时还不能关闭response，否则无法读取响应的内容
	// defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Fatalf("警告：请求（%s）的响应码不为OK：%s\n", req.URL, res.Status)
	}
	return
}

// 执行Get请求
func (client *DoClient) Get(url string, headers map[string]string) (res *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	res, err = client.Request(req, headers)
	// 因为没有后续操作，所以此处不需判断err是否为nil
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
	res, err := client.Request(req, headers)
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
