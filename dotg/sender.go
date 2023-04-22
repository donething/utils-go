// Package dotg 绑定机器人的 token 后即可发送消息
// tg := dotg.NewTGBot(token)
// tg.SendMessage(chatID, text)
package dotg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"github.com/donething/utils-go/dovideo"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TGBot TG 推送发送消息的机器人实例
type TGBot struct {
	// TG Token
	token string

	// API 地址。是本地服务地址，还是 TG 直连
	addr string
}

const (
	urlSendMsg        = "%s/%s/sendMessage"
	urlSendMediaGroup = "%s/%s/sendMediaGroup"

	// FileSizeThreshold TG 上传视频有2GB的限制，此处为了容错设置低一点
	FileSizeThreshold = 1800 * 1024 * 1024
)

var (
	client = dohttp.New(false, false)

	ErrResend = fmt.Errorf("发送过快，要求重发")
)

// NewTGBot 创建新的 Telegram 推送机器人
//
// token 新机器人的 token
func NewTGBot(token string) *TGBot {
	return &TGBot{
		// TG Token
		token: "bot" + token,

		// 默认 addr
		addr: "https://api.telegram.org",
	}
}

// SetProxy 设置网络代理
//
// 格式参考 dohttp.ProxySocks5、dohttp.ProxyHttp
func (bot *TGBot) SetProxy(proxyStr string) error {
	return client.SetProxy(proxyStr)
}

// SetAddr 设置域名。开启 telegram-bot-api 本地服务时，可用本地服务地址
//
// addr 如 "http://127.0.0.1:12345"
func (bot *TGBot) SetAddr(addr string) {
	bot.addr = addr
}

// Send 实际执行发送请求
//
// 出错会通过 chan 回传 error。如果 error 是
func (bot *TGBot) Send(url string, reader io.Reader, contentType string, chResult chan SendResult) {
	tag := "Send"
	// 发送
	sendUrl := fmt.Sprintf(url, bot.addr, bot.token)
	resp, err := client.Post(sendUrl, contentType, reader)
	if err != nil {
		chResult <- SendResult{
			Message: nil,
			Error:   fmt.Errorf("[%s]执行请求出错：%w", tag, err),
		}
		return
	}
	defer resp.Body.Close()

	if c, ok := reader.(io.ReadCloser); ok {
		c.Close()
	}

	// 读取响应
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		chResult <- SendResult{
			Message: nil,
			Error:   fmt.Errorf("[%s]读取响应内容出错：%w", tag, err),
		}
		return
	}

	// 解析响应
	var msg Message
	err = json.Unmarshal(bs, &msg)
	if err != nil {
		chResult <- SendResult{
			Message: nil,
			Error:   fmt.Errorf("[%s]解析响应内容出错：%w", tag, err),
		}
		return
	}

	// 速率限制
	if msg.ErrorCode == 429 {
		secStr := msg.Description[strings.LastIndex(msg.Description, " ")+1:]
		sec, err := strconv.Atoi(secStr)
		if err != nil {
			chResult <- SendResult{
				Message: nil,
				Error:   fmt.Errorf("[%s]解析速率限制的等待时长时出错：%w", tag, err),
			}
			return
		}

		fmt.Printf("[%s]由于速率限制，等待 %d 秒后重新发送\n", tag, sec)
		time.Sleep(time.Duration(sec+1) * time.Second)
		// 重发
		chResult <- SendResult{
			Message: nil,
			Error:   ErrResend,
		}
		return
	}

	// 发送失败
	if !msg.Ok {
		chResult <- SendResult{
			Message: nil,
			Error:   fmt.Errorf("[%s]发送失败：%d %s", tag, msg.ErrorCode, msg.Description),
		}
		return
	}

	// 成功
	chResult <- SendResult{
		Message: &msg,
		Error:   nil,
	}
}

// SendMessage 发送Markdown文本消息
//
// 注意使用 EscapeMk 来转义文本中的非法字符（即属于 Markdown 字符，而不想当做 Markdown 字符渲染）
func (bot *TGBot) SendMessage(chatID string, text string) (*Message, error) {
	tag := "SendMessage"
	form := url.Values{
		"chat_id":    []string{chatID},
		"text":       []string{text},
		"parse_mode": []string{"MarkdownV2"},
	}

	var chResult = make(chan SendResult)

	// 调用发送。因为使用了 chan 传输结果，所以需要开启新的协程发送
	go bot.Send(urlSendMsg, bytes.NewReader([]byte(form.Encode())),
		"application/x-www-form-urlencoded", chResult)

	// 等待发送完成，分情况处理
	result := <-chResult
	// 有发送错误
	if result.Error != nil && result.Error != ErrResend {
		return nil, fmt.Errorf("[%s]%s", tag, result.Error)
	}
	// 需要重发
	if result.Error == ErrResend {
		return bot.SendMessage(chatID, text)
	}

	// 发送成功
	return result.Message, nil
}

// SendMediaGroup 发送一个媒体集
//
// *只设置第一个媒体的`Caption`时，将作为该集的标题*，所有媒体可以设置`Name`属性
//
// 注意使用 EscapeMk 来转义文本中的非法字符（即属于 Markdown 字符，而不想当做 Markdown 字符渲染）
//
// 因为原生 api 限制发送文件的大小，若需发送大文件，可以运行本地 TG 服务,
// 设置 tg.SetAddr("http://127.0.0.1:1234")后，来发送
// @see https://stackoverflow.com/a/75012096
// @see https://hdcola.medium.com/telegram-bot-api-server%E4%BD%9C%E5%BC%8A%E6%9D%A1-301d40bd65ba
func (bot *TGBot) SendMediaGroup(chatID string, medias []*InputMedia) (*Message, error) {
	tag := "SendMediaGroup"
	// 在此处正确关闭输入流的 Reader
	defer func() {
		for _, m := range medias {
			if r, ok := m.Media.(io.ReadCloser); ok {
				r.Close()
			}
			if r, ok := m.Thumbnail.(io.ReadCloser); ok {
				r.Close()
			}
		}
	}()

	// 使用 bytes.Buffer{} 还是会将数据全部写入内存，所以使用 pipe 替代
	pr, pw := io.Pipe()
	// 组装 multipart 文件上传请求的参数
	writer := multipart.NewWriter(pw)

	// 接收执行请求的错误
	var chResult = make(chan SendResult)

	// 新协程执行请求
	go bot.Send(urlSendMediaGroup, pr, writer.FormDataContentType(), chResult)

	// 添加 chat_id 字段
	err := writer.WriteField("chat_id", chatID)
	if err != nil {
		return nil, fmt.Errorf("[%s]写入 chat_id 出错：%w", tag, err)
	}

	// 写入媒体
	var mediasForm = make([]MediaForm, len(medias))
	for i, m := range medias {
		// 写入媒体
		partMedia, err := writer.CreateFormFile(fmt.Sprintf("media%d", i), m.Name)
		if err != nil {
			return nil, fmt.Errorf("[%s]创建表单信息出错：%w", tag, err)
		}

		_, err = io.Copy(partMedia, m.Media)
		if err != nil {
			return nil, fmt.Errorf("[%s]复制文件流出错：%w", tag, err)
		}

		// 写入缩略图
		if m.Thumbnail != nil {
			partThumbnail, err := writer.CreateFormFile(fmt.Sprintf("thumb%d", i), m.Name)
			if err != nil {
				return nil, fmt.Errorf("[%s]创建表单信息出错：%w", tag, err)
			}

			_, err = io.Copy(partThumbnail, m.Thumbnail)
			if err != nil {
				return nil, fmt.Errorf("[%s]复制文件流出错：%w", tag, err)
			}
		}

		// 设置默认的标题解析模式 MarkdownV2
		if m.MediaData.ParseMode == "" {
			m.MediaData.ParseMode = "MarkdownV2"
		}
		mediaForm := MediaForm{
			MediaData: m.MediaData,
			Media:     fmt.Sprintf("attach://media%d", i),
			Thumbnail: fmt.Sprintf("attach://thumb%d", i),
		}
		mediasForm[i] = mediaForm
	}

	// 发送媒体组的额外信息（标题、对应媒体等）
	mediaFormBs, err := json.Marshal(mediasForm)
	if err != nil {
		return nil, fmt.Errorf("[%s]序列化媒体组的信息出错：%w", tag, err)
	}

	err = writer.WriteField("media", string(mediaFormBs))
	if err != nil {
		return nil, fmt.Errorf("[%s]写入 media 出错：%w", tag, err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("[%s]关闭 multipart writer 出错：%w", tag, err)
	}

	// 注意要关闭 pipe writer 否则会卡主
	err = pw.Close()
	if err != nil {
		return nil, fmt.Errorf("[%s]关闭 pipe writer 出错：%w", tag, err)
	}

	// 发送完数据，等待执行完成，分情况处理
	result := <-chResult
	// 有发送错误
	if result.Error != nil && result.Error != ErrResend {
		return nil, fmt.Errorf("[%s]%s", tag, result.Error)
	}
	// 需要重发
	if result.Error == ErrResend {
		return bot.SendMediaGroup(chatID, medias)
	}

	// 发送成功
	return result.Message, nil
}

// EscapeMk 转义标题中不想渲染为 Markdown V2 的字符。用于转义已经被 Markdown 字符包围的文本
//
// 用法：EscapeMk("测#试Markdown文本*消息*结束：") + "*[搜索](https://www.google.com/)* #标签"
//
// 加号前一段将转义，不渲染为 Markdown；后一段将作为 Markdown 渲染。
//
// 即结果："测#试Markdown文本*消息*结束：搜索 #标签"。其中“搜索”的字体会加粗
//
// 参考：https://core.telegram.org/bots/api#markdownv2-style
func EscapeMk(text string) string {
	// 已替换'['，就不用替换']'了
	reg := regexp.MustCompile("([_*\\[\\]()~`>#+\\-=|{}.!])")
	return reg.ReplaceAllString(text, "\\${1}")
}

// LegalMk 合法化标题中的非法 Markdown V2 字符。用于转义需要作为 Markdown 渲染的文本
//
// 否则，直接发送会报错，提示需要转义，如'\#'
func LegalMk(text string) string {
	reg := regexp.MustCompile("([#])")
	return reg.ReplaceAllString(text, "\\${1}")
}

// SendVideo 发送视频
//
// fileSizeThreshold 设置视频分段的的字节数，为 0 不分段
//
// tmpDir 设置临时文件的目录（为空""则在文件同目录）
//
// reserve 是否保留原文件（或转码后的视频文件）
func (bot *TGBot) SendVideo(chatID string, title string, path string,
	fileSizeThreshold int64, tmpDir string, reserve bool) (*Message, error) {
	tag := "SendVideo"
	// 发送完后，删除临时文件
	var delFiles = make(map[string]string)
	defer func() {
		for p := range delFiles {
			_ = os.Remove(p)
		}
	}()

	// 获取文件大小
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("[%s]获取文件信息出错：%w", tag, err)
	}
	// 如果目标视频超过了设置的最大值，就切割
	dstPaths := []string{path}
	if fileSizeThreshold != 0 && info.Size() > fileSizeThreshold {
		dstPaths, err = dovideo.Cut(path, fileSizeThreshold, tmpDir)
	}

	// 需要发送媒体
	var medias = make([]*InputMedia, len(dstPaths))
	for i, p := range dstPaths {
		// 创建媒体信息
		media, dst, thumb, err := GenTgMedia(p, "")
		if err != nil {
			return nil, fmt.Errorf("[%s]生成媒体信息：%w", tag, err)
		}
		if i == 0 {
			media.Caption = title
		}
		// 多分段时，加上 Pn 标识
		if len(dstPaths) >= 2 {
			media.Name = fmt.Sprintf("P%02d", i+1)
		}
		medias[i] = media

		// 封面图，需要删除
		delFiles[thumb] = ""
		// 分段大于2，说明是切割后的视频列表，发送成功后需要删除
		if len(dstPaths) >= 2 {
			delFiles[dst] = ""
			delFiles[p] = ""
		}
		// 不保留原文件，删除
		if !reserve {
			delFiles[dst] = ""
		}
	}

	// 发送
	return bot.SendMediaGroup(chatID, medias)
}
