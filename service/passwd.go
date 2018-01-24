package service

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/Hurricanezwf/gopass/log"
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

	user, err := user.Current()
	if err != nil || user == nil {
		return fmt.Errorf("Get current user failed, %v", err)
	}
	metaDir := filepath.Join(user.HomeDir, ".gopass")
	metaFile = filepath.Join(metaDir, metaFile)
	os.MkdirAll(metaDir, os.ModePerm)
	log.Debug("Create dir %s", metaDir)

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
	var (
		err error
		v   []byte
		k   = utils.Encrypt(key)
	)

	if err = p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		if b == nil {
			return ErrNotExist
		}
		if v = b.Get(k); len(v) <= 0 {
			return ErrNotExist
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(v) <= 0 {
		return nil, fmt.Errorf("Password for key(%s) is empty", string(key))
	}
	if v, err = utils.Decrypt(v); err != nil {
		return nil, fmt.Errorf("Decrypt password failed, %v", err)
	}
	return v, nil
}

func (p *PasswdSVC) ListKeys() ([][]byte, error) {
	keys := make([][]byte, 0)
	if err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		if b != nil {
			err := b.ForEach(func(k, v []byte) error {
				if dk, e := utils.Decrypt(k); e == nil && len(dk) > 0 {
					keys = append(keys, dk)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return keys, nil
}

func (p *PasswdSVC) Update(key, new []byte) error {
	k := utils.Encrypt(key)
	newpw := utils.Encrypt(new)

	if !p.Exist(k) {
		return fmt.Errorf("Key(%s) doesn't exist", key)
	}

	err := p.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(defaultBucket)
		if err != nil {
			return err
		}
		if err = b.Put(k, newpw); err != nil {
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

func (p *PasswdSVC) Del(key []byte) error {
	k := utils.Encrypt(key)

	if !p.Exist(k) {
		return fmt.Errorf("Key(%s) doesn't exist", key)
	}

	err := p.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(defaultBucket)
		if err != nil {
			return err
		}
		if err = b.Delete(key); err != nil {
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
