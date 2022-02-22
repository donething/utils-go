package dodb

import (
	"github.com/boltdb/bolt"
	"log"
	"strings"
	"time"
)

var (
	// 数据库的实例
	db *bolt.DB
)

// Open 根据数据库路径打开数据库
func Open(dbPath string) {
	// 打开数据库
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Fatalf("打开数据库'%s'出错：%s\n", dbPath, err)
	}
}

// Close 关闭数据库，请仅用此函数关闭，不要再其它地方调用 db.Close() 来关闭数据库
func Close() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatalf("关闭数据库出错：%s\n", err)
		}
	}
}

// Create 创建桶
func Create(bucket []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, errC := tx.CreateBucketIfNotExists(bucket)
		return errC
	})
}

// Get 获取桶中键对应的值
func Get(key []byte, bucket []byte) ([]byte, error) {
	// open a Read-only transaction with the first argument `false`
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}

	// do something ...
	got := tx.Bucket(bucket).Get(key)
	// 对值需要复制后返回，否则报错：unexpected fault address
	// @see https://github.com/boltdb/bolt/issues/204
	// ng := make([]byte, len(got))
	// copy(ng, got)

	return got, nil
}

// Put 向桶中存放数据
func Put(key []byte, value []byte, bucket []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.Put(key, value)
	})
}

// Del 删除桶中的数据
func Del(key []byte, bucket []byte) ([]byte, error) {
	var data []byte
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		data = b.Get(key)
		return b.Delete(key)
	})

	return data, err
}

// Query 查询指定桶内的所有数据
//
// param keySubStr 为需要包含的子字符串，当不为 nil 时，需要数据库中的键名为 string 类型;
// 当 keySubStr 为 nil 时，返回所有数据
func Query(bucket []byte, keySubStr *string) (map[string][]byte, error) {
	payload := make(map[string][]byte)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if keySubStr != nil && !strings.Contains(string(k), *keySubStr) {
				continue
			}
			payload[string(k)] = v
		}
		return nil
	})

	return payload, err
}
