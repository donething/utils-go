// 对http.Client的包装
// 使用dohttp.New()创建新的客户端
package dohttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"golang.org/x/net/proxy"
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

var ErrFileExist = errors.New("file is exist")

// dohttp.Client的包装
type DoClient struct {
	*http.Client
	Tr *http.Transport
}
type DoReq struct {
	*http.Request
}

// 创建新的DoClient
func New(timeout time.Duration, needCookieJar bool, checkSSL bool) *DoClient {
	// 根据参数，创建http.Client
	c := &http.Client{}
	tr := http.DefaultTransport.(*http.Transport)

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

	return &DoClient{c, tr}
}

// 设置代理
// proxyStr 代理，格式如下：
// http代理：http://127.0.0.1:1080
// https代理：https://127.0.0.1:1080
// socks5代理：socks5://127.0.0.1:1080
// 若为空字符串""，则清除之前设置的代理
func (client *DoClient) SetProxy(proxyStr string) error {
	if strings.Index(proxyStr, "http") == 0 {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			return err
		}
		client.Tr.Proxy = http.ProxyURL(proxyURL)
	} else if strings.Index(proxyStr, "socks5") == 0 {
		ipport := proxyStr[strings.Index(proxyStr, "//")+2:]
		dialer, err := proxy.SOCKS5("tcp", ipport, nil, proxy.Direct)
		if err != nil {
			return err
		}
		client.Tr.Dial = dialer.Dial
	}
	client.Transport = client.Tr
	return nil
}

// 设置cookie。若不设置，则使用默认jar
func (client *DoClient) SetJar(jar *http.CookieJar) {
	client.Jar = *jar
}

// 执行请求
// 此函数没有关闭response.Body
func (client *DoClient) Request(req *http.Request, headers map[string]string) (*http.Response, error) {
	//	// 填充请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	// 执行请求
	// 此时还不能关闭response，否则后续方法无法读取响应的内容
	return client.Do(req)
}

// 执行Get请求
func (client *DoClient) Get(url string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, 0, err
	}
	resp, err := client.Request(req, headers)
	if err != nil {
		return []byte{}, 0, err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	return bs, resp.StatusCode, err
}

// 读取文本类型
func (client *DoClient) GetText(url string, headers map[string]string) (string, int, error) {
	data, statusCode, err := client.Get(url, headers)
	return string(data), statusCode, err
}

// 下载文件到本地
// 并非一次读取、下载到内存，所以不用考虑网络上文件的大小
func (client *DoClient) DownFile(url string, savePath string, override bool, headers map[string]string) (int64, int, error) {
	exist, err := dofile.Exists(savePath)
	if exist && !override {
		return 0, 0, ErrFileExist
	}
	// 网络文件流
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, err
	}
	resp, err := client.Request(req, headers)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// 存储文件，需要放在网络连接后面，连接成功才创建新文件
	out, err := os.OpenFile(savePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, resp.StatusCode, err
	}
	defer out.Close()
	n, err := io.Copy(out, resp.Body)
	return n, resp.StatusCode, err
}

// Post请求
func (client *DoClient) post(req *http.Request, headers map[string]string) ([]byte, int, error) {
	resp, err := client.Request(req, headers)
	if err != nil {
		return []byte{}, 0, err
	}
	defer resp.Body.Close()
	n, err := ioutil.ReadAll(resp.Body)
	return n, resp.StatusCode, err
}

// Post表单
// form格式:a=1&b=2
func (client *DoClient) PostForm(url string, form string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form))
	if err != nil {
		return []byte{}, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return client.post(req, headers)
}

// post JSON字符串
func (client *DoClient) PostJSONString(url string, jsonStr string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonStr))
	if err != nil {
		return []byte{}, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.post(req, headers)
}

// POST map、struct等数据结构的json字符串
func (client *DoClient) PostJSONObj(url string, jsonObj interface{}, headers map[string]string) ([]byte, int, error) {
	jsonBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return []byte{}, 0, err
	}
	return client.PostJSONString(url, string(jsonBytes), headers)
}

// Post文件
// path：待上传文件的绝对路径
// fieldname：表单中文件对应的的键名
// otherForm：其它表单的键值
// headers：请求头
// https://www.golangnote.com/topic/124.html
func (client *DoClient) PostFile(url string, path string, fieldname string,
	otherForm map[string]string, headers map[string]string) ([]byte, int, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()

	// 添加其它表单值
	for k, v := range otherForm {
		_ = bodyWriter.WriteField(k, v)
	}

	// 添加文件表单值
	// use the bodyWriter to write the Part headers to the buffer
	_, err := bodyWriter.CreateFormFile(fieldname, filepath.Base(path))
	if err != nil {
		return []byte{}, 0, err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(path)
	if err != nil {
		return []byte{}, 0, err
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
		return []byte{}, 0, err
	}

	// Set headers for multipart, and Content Length
	fi, err := fh.Stat()
	if err != nil {
		return []byte{}, 0, err
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

// 检测响应码是否在200-399间
func CheckCode(code int) bool {
	return code >= 200 && code < 400
}
