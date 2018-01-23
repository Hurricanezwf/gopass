package service

import (
	"sync"
)

var (
	MetaFile string = "./meta/meta.db"
)

var (
	passwdMutex sync.RWMutex
	passwd      *PasswdSVC
)

//func Open() error {
//	if Passwd == nil {
//		Passwd = &PasswdSVC{}
//	}
//
//	var err error
//	if err = Passwd.Open(MetaFile); err != nil {
//		log.Error("Open passwd service failed, %v", err)
//		return err
//	}
//
//	return nil
//}

func CloseAll() error {
	passwdMutex.Lock()
	if passwd != nil {
		passwd.Close()
		passwd = nil
	}
	passwdMutex.Unlock()
	return nil
}

func AddPassword(key, password []byte) error {
	passwdMutex.Lock()
	defer passwdMutex.Unlock()

	if passwd == nil {
		passwd = &PasswdSVC{}
		if err := passwd.Open(MetaFile); err != nil {
			return err
		}
	}

	return passwd.Add(key, password)
}

func DelPassword(key []byte) error {
	passwdMutex.Lock()
	defer passwdMutex.Unlock()

	if passwd == nil {
		passwd = &PasswdSVC{}
		if err := passwd.Open(MetaFile); err != nil {
			return err
		}
	}

	return passwd.Del(key)
}

func UpdatePassword(key, newPass []byte) error {
	passwdMutex.Lock()
	defer passwdMutex.Unlock()

	if passwd == nil {
		passwd = &PasswdSVC{}
		if err := passwd.Open(MetaFile); err != nil {
			return err
		}
	}

	return passwd.Update(key, newPass)
}

func GetPassword(key []byte) ([]byte, error) {
	passwdMutex.RLock()
	defer passwdMutex.RUnlock()

	if passwd == nil {
		passwd = &PasswdSVC{}
		if err := passwd.Open(MetaFile); err != nil {
			return nil, err
		}
	}

	return passwd.Get(key)
}

func ListKeys() ([]string, error) {
	passwdMutex.RLock()
	defer passwdMutex.RUnlock()

	if passwd == nil {
		passwd = &PasswdSVC{}
		if err := passwd.Open(MetaFile); err != nil {
			return nil, err
		}
	}

	return passwd.ListKeys()
}
