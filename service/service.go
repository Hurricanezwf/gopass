package service

import (
	"github.com/Hurricanezwf/gopass/log"
)

var (
	MetaFile string = "./meta/meta.db"
)

var (
	passwd *Passwd
)

func Open() error {
	if passwd == nil {
		passwd = &Passwd{}
	}

	var err error
	if err = passwd.Open(MetaFile); err != nil {
		log.Error("Open passwd failed, %v", err)
		return err
	}

	return nil
}

func Close() error {
	if passwd != nil {
		passwd.Close()
		passwd = nil
	}
	return nil
}
