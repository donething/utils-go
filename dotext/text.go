// Package dotext 文本处理
// 将 UTF8 编码转为 GBK、UTF16 等为 编码(encode)，反之为解码(decode)
package dotext

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"regexp"
	"time"
)

const (
	// TimeFormat 转换时间的常用格式
	TimeFormat = "2006-01-02 15:04:05"
)

// HasUTF8BOM 判断指定 UTF-8 编码的文本数据是否含 BOM
// 参考：https://www.jianshu.com/p/5d8771da218b
func HasUTF8BOM(bs []byte) bool {
	if len(bs) >= 3 && bs[0] == 239 && bs[1] == 187 && bs[2] == 191 {
		return true
	}
	return false
}

// DetectFileCoding 检测指定路径的文本文件的编码
//
// 返回 编码、地区、准确度（如 GB-18030、zh、100）
func DetectFileCoding(path string) (*chardet.Result, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DetectTextCoding(bs)
}

// DetectTextCoding 检测文本的编码，返回 编码、地区、准确度（如 GB-18030、zh、100）
//
// [chardet: Charset detector library for golang derived from ICU](https://github.com/saintfish/chardet)
func DetectTextCoding(data []byte) (result *chardet.Result, err error) {
	detector := chardet.NewTextDetector()
	return detector.DetectBest(data)
}

// TransformText 转换编码为无 BOM 的 UTF-8
//
// 返回结果、原编码、可能的错误
func TransformText(bs []byte) ([]byte, string, error) {
	// 检测文本的编码
	result, err := DetectTextCoding(bs)
	if err != nil {
		return nil, "", err
	}
	// 若本来就是无 BOM 的 UTF-8 编码，不需修改直接返回
	if result.Charset == "UTF-8" && !HasUTF8BOM(bs) {
		return nil, result.Charset, nil
	}

	// 按指定编码读取数据为 UTF-8 编码（可能含有 BOM，需要通过下面的方法去除）
	// 参考：https://stackoverflow.com/a/44298295
	byteReader := bytes.NewReader(bs)
	reader, _ := charset.NewReaderLabel(result.Charset, byteReader)

	// 读取结果
	nbs, err := ioutil.ReadAll(reader)

	// 如果是带有 BOM UTF-8，去除 BOM
	if HasUTF8BOM(nbs) {
		nbs = nbs[3:]
	}

	return nbs, result.Charset, err
}

// FormatDate 格式化时间
//
// 参数 format 为时间的格式，可使用 dostr.TimeFormat
func FormatDate(t time.Time, format string) string {
	return t.Format(format)
}

// BeiJingTime 将当前时间转为北京时间
func BeiJingTime(t time.Time) time.Time {
	// 东八区
	var cstZone = time.FixedZone("GMT", 8*3600)
	return t.UTC().In(cstZone)
}

// Base64Encode base64 编码
func Base64Encode(str string) string {
	// https://stackoverflow.com/a/28672789/8179418
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64 解码
func Base64Decode(str string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// BytesHumanReadable 将文件大小的字节转为可读的字符，如"102 MB"
//
// https://stackoverflow.com/a/30822306
func BytesHumanReadable(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ResolveFanhao 解析文本中的番号
//
// 若没有找到番号就返回""
func ResolveFanhao(text string) string {
	// 正则写3个小括号，是为了之后判断番号的数字部分是否为3位
	reg := regexp.MustCompile(`([a-zA-Z]+)([-_\s]?)([0-9]+)`)
	result := reg.FindStringSubmatch(text)
	if result != nil {
		// 如果番号的数字部分不为3位，则需要用"0"填充
		if len(result[3]) == 1 {
			result[3] = "00" + result[3]
		} else if len(result[3]) == 2 {
			result[3] = "0" + result[3]
		}
		return result[1] + "-" + result[3]
	}
	return ""
}
