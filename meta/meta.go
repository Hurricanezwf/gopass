package meta

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Hurricanezwf/gopass/log"
	"github.com/boltdb/bolt"
)

var (
	ErrNotExist = errors.New("Doesn't exist")
	ErrExisted  = errors.New("Existed")
)

type Pair struct {
	K []byte
	V []byte
}

type Meta struct {
	// 元数据存储文件
	metaFile string

	db *bolt.DB
}

func Open(metaFile string, readOnly bool) (*Meta, error) {
	if len(metaFile) <= 0 {
		return nil, errors.New("Empty metafile")
	}

	metaDir := filepath.Dir(metaFile)
	os.MkdirAll(metaDir, os.ModePerm)
	log.Debug("Create dir %s", metaDir)

	db, err := bolt.Open(metaFile, 0600,
		&bolt.Options{
			ReadOnly: readOnly,
			Timeout:  1 * time.Second,
		})
	if err != nil {
		return nil, fmt.Errorf("Open db(%s) failed, %v", metaFile, err)
	}
	m := &Meta{
		metaFile: metaFile,
		db:       db,
	}
	return m, nil
}

func (m *Meta) Close() error {
	if m.db != nil {
		tmpDB := m.db
		m.db = nil
		return tmpDB.Close()
	}
	return nil
}

func (m *Meta) Add(bucket, k, v []byte) error {
	if m.Exist(bucket, k) {
		return ErrExisted
	}
	err := m.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		if err = b.Put(k, v); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err = m.db.Sync(); err != nil {
		return err
	}
	return nil
}

func (m *Meta) Del(bucket, k []byte) error {
	if !m.Exist(bucket, k) {
		return ErrNotExist
	}
	err := m.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		if err = b.Delete(k); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err = m.db.Sync(); err != nil {
		return err
	}
	return nil
}

func (m *Meta) Get(bucket, k []byte) ([]byte, error) {
	var v []byte
	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrNotExist
		}
		if tmp := b.Get(k); len(tmp) <= 0 {
			return ErrNotExist
		} else {
			v = make([]byte, len(tmp))
			copy(v, tmp)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (m *Meta) ListKeys(bucket []byte) ([][]byte, error) {
	keys := make([][]byte, 0)
	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b != nil {
			err := b.ForEach(func(k, v []byte) error {
				keys = append(keys, k)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (m *Meta) ListAll(bucket []byte) ([]Pair, error) {
	pairs := make([]Pair, 0, 32)
	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b != nil {
			err := b.ForEach(func(k, v []byte) error {
				pairs = append(pairs, Pair{K: k, V: v})
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return pairs[:len(pairs)], nil
}

func (m *Meta) Update(bucket, k, v []byte) error {
	if !m.Exist(bucket, k) {
		return ErrNotExist
	}
	err := m.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		if err = b.Put(k, v); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err = m.db.Sync(); err != nil {
		return err
	}
	return nil
}

func (m *Meta) Exist(bucket, k []byte) bool {
	if err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrNotExist
		}
		if v := b.Get(k); len(v) <= 0 {
			return ErrNotExist
		}
		return nil
	}); err != nil {
		return false
	}
	return true
}

// 密码加密后存储
//func (p *MetaProxy) Add(key, password []byte) error {
//	// TODO: 加密方式可选，保存加密种类
//	k := utils.Encrypt(key)
//	pw := utils.Encrypt(password)
//
//	if p.Exist(k) {
//		return fmt.Errorf("Key(%s) had been existed", key)
//	}
//
//	err := p.db.Update(func(tx *bolt.Tx) error {
//		b, err := tx.CreateBucketIfNotExists(defaultBucket)
//		if err != nil {
//			return err
//		}
//		if err = b.Put(k, pw); err != nil {
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//	if err = p.db.Sync(); err != nil {
//		return err
//	}
//
//	return nil
//}

// 解密后返回
//func (p *MetaProxy) Get(key []byte) ([]byte, error) {
//	var (
//		err error
//		v   []byte
//		k   = utils.Encrypt(key)
//	)
//
//	if err = p.db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket(defaultBucket)
//		if b == nil {
//			return ErrNotExist
//		}
//		if v = b.Get(k); len(v) <= 0 {
//			return ErrNotExist
//		}
//		return nil
//	}); err != nil {
//		return nil, err
//	}
//	if len(v) <= 0 {
//		return nil, fmt.Errorf("Password for key(%s) is empty", string(key))
//	}
//	if v, err = utils.Decrypt(v); err != nil {
//		return nil, fmt.Errorf("Decrypt password failed, %v", err)
//	}
//	return v, nil
//}
//
//func (p *MetaProxy) ListKeys() ([][]byte, error) {
//	keys := make([][]byte, 0)
//	if err := p.db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket(defaultBucket)
//		if b != nil {
//			err := b.ForEach(func(k, v []byte) error {
//				if dk, e := utils.Decrypt(k); e == nil && len(dk) > 0 {
//					keys = append(keys, dk)
//				}
//				return nil
//			})
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	}); err != nil {
//		return nil, err
//	}
//	return keys, nil
//}

//func (p *MetaProxy) Update(key, new []byte) error {
//	k := utils.Encrypt(key)
//	newpw := utils.Encrypt(new)
//
//	if !p.Exist(k) {
//		return fmt.Errorf("Key(%s) doesn't exist", key)
//	}
//
//	err := p.db.Update(func(tx *bolt.Tx) error {
//		b, err := tx.CreateBucketIfNotExists(defaultBucket)
//		if err != nil {
//			return err
//		}
//		if err = b.Put(k, newpw); err != nil {
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//	if err = p.db.Sync(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (p *MetaProxy) Del(key []byte) error {
//	k := utils.Encrypt(key)
//
//	if !p.Exist(k) {
//		return fmt.Errorf("Key(%s) doesn't exist", key)
//	}
//
//	err := p.db.Update(func(tx *bolt.Tx) error {
//		b, err := tx.CreateBucketIfNotExists(defaultBucket)
//		if err != nil {
//			return err
//		}
//		if err = b.Delete(k); err != nil {
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//	if err = p.db.Sync(); err != nil {
//		return err
//	}
//
//	return nil
//}

//func (p *MetaProxy) Exist(key []byte) bool {
//	if err := p.db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket(defaultBucket)
//		if b == nil {
//			return ErrNotExist
//		}
//		if v := b.Get(key); len(v) <= 0 {
//			return ErrNotExist
//		}
//		return nil
//	}); err != nil {
//		return false
//	}
//	return true
//}
