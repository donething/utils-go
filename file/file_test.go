package file

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	file, err := Get("/home/doneth/MyData")
	if err != nil {
		t.Fatal(err)
	}
	filesList, err := file.List(DIR)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v\n", filesList)
}

func TestDFile_ListPaths(t *testing.T) {
	file, err := Get("D:/Temp")
	if err != nil {
		t.Fatal(err)
	}
	pathsList, err := file.List(DIR)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v\n", pathsList[0])
}

func TestDFile_Md5(t *testing.T) {
	type fields struct {
		Path     string
		FileInfo os.FileInfo
	}
	path := "D:/Temp/ishike.json"
	info, _ := os.Stat(path)

	tests := []struct {
		name       string
		fields     fields
		wantMd5Str string
		wantErr    bool
	}{
		{
			name:       "计算文件md5",
			fields:     fields{Path: path, FileInfo: info},
			wantMd5Str: "47982d3a26b4fcf0cd84aced7eaef2af",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &dFile{
				Path:     tt.fields.Path,
				FileInfo: tt.fields.FileInfo,
			}
			gotMd5Str, err := f.Md5()
			if (err != nil) != tt.wantErr {
				t.Errorf("dFile.Md5() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMd5Str != tt.wantMd5Str {
				t.Errorf("dFile.Md5() = %v, want %v", gotMd5Str, tt.wantMd5Str)
			}
		})
	}
}
