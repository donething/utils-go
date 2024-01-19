package dotext

import (
	"github.com/saintfish/chardet"
	"reflect"
	"testing"
	"time"
)

func TestDetectFileCoding(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *chardet.Result
		wantErr bool
	}{
		{
			"UTF-8",
			args{"D:/影视/连续剧/《金牌冰人》国语外挂GOTV_720P_TS_800M/01.srt"},
			&chardet.Result{"UTF-8", "", 100},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectFileCoding(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectFileCoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectFileCoding() got = %v, want %v", got, tt.want)
			}
		})
	}
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
		bytes uint64
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

func TestFile2UTF8(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   string
		wantErr bool
	}{
		{
			name:    "Test GB18030",
			args:    args{path: "D:/Data/GoTest/入殓师.ssa"},
			want:    true,
			want1:   "GB18030",
			wantErr: false,
		},
		{
			name:    "Test UTF-16LE",
			args:    args{path: "D:/Data/GoTest/u16le.txt"},
			want:    true,
			want1:   "UTF-16LE",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := File2UTF8(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("File2UTF8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("File2UTF8() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("File2UTF8() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
