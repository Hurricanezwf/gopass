package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Hurricanezwf/gopass/utils"
	"github.com/boltdb/bolt"
)

var (
	ErrNotExist = errors.New("Not existed")
)

var (
	defaultBucket = []byte("gopass")
)

type PasswdSVC struct {
	// 元数据存储文件
	metaFile string

	db *bolt.DB
}

func (p *PasswdSVC) Open(metaFile string) error {
	if len(metaFile) <= 0 {
		return errors.New("Empty metafile")
	}

	metaPath := filepath.Dir(metaFile)
	os.MkdirAll(metaPath, os.ModePerm)

	db, err := bolt.Open(metaFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("Open db(%s) failed, %v", metaFile, err)
	}

	p.db = db
	p.metaFile = metaFile
	return nil
}

func (p *PasswdSVC) Close() error {
	return p.db.Close()
}

// 密码加密后存储
func (p *PasswdSVC) Add(key, password []byte) error {
	// TODO: 加密方式可选，保存加密种类
	k := utils.Encrypt(key)
	pw := utils.Encrypt(password)

	if p.Exist(k) {
		return fmt.Errorf("Key(%s) had been existed", key)
	}

	err := p.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(defaultBucket)
		if err != nil {
			return err
		}
		if err = b.Put(k, pw); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err = p.db.Sync(); err != nil {
		return err
	}

	return nil
}

// 解密后返回
func (p *PasswdSVC) Get(key []byte) ([]byte, error) {
	var v []byte
	err := p.db.View(func(tx *bolt.Tx) error {
		v = tx.Bucket(defaultBucket).Get(key)
		if len(v) <= 0 {
			return ErrNotExist
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// TODO: 解密

	return v, nil
}

func (p *PasswdSVC) Update(key, password []byte) error {
	if !p.Exist(key) {
		return fmt.Errorf("Password for key(%s) is not existed", key)
	}

	// TODO: 加密

	err := p.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(defaultBucket)
		if err != nil {
			return err
		}
		if err = b.Put(key, password); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err = p.db.Sync(); err != nil {
		return err
	}

	return nil
}

func (p *PasswdSVC) Exist(key []byte) bool {
	if err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		if b == nil {
			return ErrNotExist
		}
		if v := b.Get(key); len(v) <= 0 {
			return ErrNotExist
		}
		return nil
	}); err != nil {
		return false
	}
	return true
}
