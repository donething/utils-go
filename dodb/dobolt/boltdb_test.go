package dobolt

import (
	"fmt"
	"reflect"
	"testing"
)

var db *DoBolt

var (
	bucketName = []byte("classes")

	classesFruit      = []byte("fruit")
	classesFruitValue = []byte("Fruit")

	classesTools      = []byte("tools")
	classesToolsValue = []byte("Tools")
)

func init() {
	var err error
	db, err = Open("./dbtest.db", nil, nil)
	if err != nil {
		panic(err)
	}
}

func TestCreate(t *testing.T) {
	err := db.Create(bucketName)
	if err != nil {
		t.Errorf("创建桶失败：%s\n", err)
	}
	t.Logf("创建桶成功")
	db.Close()
}

func TestGet(t *testing.T) {
	defer db.Close()
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
			got, err := db.Get(tt.args.key, tt.args.bucket)
			fmt.Printf("获取的数据：%s\n", got)
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
			if err := db.Set(tt.args.key, tt.args.value, tt.args.bucket); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.Close()
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
			got, err := db.Del(tt.args.key, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Del() got = %v, want %v", got, tt.want)
			}
		})
	}
	db.Close()
}

func TestQuery(t *testing.T) {
	var f = "fruit"
	type args struct {
		bucket    []byte
		keySubStr string
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
				keySubStr: f,
			},
			want:    map[string][]byte{string(classesFruit): classesFruitValue},
			wantErr: false,
		},
		{
			name: "Test query all",
			args: args{
				bucket:    bucketName,
				keySubStr: "",
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
			got, err := db.Query(tt.args.keySubStr, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
		})
	}
	db.Close()
}

func TestDoBolt_QueryPrefix(t *testing.T) {
	got, err := db.QueryPrefix([]byte("too"), bucketName)
	if err != nil {
		t.Fatal(err)
	}

	for _, bs := range got {
		t.Logf("值：%s\n", string(bs))
	}
}
