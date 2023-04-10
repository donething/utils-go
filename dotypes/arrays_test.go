package dotypes

import (
	"reflect"
	"testing"
)

type Fruit struct {
	Name string
}

func TestFindIndex(t *testing.T) {
	type args[T any] struct {
		data  []T
		equal func(a T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[string]{
		{
			name: "比较字符串",
			args: args[string]{
				data:  []string{"苹果", "橘子", "香蕉"},
				equal: func(item string) bool { return item == "橘子" },
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindIndex(tt.args.data, tt.args.equal); got != tt.want {
				t.Errorf("FindIndex() = %v, want %v", got, tt.want)
			}
		})
	}

	tests2 := []testCase[Fruit]{
		{
			name: "比较对象",
			args: args[Fruit]{
				data:  []Fruit{{Name: "苹果"}, {Name: "橘子"}},
				equal: func(item Fruit) bool { return item.Name == "橘子" },
			},
			want: 1,
		},
	}

	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindIndex(tt.args.data, tt.args.equal); got != tt.want {
				t.Errorf("FindIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelItem(t *testing.T) {
	type args[T any] struct {
		data  []T
		equal func(item T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "删除数组中的整数",
			args: args[int]{
				data:  []int{1, 2, 3},
				equal: func(item int) bool { return item == 2 },
			},
			want: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DelItem(tt.args.data, tt.args.equal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DelItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
