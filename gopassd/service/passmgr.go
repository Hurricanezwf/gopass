package service

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Hurricanezwf/gopass/gopassd/g"
	"github.com/boltdb/bolt"
)

var (
	ErrKeyExisted = errors.New("key existed")
)

type Password struct {
	Package string
	Key     string
	Value   string
}

type PassMgr struct {
	mgrMutex sync.RWMutex

	dbMetaPath string
	db         *bolt.DB

	//pwdByKey map[string]*Password // map[key]*Password
	//pwdByPkg map[string]*Password // map[key]*Password

	//buf *bytes.Buffer
}

func (m *PassMgr) Open(metaPath string) error {
	m.dbMetaPath = metaPath
	//m.pwdByKey = make(map[string]*Password)
	//m.pwdByPkg = make(map[string]*Password)
	//m.buf = bytes.NewBuffer(nil)

	var err error

	dir := filepath.Dir(m.dbMetaPath)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		g.Log.Error("MkdirAll %s failed, %v", dir, err)
		return err
	}

	if m.db, err = bolt.Open(m.dbMetaPath, 0644, &bolt.Options{
		Timeout: time.Minute,
	}); err != nil {
		g.Log.Error("Open db %s failed, %v", m.dbMetaPath, err)
		return err
	}

	return nil
}

func (m *PassMgr) Close() error {
	var err error
	if err = m.db.Close(); err != nil {
		g.Log.Error("Close DB failed, %v", err)
		return err
	}
	return nil
}

func (m *PassMgr) Add(bucket, key, value string) error {
	var (
		err error
		b   *bolt.Bucket
		//pkgNameB64 = utils.Base64([]byte(pkgName))
		//keyB64     = utils.Base64([]byte(key))
		//valueB64   = utils.Base64([]byte(value))
	)

	if err = m.db.Update(func(tx *bolt.Tx) error {
		if b, err = tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
			return err
		}
		if err = b.Put([]byte(key), []byte(value)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		g.Log.Error("Add Key(%s) to bucket(%s) failed, %v", key, bucket, err)
		return err
	}
	return nil
}
