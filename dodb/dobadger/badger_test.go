package dobadger

import (
	"reflect"
	"testing"
)

var db *DoBadger

func init() {
	var err error
	db, err = Open("./dbtest", nil)
	if err != nil {
		panic(err)
	}
}

func TestSet(t *testing.T) {
	type args struct {
		key   []byte
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test set 1",
			args: args{
				key:   []byte("test123"),
				value: []byte("test123测试123"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.Close()
}

func TestBatchSet(t *testing.T) {
	type args struct {
		data map[string][]byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test BatchSet 1",
			args: args{data: map[string][]byte{
				"test1": []byte("test111"),
				"test2": []byte("test222"),
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.BatchSet(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("BatchSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	db.Close()
}

func TestGet(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Test GetBytes 1",
			args:    args{key: []byte("test123")},
			want:    []byte("test123测试123"),
			wantErr: false,
		},
		{
			name:    "Test GetBytes 2",
			args:    args{key: []byte("test2")},
			want:    []byte("test222"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBytes() got = %v, want %v", got, tt.want)
			}
		})
	}

	db.Close()
}

func TestDel(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test Del 1",
			args:    args{key: []byte("test123")},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.Del(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	db.Close()
}

func TestQuery(t *testing.T) {
	got, err := db.Query("")
	if err != nil {
		t.Errorf("fail: %s\n", err)
		return
	}

	for key, bs := range got {
		t.Logf("%s: %s\n", key, string(bs))
	}

	db.Close()
}

func TestQueryPrefix(t *testing.T) {
	got, err := db.QueryPrefix("test1", "")
	if err != nil {
		t.Errorf("fail: %s\n", err)
		return
	}

	for key, bs := range got {
		t.Logf("%s: %s\n", key, string(bs))
	}

	db.Close()
}
