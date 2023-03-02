package doconf

import (
	"testing"
)

type DoConf struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Extra string `json:"extra"`
}

const confPath = "./test_conf.json"

var Conf DoConf

func TestInit(t *testing.T) {
	exist, err := Init(confPath, &Conf)
	if err != nil {
		t.Fatal(err)
	}

	if !exist {
		t.Logf("似乎是首次运行，请先填写配置后，再运行程序")
		return
	}

	t.Logf("配置：%+v\n", Conf)
}
