// Package dotext 文本处理
// 将 UTF8 编码转为 GBK、UTF16 等为 编码(encode)，反之为解码(decode)
package dotext

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"time"
)

const (
	// TimeFormat 转换时间的常用格式
	TimeFormat = "2006-01-02 15:04:05"
)

// GBKToUTF8 GBK 编码转 UTF-8
func GBKToUTF8(s []byte) ([]byte, error) {
	// 编码转换：http://mengqi.info/html/2015/201507071345-using-golang-to-convert-text-between-gbk-and-utf-8.html
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF8ToGBK UTF-8 编码转 GBK
func UTF8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF16ToUTF8 将 UTF16 编码转换为 UTF8 编码
//
// 参数 endian: unicode.LittleEndian、unicode.BigEndian
//
// 参数 bom: unicode.IgnoreBOM、unicode.UseBOM、unicode.ExpectBOM
//
// 参考：https://gist.github.com/bradleypeabody/185b1d7ed6c0c2ab6cec#gistcomment-2780177
// 参考 https://blog.csdn.net/wyzxg/article/details/5349896
func UTF16ToUTF8(bs []byte, endian unicode.Endianness, bom unicode.BOMPolicy) ([]byte, error) {
	decoder := unicode.UTF16(endian, bom).NewDecoder()
	bs8, err := decoder.Bytes(bs)
	return bs8, err
}

// UTF8ToUTF16 将 UTF8 编码转换为 UTF16 编码
//
// 参数 endian: unicode.LittleEndian、unicode.BigEndian
//
// 参数 bom: unicode.IgnoreBOM、unicode.UseBOM、unicode.ExpectBOM
//
// 参考 https://forum.golangbridge.org/t/how-to-convert-utf-8-string-to-utf-16-be-string/7072/2
// 参考 https://blog.csdn.net/wyzxg/article/details/5349896
func UTF8ToUTF16(bs []byte, endian unicode.Endianness, bom unicode.BOMPolicy) ([]byte, error) {
	encoder := unicode.UTF16(endian, bom).NewEncoder()
	bs16, err := encoder.Bytes(bs)
	return bs16, err
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

// BytesHumanReadable 将文件大小的字节转为可读的字符（如"102 MB"）
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
