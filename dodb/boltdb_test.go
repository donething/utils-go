package dodb

import (
	"reflect"
	"testing"
)

var (
	bucketName = []byte("classes")

	classesFruit      = []byte("fruit")
	classesFruitValue = []byte("Fruit")

	classesTools      = []byte("tools")
	classesToolsValue = []byte("Tools")
)

func init() {
	Open("test.db")
}

func TestCreate(t *testing.T) {
	err := Create(bucketName)
	if err != nil {
		t.Errorf("创建桶失败：%s\n", err)
	}
	t.Logf("创建桶成功")
	Close()
}

func TestGet(t *testing.T) {
	type args struct {
		key    []byte
		bucket []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Test get",
			args: args{
				key:    classesFruit,
				bucket: bucketName,
			},
			want:    classesFruitValue,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.key, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPut(t *testing.T) {
	type args struct {
		key    []byte
		value  []byte
		bucket []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test put 1",
			args: args{
				key:    classesFruit,
				value:  classesFruitValue,
				bucket: bucketName,
			},
		},
		{
			name: "Test put 2",
			args: args{
				key:    classesTools,
				value:  classesToolsValue,
				bucket: bucketName,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Put(tt.args.key, tt.args.value, tt.args.bucket); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDel(t *testing.T) {
	type args struct {
		key    []byte
		bucket []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Test del",
			args: args{
				key:    classesFruit,
				bucket: bucketName,
			},
			want:    classesFruitValue,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Del(tt.args.key, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Del() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery(t *testing.T) {
	var f = "fruit"
	type args struct {
		bucket    []byte
		keySubStr *string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		wantErr bool
	}{
		{
			name: "Test query fruit",
			args: args{
				bucket:    bucketName,
				keySubStr: &f,
			},
			want:    map[string][]byte{string(classesFruit): classesFruitValue},
			wantErr: false,
		},
		{
			name: "Test query all",
			args: args{
				bucket:    bucketName,
				keySubStr: nil,
			},
			want: map[string][]byte{
				string(classesFruit): classesFruitValue,
				string(classesTools): classesToolsValue,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.args.keySubStr, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
		})
	}
}
