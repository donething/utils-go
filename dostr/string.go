package dostr

import (
	"bytes"
	"encoding/base64"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"time"
)

// GBK编码转UTF-8
func GBKToUTF8(s []byte) ([]byte, error) {
	// 编码转换：http://mengqi.info/html/2015/201507071345-using-golang-to-convert-text-between-gbk-and-utf-8.html
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8编码转GBK
func UTF8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 格式化时间
// 如果format为空白字符，则默认设为"2006-01-02 15:04:05"
func FormatDate(t time.Time, format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return t.Format(format)
}

// base64编码
func Base64Encode(str string) string {
	// https://stackoverflow.com/a/28672789/8179418
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// base64解码
func Base64Decode(str string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return data, nil
}
