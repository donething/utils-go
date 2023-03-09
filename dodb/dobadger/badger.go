// Package dobadger Fast key-value DB badger 的帮助函数
//
// @see https://github.com/dgraph-io/badger
// @see https://dgraph.io/docs/badger/get-started/
// @see https://blog.csdn.net/qq_42057154/article/details/123739565
package dobadger

import (
	"github.com/dgraph-io/badger/v4"
	"strings"
)

// DoBadger badger 的包装。执行数据操作后不会关闭数据库，可按需调用`db.DB.Close()`关闭。还可以在程序退出时，才关闭数据库
type DoBadger struct {
	DB *badger.DB
}

// Open 根据数据库路径打开数据库。当 optsNullable 为 nil 时，使用默认值
func Open(dbDirPath string, optsNullable *badger.Options) (*DoBadger, error) {
	// 处理选项
	var opts = optsNullable
	if opts == nil {
		op := badger.DefaultOptions(dbDirPath)
		opts = &op
		opts.CompactL0OnClose = true
		opts.ValueLogFileSize = 1024 * 1024 * 10
		opts.Logger = nil
	}

	// 打开数据库
	var err error
	db, err := badger.Open(*opts)

	return &DoBadger{DB: db}, err
}

// Close 关闭数据库。请仅用此函数关闭，不要再其它地方调用 db.DB.Close() 来关闭数据库
func (db *DoBadger) Close() error {
	if db.DB != nil {
		err := db.DB.Close()
		db.DB = nil
		return err
	}

	return nil
}

// Get 获取数据。不存在该键时返回错误`ErrKeyNotFound`
func (db *DoBadger) Get(key []byte) ([]byte, error) {
	var valCopy []byte
	err := db.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		// Alternatively, you could also use item.ValueCopy().
		valCopy, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, err
	}

	return valCopy, nil
}

// Has 是否存在指定的键
func (db *DoBadger) Has(key []byte) (bool, error) {
	_, err := db.Get(key)

	if err == badger.ErrKeyNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Set 存放数据
func (db *DoBadger) Set(key []byte, value []byte) error {
	return db.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Del 删除数据
func (db *DoBadger) Del(key []byte) error {
	return db.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Query 查找匹配的键及其数据
//
// param keySubStr 为需要包含的子字符串，当不为空""时，需要数据库中的键名为 string 类型;
// 当 keySubStr 为 空""时，返回所有数据
func (db *DoBadger) Query(keySubStr string) (map[string][]byte, error) {
	payload := make(map[string][]byte)

	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if keySubStr != "" &&
				!strings.Contains(strings.ToLower(string(item.Key())), strings.ToLower(keySubStr)) {
				continue
			}

			valueCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			payload[string(item.Key())] = valueCopy
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return payload, err
}

// QueryPrefix 扫描前缀和关键字
//
// keySub 为空""时，返回带有指定 prefixStr 的所有项
//
// 针对多个场景共用一个数据库，根据前缀区分场景的情况
func (db *DoBadger) QueryPrefix(prefixStr string, keySub string) (map[string][]byte, error) {
	payload := make(map[string][]byte)

	err := db.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(prefixStr)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			if keySub != "" &&
				!strings.Contains(strings.ToLower(string(item.Key())), strings.ToLower(keySub)) {
				continue
			}

			valueCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			payload[string(item.Key())] = valueCopy
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return payload, nil
}

// BatchSet 批量插入数据
func (db *DoBadger) BatchSet(data map[string][]byte) error {
	wb := db.DB.NewWriteBatch()
	defer wb.Cancel()

	for key, bs := range data {
		errPut := wb.Set([]byte(key), bs)
		if errPut != nil {
			return errPut
		}
	}
	return wb.Flush()
}
