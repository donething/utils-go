package dostr

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"time"
	"unicode/utf16"
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

// 将UTF16编码转换为UTF8编码
// 参数order: binary.LittleEndian、binary.BigEndian等
// 参考：https://gist.github.com/bradleypeabody/185b1d7ed6c0c2ab6cec#gistcomment-2780177
func UTF16ToUTF8(bs []byte, order binary.ByteOrder) ([]rune, error) {
	ints := make([]uint16, len(bs)/2)
	if err := binary.Read(bytes.NewReader(bs), order, &ints); err != nil {
		return nil, err
	}
	// 可通过string()将其转换为string
	// 再通过[]byte()将string转为[]byte
	return utf16.Decode(ints), nil
}

// 将UTF8编码转换为UTF16编码
// 参数：endian: unicode.LittleEndian、unicode.LittleEndian
// 参数：bom: unicode.IgnoreBOM、unicode.UseBOM、unicode.ExpectBOM
// 参考：https://forum.golangbridge.org/t/how-to-convert-utf-8-string-to-utf-16-be-string/7072/2
func UTF8ToUTF16(bs []byte, endian unicode.Endianness, bom unicode.BOMPolicy) ([]byte, error) {
	decoder := unicode.UTF16(endian, bom).NewDecoder()
	bs16, err := decoder.Bytes(bs)
	return bs16, err
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
