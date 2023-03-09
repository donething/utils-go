// Package dobolt boltdb 的帮助函数
//
// @see https://github.com/etcd-io/bbolt
package dobolt

import (
	bolt "go.etcd.io/bbolt"
	"os"
	"strings"
	"time"
)

// DoBolt boltdb 的包装。执行数据操作后不会关闭数据库，可按需调用`db.DB.Close()`关闭。还可以在程序退出时，才关闭数据库
type DoBolt struct {
	DB *bolt.DB
}

// Open 根据数据库路径打开数据库。当 `mode`、`options` 为 `nil` 时，将使用默认值
func Open(dbPath string, mode *os.FileMode, options *bolt.Options) (*DoBolt, error) {
	// 设置选项
	var md = mode
	if md == nil {
		// 分配内存，避免赋值时发生空指针异常
		md = new(os.FileMode)
		*md = 0600
	}

	var opts = options
	if opts == nil {
		opts = &bolt.Options{Timeout: 3 * time.Second}
	}

	// 打开数据库
	db, err := bolt.Open(dbPath, *md, opts)
	return &DoBolt{DB: db}, err
}

// Close 关闭数据库。请仅用此函数关闭，不要再其它地方调用 db.DB.Close() 来关闭数据库
func (db *DoBolt) Close() error {
	if db.DB != nil {
		err := db.DB.Close()
		db.DB = nil
		return err
	}

	return nil
}

// Create 创建桶
func (db *DoBolt) Create(bucket []byte) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		_, errC := tx.CreateBucketIfNotExists(bucket)
		return errC
	})
}

// Get 获取桶中键对应的值
func (db *DoBolt) Get(key []byte, bucket []byte) ([]byte, error) {
	var value []byte
	err := db.DB.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(bucket).Get(key)
		// 获取的值只在事务内有效，之外需要复制后返回，否则可能报错：unexpected fault address
		// @see https://github.com/boltdb/bolt/issues/204
		// got := make([]byte, len(got))
		// copy(got, value)
		return nil
	})

	return value, err
}

// Put 向桶中存放数据
func (db *DoBolt) Put(key []byte, value []byte, bucket []byte) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.Put(key, value)
	})
}

// Del 删除桶中的数据
func (db *DoBolt) Del(key []byte, bucket []byte) ([]byte, error) {
	var data []byte
	err := db.DB.Update(func(tx *bolt.Tx) error {
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
func (db *DoBolt) Query(keySubStr *string, bucket []byte) (map[string][]byte, error) {
	payload := make(map[string][]byte)
	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if keySubStr != nil && !strings.Contains(strings.ToLower(string(k)), strings.ToLower(*keySubStr)) {
				continue
			}
			payload[string(k)] = v
		}
		return nil
	})

	return payload, err
}

// Batch 批量插入数据
func (db *DoBolt) Batch(data map[string][]byte, bucket []byte) error {
	err := db.DB.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		for key, bs := range data {
			errPut := b.Put([]byte(key), bs)
			if errPut != nil {
				return errPut
			}
		}
		return nil
	})

	return err
}
