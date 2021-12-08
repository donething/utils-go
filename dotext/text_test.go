package dotext

import (
	"golang.org/x/text/encoding/unicode"
	"testing"
	"time"
)

func TestUTF8ToGBK(t *testing.T) {
	str := "UTF8 和 GBK 编码转换测试"
	t.Log(str)

	gbk, err := UTF8ToGBK([]byte(str))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(gbk))

	utf8, err := GBKToUTF8(gbk)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(utf8))

	if string(utf8) != str {
		t.Fatal("UTF8 和 GBK 编码转换测试：编码转换失败")
	}
	t.Log("UTF8 和 GBK 编码转换测试：编码转换成功")
}

func TestUTF16ToUTF8(t *testing.T) {
	str := "UTF8 和 UTF16 编码转换测试"
	t.Log(str)

	utf16, err := UTF8ToUTF16([]byte(str), unicode.LittleEndian, unicode.IgnoreBOM)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(utf16))

	utf8, err := UTF16ToUTF8(utf16, unicode.LittleEndian, unicode.IgnoreBOM)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(utf8))

	if string(utf8) != str {
		t.Fatal("UTF8 和 UTF16 编码转换测试：编码转换失败")
	}
	t.Log("UTF8 和 UTF16 编码转换测试：编码转换成功")
}

func TestFormatDate(t *testing.T) {
	type args struct {
		t      time.Time
		format string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test FormatDate",
			args: args{
				time.Date(2018, 11, 11, 1, 10, 21, 0, time.Local),
				TimeFormat,
			},
			want: "2018-11-11 01:10:21",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDate(tt.args.t, tt.args.format); got != tt.want {
				t.Errorf("FormatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeiJingTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Test BeiJingTime",
			args: args{
				time.Date(2018, 11, 11, 1, 10, 21, 0,
					time.UTC),
			},
			want: time.Date(2018, 11, 11, 1, 10, 21, 0,
				time.UTC).Add(8 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BeiJingTime(tt.args.t)
			if got.Hour() != tt.want.Hour() ||
				got.Minute() != tt.want.Minute() ||
				got.Second() != tt.want.Second() {
				t.Errorf("BeiJingTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64Encode(t *testing.T) {
	var str = "测试 Base64 编码转换"
	var base64 = Base64Encode(str)
	var source, err = Base64Decode(base64)
	if err != nil {
		t.Fatal(err)
	}

	if string(source) != str {
		t.Fatal("测试 Base64 编码转换 编码转换失败：", string(source))
	}
	t.Log("测试 Base64 编码转换：编码转换成功")
}

func TestBytesHumanReadable(t *testing.T) {
	type args struct {
		bytes int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test BytesHumanReadable",
			args: args{bytes: 123456789},
			want: "117.74 MB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesHumanReadable(tt.args.bytes); got != tt.want {
				t.Errorf("BytesHumanReadable() = %v, want %v", got, tt.want)
			}
		})
	}
}
