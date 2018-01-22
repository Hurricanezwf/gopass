package service

import (
	"github.com/Hurricanezwf/gopass/log"
)

var (
	MetaFile string = "./meta/meta.db"
)

var (
	Passwd *PasswdSVC
)

func Open() error {
	if Passwd == nil {
		Passwd = &PasswdSVC{}
	}

	var err error
	if err = Passwd.Open(MetaFile); err != nil {
		log.Error("Open passwd service failed, %v", err)
		return err
	}

	return nil
}

func Close() error {
	if Passwd != nil {
		Passwd.Close()
		Passwd = nil
	}
	return nil
}
