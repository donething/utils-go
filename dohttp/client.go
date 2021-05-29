// Package dohttp 对 http.Client 的包装
// 使用dohttp.New()创建新的客户端
package dohttp

import (
	"bytes"
	"context"
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

// errors
var (
	// ErrFileExists 文件已存在
	ErrFileExists = errors.New("file already exist")
	// ErrStatusCode 网络响应代码不在 200-299 之间
	ErrStatusCode = errors.New("incorrect network status code")
	// ErrProxyString 网络代理错误
	ErrProxyString = errors.New("unknown proxy string")
)

type DoClient struct {
	*http.Client
}

// SetProxy 设置代理
//
// 参数 proxyStr string 代理地址，如 http://127.0.0.1:1081 socks5://127.0.0.1:1080 等
//
// 参数 auth *proxy.Auth 代理的用户名、密码，可空
func (c *DoClient) SetProxy(proxyStr string, auth *proxy.Auth) error {
	proxyStr = strings.ToLower(strings.TrimSpace(proxyStr))
	if strings.Index(proxyStr, "http") == 0 {
		proxyUrl, err := url.Parse(proxyStr)
		if err != nil {
			return err
		}
		c.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyUrl)
		return nil
	} else if strings.Index(proxyStr, "socks5") == 0 {
		// 只取后面的 127.0.0.1:1080
		addr := proxyStr[strings.LastIndex(proxyStr, "/")+1:]
		dialer, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
		if err != nil {
			return err
		}
		dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}
		c.Transport.(*http.Transport).DialContext = dialContext
		return nil
	}
	return ErrProxyString
}

// New 初始化
func New(timeout time.Duration, needCookieJar bool, checkSSL bool) DoClient {
	c := &http.Client{Transport: http.DefaultTransport}
	// 超时时间
	c.Timeout = timeout
	// 需要管理Cookie
	if needCookieJar {
		cookieJar, _ := cookiejar.New(nil)
		c.Jar = cookieJar
	}
	// 不需要检查SSL
	if !checkSSL {
		// 圆括号内为类型断言
		c.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return DoClient{c}
}

// Exec 执行请求
//
// 此函数没有关闭 response.Body
func (c *DoClient) Exec(req *http.Request, headers map[string]string) (*http.Response, error) {
	//	// 填充请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 执行请求
	// 此时还不能关闭response，否则后续方法无法读取响应的内容
	return c.Do(req)
}

// Get 执行Get请求
// 如果状态码不在 200-299 之间，会返回错误 ErrStatusCode
func (c *DoClient) Get(url string, headers map[string]string) ([]byte, error) {
	// 生成请求
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// 执行请求
	resp, err := c.Exec(req, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应内容
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 判断状态码，如果不在 200-299 之间，就返回读取的响应内容和错误 ErrStatusCode
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return bs, fmt.Errorf("%w %d", ErrStatusCode, resp.StatusCode)
	}
	return bs, nil
}

// GetText 读取文本类型
func (c *DoClient) GetText(url string, headers map[string]string) (string, error) {
	bs, err := c.Get(url, headers)
	return string(bs), err
}

// Download 下载文件到本地
//
// 如果本地存在此文件，且 override 参数为 false，会返回错误 ErrFileExists
//
// 如果状态码不在 200-299 之间，即使将文件保存到了本地，也会返回错误 ErrStatusCode
//
// 并非一次读取、下载到内存，所以不用考虑网络上文件的大小
func (c *DoClient) Download(url string, savePath string, override bool,
	headers map[string]string) (int64, error) {
	exist, err := dofile.Exists(savePath)
	if exist && !override {
		return 0, ErrFileExists
	}
	// 网络文件流
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.Exec(req, headers)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 存储文件，需要放在网络连接后面，连接成功才创建新文件
	out, err := os.OpenFile(savePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer out.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return 0, err
	}
	// 判断状态码，如果不在200-299间，依然返回ErrStatusCode error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return n, fmt.Errorf("%w %d", ErrStatusCode, resp.StatusCode)
	}
	return n, nil
}

// POST 请求
func (c *DoClient) post(req *http.Request, headers map[string]string) ([]byte, error) {
	resp, err := c.Exec(req, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 判断状态码，如果不在200-299 之间，就返回读取的响应内容和ErrStatusCode error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return bs, fmt.Errorf("%w %d", ErrStatusCode, resp.StatusCode)
	}
	return bs, nil
}

// PostForm POST 表单
// form 格式 a=1&b=2
func (c *DoClient) PostForm(url string, form string,
	headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.post(req, headers)
}

// PostJSONString POST JSON 字符串
func (c *DoClient) PostJSONString(url string, jsonStr string,
	headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.post(req, headers)
}

// PostJSONObj POST map、struct 等数据结构的 JSON 字符串
func (c *DoClient) PostJSONObj(url string, jsonObj interface{},
	headers map[string]string) ([]byte, error) {
	jsonBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return nil, err
	}
	return c.PostJSONString(url, string(jsonBytes), headers)
}

// PostFile POST 文件
//
// path：待上传文件的绝对路径
//
// fieldname：表单中文件对应的的键名
//
// otherForm：其它表单的键值
//
// headers：请求头
//
// 参考 https://www.golangnote.com/topic/124.html
func (c *DoClient) PostFile(url string, path string, fieldname string,
	otherForm map[string]string, headers map[string]string) ([]byte, error) {
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
		return nil, err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	// need to know the boundary to properly close the part myself.
	boundary := bodyWriter.Boundary()
	// close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	requestReader := io.MultiReader(bodyBuf, fh, closeBuf)
	req, err := http.NewRequest(http.MethodPost, url, requestReader)
	if err != nil {
		return nil, err
	}

	// Set headers for multipart, and Content Length
	fi, err := fh.Stat()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())

	return c.post(req, headers)
}

// CheckNetworkConn 检测网络是否可用
// 参考 https://stackoverflow.com/a/42227115
func CheckNetworkConn() bool {
	timeout := 10 * time.Second
	// 需要使用：ip:port 的格式
	// 此处使用百度搜索的IP和端口：baidu.com:80
	conn, err := net.DialTimeout("tcp", "baidu.com:80", timeout)
	if err != nil {
		return false
	}
	defer conn.Close() // 因为需要关闭连接，所有不能直接返回：return err!=nil
	return true
}

// CheckCode 检测响应码是否在 200-299 之间
func CheckCode(code int) bool {
	return code >= 200 && code <= 299
}

// LocalAddr 返回本机的局域网络地址
func LocalAddr() (string, error) {
	conn, err := net.Dial("ip:icmp", "google.com")
	if err != nil {
		return "", nil
	}
	return conn.LocalAddr().String(), nil
}
