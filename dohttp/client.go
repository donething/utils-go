// 对http.Client的包装
// 使用dohttp.New()创建新的客户端
package dohttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 状态码不在200-399内
var ErrStatusCode = errors.New("错误的状态码")

// dohttp.Client的包装
type DoClient struct {
	*http.Client
}
type DoReq struct {
	*http.Request
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

// 设置cookie，不设置，则使用默认jar
func (client *DoClient) SetJar(jar *http.CookieJar) {
	client.Jar = *jar
}

// 执行请求
// 此函数没有关闭response.Body
func (client *DoClient) Request(req *http.Request, headers map[string]string) (resp *http.Response, err error) {
	//	// 填充请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	// 执行请求
	// 此时还不能关闭response，否则后续方法无法读取响应的内容
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		err = ErrStatusCode
	}
	return
}

// 执行Get请求
func (client *DoClient) Get(url string, headers map[string]string) (data []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := client.Request(req, headers)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return
}

// 读取文本类型
func (client *DoClient) GetText(url string, headers map[string]string) (string, error) {
	data, err := client.Get(url, headers)
	return string(data), err
}

// 下载文件到本地
func (client *DoClient) GetFile(url string, headers map[string]string, savePath string) (size int64, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := client.Request(req, headers)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer out.Close()
	size, err = io.Copy(out, resp.Body)
	return
}

// Post请求
func (client *DoClient) post(req *http.Request, headers map[string]string) (data []byte, err error) {
	res, err := client.Request(req, headers)
	if err != nil {
		return
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	return
}

// Post表单
// form格式:a=1&b=2
func (client *DoClient) PostForm(url string, form string, headers map[string]string) (data []byte, err error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return client.post(req, headers)
}

// post JSON字符串
func (client *DoClient) PostJSONString(url string, jsonStr string, headers map[string]string) (data []byte, err error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonStr))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return client.post(req, headers)
}

// POST map、struct等数据结构的json字符串
func (client *DoClient) PostJSONObj(url string, jsonObj interface{}, headers map[string]string) (data []byte, err error) {
	jsonBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return
	}
	return client.PostJSONString(url, string(jsonBytes), headers)
}

// Post文件
// filename：待上传文件的路径
// fileFormField：表单中表示上传文件的键
// otherForm：其它表单
// https://www.golangnote.com/topic/124.html
func (client *DoClient) PostFile(url string, filename string, fileFormField string, otherForm map[string]string, headers map[string]string) (data []byte, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()

	// 添加其它表单值
	for k, v := range otherForm {
		_ = bodyWriter.WriteField(k, v)
	}

	// 添加文件表单值
	// use the bodyWriter to write the Part headers to the buffer
	_, err = bodyWriter.CreateFormFile(fileFormField, filepath.Base(filename))
	if err != nil {
		return
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fh.Close()

	// need to know the boundary to properly close the part myself.
	boundary := bodyWriter.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	requestReader := io.MultiReader(bodyBuf, fh, closeBuf)
	req, err := http.NewRequest(http.MethodPost, url, requestReader)
	if err != nil {
		return
	}

	// Set headers for multipart, and Content Length
	fi, err := fh.Stat()
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())

	return client.post(req, headers)
}

// 检测网络是否可用
// 参考：https://stackoverflow.com/a/42227115
func CheckNetworkConn() bool {
	timeout := 3 * time.Second
	// 需要使用：ip:port 的格式
	// 此处使用百度搜索的IP和端口：123.125.115.110:80
	conn, err := net.DialTimeout("tcp", "123.125.115.110:80", timeout)
	if err != nil {
		return false
	}
	defer conn.Close() // 因为需要关闭连接，所有不能直接返回：return err!=nil
	return true
}
