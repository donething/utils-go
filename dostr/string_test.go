package dostr

import (
	"testing"
	"time"
)

func TestUTF8ToGBK(t *testing.T) {
	str := "UTF8和GBK编码转换测试"
	gbk, err := UTF8ToGBK([]byte(str))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(gbk))

	utf8, err := GBKToUTF8([]byte(gbk))
	if err != nil {
		t.Fatal(err)
	}
	if string(utf8) != str {
		t.Error("编码转换失败")
	}
	t.Log("编码转换成功")
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
			name: "default_format",
			args: args{time.Date(2018, 11, 11, 0, 49, 1, 0, time.Local), ""},
			want: "2018-11-11 00:49:01",
		},
		{
			name: "2006-01-02 15:04:05",
			args: args{time.Date(2018, 11, 11, 1, 10, 21, 0, time.Local), "2006/01/02 15:04:05"},
			want: "2018/11/11 01:10:21",
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
