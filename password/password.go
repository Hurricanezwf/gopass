package password

import (
	"crypto/md5"
	"fmt"
	"path/filepath"

	"github.com/Hurricanezwf/gopass/g"
	"github.com/Hurricanezwf/gopass/meta"
	"github.com/Hurricanezwf/gopass/utils"
)

const MetaFile = "meta.db"

var (
	bucketDefault = []byte("gopass")
	bucketAuth    = []byte("auth_gopass")
)

func Add(sk, key, password []byte) error {
	var (
		err error
		m   *meta.Meta
	)
	if key, err = utils.Encrypt(sk, key, utils.TypeXORBase64); err != nil {
		return err
	}
	if password, err = utils.Encrypt(sk, password, utils.TypeCBC); err != nil {
		return err
	}
	if m, err = meta.Open(makeMetaPath(), false); err != nil {
		return err
	}
	defer m.Close()
	return m.Add(bucketDefault, key, password)
}

func Del(sk, key []byte) error {
	var (
		err error
		m   *meta.Meta
	)
	if key, err = utils.Encrypt(sk, key, utils.TypeXORBase64); err != nil {
		return err
	}
	if m, err = meta.Open(makeMetaPath(), false); err != nil {
		return err
	}
	defer m.Close()
	return m.Del(bucketDefault, key)
}

func Update(sk, key, newPass []byte) error {
	var (
		err error
		m   *meta.Meta
	)
	if key, err = utils.Encrypt(sk, key, utils.TypeXORBase64); err != nil {
		return err
	}
	if newPass, err = utils.Encrypt(sk, newPass, utils.TypeCBC); err != nil {
		return err
	}
	if m, err = meta.Open(makeMetaPath(), false); err != nil {
		return err
	}
	defer m.Close()
	return m.Update(bucketDefault, key, newPass)
}

func Get(sk, key []byte) ([]byte, error) {
	var (
		err error
		p   []byte
		m   *meta.Meta
	)
	if key, err = utils.Encrypt(sk, key, utils.TypeXORBase64); err != nil {
		return nil, err
	}
	if m, err = meta.Open(makeMetaPath(), false); err != nil {
		return nil, err
	}
	defer m.Close()

	if p, err = m.Get(bucketDefault, key); err != nil {
		return nil, err
	}
	if p, err = utils.Decrypt(sk, p); err != nil {
		return nil, err
	}
	return p, nil
}

func ListKeys(sk []byte) ([][]byte, error) {
	var (
		err  error
		keys [][]byte
		m    *meta.Meta
	)
	if m, err = meta.Open(makeMetaPath(), false); err != nil {
		return nil, err
	}
	defer m.Close()

	if keys, err = m.ListKeys(bucketDefault); err != nil {
		return nil, err
	}
	for i, k := range keys {
		if k, err = utils.Decrypt(sk, k); err != nil {
			return nil, err
		}
		keys[i] = k
	}
	return keys, nil
}

func GetAuthSK() ([]byte, error) {
	m, err := meta.Open(makeMetaPath(), false)
	if err != nil {
		return nil, err
	}
	defer m.Close()
	return m.Get(bucketAuth, []byte("sk"))
}

func InitAuthSK(sk []byte) ([]byte, error) {
	m, err := meta.Open(makeMetaPath(), false)
	if err != nil {
		return nil, err
	}
	defer m.Close()

	sk = []byte(fmt.Sprintf("%x", md5.Sum(sk)))
	err = m.Add(bucketAuth, []byte("sk"), sk)
	return sk, err
}

func makeMetaPath() string {
	return filepath.Join(g.MetaDir, MetaFile)
}
