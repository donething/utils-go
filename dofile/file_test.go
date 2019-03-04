package dofile

import (
	"log"
	"os"
	"reflect"
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

func TestRead(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantBs  []byte
		wantErr bool
	}{
		{
			"Read Exist",
			args{`E:/Temp/temp.txt`},
			[]byte{49, 50, 51},
			false,
		},
		{
			"Read Not Exist",
			args{`E:/Temp/temp12.txt`},
			[]byte{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBs, err := Read(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBs, tt.wantBs) {
				t.Errorf("Read() = %v, want %v", gotBs, tt.wantBs)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	type args struct {
		bs   []byte
		path string
		mode int
		perm os.FileMode
	}
	tests := []struct {
		name    string
		args    args
		wantN   int
		wantErr bool
	}{
		{
			"Test Create",
			args{[]byte{54, 55}, `E:/Temp/go/utils-go/dofile/test_cretae.txt`, WRITE_CREATE, 0644},
			2,
			false,
		},
		{
			"Test Append",
			args{[]byte{54, 55}, `E:/Temp/go/utils-go/dofile/test_append.txt`, WRITE_APPEND, 0644},
			2,
			false,
		},
		{
			"Test Trunc",
			args{[]byte{54, 55}, `E:/Temp/go/utils-go/dofile/test_trunc.txt`, WRITE_TRUNC, 0644},
			2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := Write(tt.args.bs, tt.args.path, tt.args.mode, tt.args.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestShowInExplorer(t *testing.T) {
	err := ShowInExplorer(`E:/Temp/get_file.txt`)
	if err != nil {
		t.Fatal(err)
	}
}
