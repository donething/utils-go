package dofile

import (
	"log"
	"testing"
)

func TestMd5(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		wantMd5Str string
		wantErr    bool
	}{
		{
			"Test md5",
			args{"E:/Temp/20190226_PolarBearDay_ZH-CN5185516722_1920x1080.jpg"},
			"bf9d0939df1039e2893b1004d5a169d7",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMd5Str, err := Md5(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Md5() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMd5Str != tt.wantMd5Str {
				t.Errorf("Md5() = %v, want %v", gotMd5Str, tt.wantMd5Str)
			}
		})
	}
}

func Test_isDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"Test File",
			args{"E:/Temp/20190226_PolarBearDay_ZH-CN5185516722_1920x1080.jpg"},
			false,
			false,
		},
		{
			"Test Dir",
			args{"E:/Temp"},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("isDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"Test exists file",
			args{"E:/Temp"},
			true,
			false,
		},
		{
			"Test not exists file",
			args{"E:/Temp/123"},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PathExists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PathExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckIntegrity(t *testing.T) {
	result, err := CheckIntegrity("E:/Temp/20190226_PolarBearDay_ZH-CN5185516722_1920x1080.jpg")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(result)
}
