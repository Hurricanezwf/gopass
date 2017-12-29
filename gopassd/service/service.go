package service

import (
	"path/filepath"

	"github.com/Hurricanezwf/gopass/gopassd/g"
)

var (
	PasswordMgr PassMgr
)

func Run() error {
	var err error

	// 加载元数据
	p := filepath.Join(g.MetaDir, g.MetaName)
	if err = PasswordMgr.Open(p); err != nil {
		g.Log.Error("Service: Load password failed, %v", err)
		return err
	}
	return nil
}

func Close() error {
	PasswordMgr.Close()
	return nil
}
