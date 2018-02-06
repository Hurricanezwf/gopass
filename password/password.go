package password

import (
	"path/filepath"

	"github.com/Hurricanezwf/gopass/g"
)

const MetaFile = "meta.db"

func Add(key, password []byte) error {
	p := &PasswdProxy{}
	if err := p.Open(makeMetaPath(), false); err != nil {
		return err
	}
	defer p.Close()
	return p.Add(key, password)
}

func Del(key []byte) error {
	p := &PasswdProxy{}
	if err := p.Open(makeMetaPath(), false); err != nil {
		return err
	}
	defer p.Close()
	return p.Del(key)
}

func Update(key, newPass []byte) error {
	p := &PasswdProxy{}
	if err := p.Open(makeMetaPath(), false); err != nil {
		return err
	}
	defer p.Close()
	return p.Update(key, newPass)
}

func Get(key []byte) ([]byte, error) {
	p := &PasswdProxy{}
	if err := p.Open(makeMetaPath(), true); err != nil {
		return nil, err
	}
	defer p.Close()
	return p.Get(key)
}

func ListKeys() ([][]byte, error) {
	p := &PasswdProxy{}
	if err := p.Open(makeMetaPath(), true); err != nil {
		return nil, err
	}
	defer p.Close()
	return p.ListKeys()
}

func makeMetaPath() string {
	return filepath.Join(g.MetaDir, MetaFile)
}
