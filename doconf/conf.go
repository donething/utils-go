package doconf

import (
	"encoding/json"
	"github.com/donething/utils-go/dofile"
	"os"
)

// Init 初始化配置
//
// 返回是否存在配置，以便提示填写配置后再运行。返回`false`需要终止程序以填写配置；为`true`才可继续执行
func Init[T any](confPath string, pConf *T) (bool, error) {
	// 需要先判断配置文件是否存在，以免后面覆盖文件后，影响判断
	exist, err := dofile.Exists(confPath)
	if err != nil {
		return false, err
	}

	// 先判断配置文件是否存在：
	// 不存在时，根据`pConf`生成配置文件后，提示填写配置后，重新运行程序；
	// 存在时，先读取配置文件到`pConf`，再重命名配置为`*.json.bak`，作为文件备份
	// 然后根据新的`pConf`结构，重新产生配置文件`*.json`

	if exist {
		// 读取配置文件
		bs, err := dofile.Read(confPath)
		if err != nil {
			return false, err
		}

		err = json.Unmarshal(bs, pConf)
		if err != nil {
			return false, err
		}

		// 生成`*.json.bak`
		err = os.Rename(confPath, confPath+".bak")
		if err != nil {
			return false, err
		}
	}

	// 生成新的`*.json`
	bs, err := json.MarshalIndent(*pConf, "", "  ")
	if err != nil {
		return false, err
	}

	_, err = dofile.Write(bs, confPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}

	return exist, nil
}
