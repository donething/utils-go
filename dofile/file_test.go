package dofile

import (
	"encoding/hex"
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
			got, err := Exists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
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
			args{[]byte{54, 55}, `E:/Temp/go/utils-go/dofile/test_cretae.txt`, WriteCreate, 0644},
			2,
			false,
		},
		{
			"Test Append",
			args{[]byte{54, 55}, `E:/Temp/go/utils-go/dofile/test_append.txt`, WriteAppend, 0644},
			2,
			false,
		},
		{
			"Test Trunc",
			args{[]byte{54, 55}, `E:/Temp/go/utils-go/dofile/test_trunc.txt`, WriteTrunc, 0644},
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

func TestCopyFile(t *testing.T) {
	n, err := CopyFile("D:/MyData/Setting/Windows/bash.bashrc", "E:/Temp/bash.bashrc", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("文件复制完成：", n)
}

func TestReadHeadTailBytes(t *testing.T) {
	n1, n2, err := ReadHeadTailBytes("D:/Users/Doneth/Downloads/006Yjd8ogy1g0zvj6ny49g304604oe1g.gif", 0, 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s,%s\n", hex.EncodeToString(n1), hex.EncodeToString(n2))
}
