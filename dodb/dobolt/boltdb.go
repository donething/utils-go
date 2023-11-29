// Package dobolt boltdb 的帮助函数
//
// 读写前，务必先创建桶，否则会报空指针错误
//
// @see https://github.com/etcd-io/bbolt
package dobolt

import (
	"bytes"
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

// Get 获取桶中键对应的值。不存在该键时返回`nil`
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

// Set 向桶中存放数据
func (db *DoBolt) Set(key []byte, value []byte, bucket []byte) error {
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

// DelPrefixAll 删除桶中的为指定前缀的所有数据
func (db *DoBolt) DelPrefixAll(prefix []byte, bucket []byte) error {
	// 在事务中执行删除操作
	err := db.DB.Update(func(tx *bolt.Tx) error {
		// 找到要操作的 bucket
		b := tx.Bucket(bucket)
		if b == nil {
			// 如果bucket不存在，直接返回
			return bolt.ErrBucketNotFound
		}

		// 遍历 bucket
		c := b.Cursor()
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			// 删除匹配前缀的键
			if err := b.Delete(k); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// Query 查询指定桶内的所有数据
//
// 参数 keySub 需要包含的的子串。当为空时，返回所有数据
func (db *DoBolt) Query(keySub string, bucket []byte) (map[string][]byte, error) {
	payload := make(map[string][]byte)
	err := db.DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucket).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if keySub != "" && !strings.Contains(strings.ToUpper(string(k)), strings.ToUpper(keySub)) {
				continue
			}

			payload[string(k)] = v
		}

		return nil
	})

	return payload, err
}

// QueryPrefix 指定桶内的前缀扫描
//
// 参数 prefix 为需要扫描的前缀
func (db *DoBolt) QueryPrefix(prefix []byte, bucket []byte) (map[string][]byte, error) {
	payload := make(map[string][]byte)
	err := db.DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucket).Cursor()

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
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

// Clear 清空桶
func (db *DoBolt) Clear(bucket []byte) error {
	// 获取事务
	err := db.DB.Update(func(tx *bolt.Tx) error {
		// 直接删除桶
		if err := tx.DeleteBucket(bucket); err != nil {
			return err
		}

		// 重新创建同名的桶
		_, err := tx.CreateBucket(bucket)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
