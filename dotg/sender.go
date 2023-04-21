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

	// API
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
	bs, err := client.PostForm(fmt.Sprintf(urlSendMsg, bot.addr, bot.token), form.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("[%s]发送消息出错：%w", tag, err)
	}

	// 返回的消息
	var msg Message
	err = json.Unmarshal(bs, &msg)
	if err != nil {
		return nil, fmt.Errorf("[%s]解析响应出错：%w", tag, err)
	}

	return &msg, nil
}

// SendMediaGroup 发送一个媒体集
//
// *只设置第一个媒体的`Caption`时，将作为该集的标题*
//
// 注意使用 EscapeMk 来转义文本中的非法字符（即属于 Markdown 字符，而不想当做 Markdown 字符渲染）
//
// 因为原生 api 限制发送文件的大小，若需发送大文件，可以运行本地 TG 服务,
// 设置 tg.SetAddr("http://127.0.0.1:1234")后，传递大文件的路径(file:///home/output.mp4)来发送
// @see https://stackoverflow.com/a/75012096
// @see https://hdcola.medium.com/telegram-bot-api-server%E4%BD%9C%E5%BC%8A%E6%9D%A1-301d40bd65ba
func (bot *TGBot) SendMediaGroup(chatID string, medias []*InputMedia) (*Message, error) {
	tag := "SendMediaGroup"
	// 正确关闭 Reader
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

	// 组装 multipart 文件上传请求参数
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

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

			_, err = io.Copy(partThumbnail, m.Media)
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

	// 发送
	sendUrl := fmt.Sprintf(urlSendMediaGroup, bot.addr, bot.token)
	resp, err := client.Post(sendUrl, writer.FormDataContentType(), buf)
	if err != nil {
		return nil, fmt.Errorf("[%s]执行发送媒体组的请求出错：%w", tag, err)
	}
	defer resp.Body.Close()

	// 处理返回的消息
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s]读取响应内容出错：%w", tag, err)
	}

	var msg Message
	err = json.Unmarshal(bs, &msg)
	if err != nil {
		return nil, fmt.Errorf("[%s]解析响应内容出错：%w", tag, err)
	}

	// 速率限制
	if msg.ErrorCode == 429 {
		secStr := msg.Description[strings.LastIndex(msg.Description, " ")+1:]
		sec, err := strconv.Atoi(secStr)
		if err != nil {
			return nil, fmt.Errorf("[%s]解析速率限制的等待时长时出错：%w", tag, err)
		}

		fmt.Printf("[%s]由于速率限制，等待 %d 秒后重新发送\n", tag, sec)
		time.Sleep(time.Duration(sec+1) * time.Second)
		return bot.SendMediaGroup(chatID, medias)
	}

	// 发送失败
	if !msg.Ok {
		return &msg, fmt.Errorf("[%s]发送失败：%d %s", tag, msg.ErrorCode, msg.Description)
	}

	return &msg, nil
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
	fileSizeThreshold int64, tmpDir string, reserve bool) error {
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
		return fmt.Errorf("[%s]获取文件信息出错：%w", tag, err)
	}
	// 如果目标视频超过了设置的最大值，就切割
	dstPaths := []string{path}
	if fileSizeThreshold != 0 && info.Size() > fileSizeThreshold {
		dstPaths, err = dovideo.Cut(path, fileSizeThreshold, tmpDir)
	}

	// 需要发送媒体
	var medias = make([]*InputMedia, len(dstPaths))
	for i, p := range dstPaths {
		media, dst, thumb, err := GenTgMedia(p, fmt.Sprintf("%s P%02d", title, i+1))
		if err != nil {
			return fmt.Errorf("[%s]生成媒体信息：%w", tag, err)
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
	_, err = bot.SendMediaGroup(chatID, medias)

	return err
}
