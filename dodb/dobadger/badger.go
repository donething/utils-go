// Package dobadger Fast key-value DB badger 的帮助函数
//
// @see https://github.com/dgraph-io/badger
// @see https://dgraph.io/docs/badger/get-started/
// @see https://blog.csdn.net/qq_42057154/article/details/123739565
package dobadger

import (
	"github.com/dgraph-io/badger/v3"
	"log"
	"strings"
)

// DoBadger badger 的包装
type DoBadger struct {
	DB *badger.DB
}

// Open 根据数据库路径打开数据库，当 Options 为 nil 时，使用默认 Options
func Open(dbDirPath string, optsNullable *badger.Options) *DoBadger {
	// 处理选项
	var opts badger.Options
	if optsNullable != nil {
		opts = *optsNullable
	} else {
		opts = badger.DefaultOptions(dbDirPath)
		opts.CompactL0OnClose = true
		opts.ValueLogFileSize = 1024 * 1024 * 10
		opts.Logger = nil
	}

	// 打开数据库
	var err error
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("打开数据库'%s'出错：%s\n", dbDirPath, err)
	}

	return &DoBadger{DB: db}
}

// Close 关闭数据库，请仅用此函数关闭，不要再其它地方调用 db.Close() 来关闭数据库
func (db *DoBadger) Close() error {
	if db.DB == nil {
		return nil
	}

	return db.DB.Close()
}

// Get 获取数据
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
			if keySubStr == "" &&
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

// QueryPrefix 前缀扫描
func (db *DoBadger) QueryPrefix(prefixStr string) (map[string][]byte, error) {
	payload := make(map[string][]byte)

	err := db.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(prefixStr)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
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
