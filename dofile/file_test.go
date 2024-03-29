package dofile

import (
	"os"
	"reflect"
	"testing"
)

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
			got, err := IsDir(tt.args.path)
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
			args{[]byte{54, 55}, `E:/Temp/go/utils-go/dofile/test_cretae.txt`, OCreate, 0644},
			2,
			false,
		},
		{
			"Test Append",
			args{[]byte{54, 55}, `E:/Temp/go/utils-go/dofile/test_append.txt`, OAppend, 0644},
			2,
			false,
		},
		{
			"Test Trunc",
			args{[]byte{54, 55}, `E:/Temp/go/utils-go/dofile/test_trunc.txt`, OTrunc, 0644},
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
	err := ShowInExplorer(`D:/1925 年北洋“中国丧失领土领海图.jpg`)
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

func TestValidFileName1(t *testing.T) {
	type args struct {
		src  string
		repl string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{`query|User\Info《Seven`, ""}, "queryUserInfo《Seven"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidFileName(tt.args.src, tt.args.repl); got != tt.want {
				t.Errorf("ValidFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenAs(t *testing.T) {
	err := OpenAs(`D:/1925 年北洋“中国丧失领土领海图.jpg`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMD5(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "测试",
			args:    args{path: "D:/Tmp/VpsGo/video.mkv"},
			want:    "484e19309253d455917ff1bae9de04d3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Md5(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Md5() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Md5() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniquePath(t *testing.T) {
	path1 := UniquePath("D:/Tmp/live/test.mp4")
	path2 := UniquePath("D:/Tmp/live/test")

	t.Log(path1, path2)
}
