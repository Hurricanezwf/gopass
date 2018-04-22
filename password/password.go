package password

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Hurricanezwf/gopass/crypt"
	"github.com/Hurricanezwf/gopass/g"
	"github.com/Hurricanezwf/gopass/log"
	"github.com/Hurricanezwf/gopass/meta"
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
	if key, err = crypt.EncKey(sk, key); err != nil {
		return err
	}
	if password, err = crypt.EncPassword(sk, password); err != nil {
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
	if key, err = crypt.EncKey(sk, key); err != nil {
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
	if key, err = crypt.EncKey(sk, key); err != nil {
		return err
	}
	if newPass, err = crypt.EncPassword(sk, newPass); err != nil {
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
	if key, err = crypt.EncKey(sk, key); err != nil {
		return nil, err
	}
	if m, err = meta.Open(makeMetaPath(), false); err != nil {
		return nil, err
	}
	defer m.Close()

	if p, err = m.Get(bucketDefault, key); err != nil {
		return nil, err
	}
	if p, err = crypt.DecPassword(sk, p); err != nil {
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
		if k, err = crypt.DecKey(sk, k); err != nil {
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

	skEnc := crypt.EncryptSK(sk)
	err = m.Add(bucketAuth, []byte("sk"), skEnc)
	return skEnc, err
}

func ChangeSK(old, new []byte) error {
	var (
		err           error
		oldMetaPath   string
		newMetaPath   string
		oldBackupPath string
		oldMetaFile   *os.File   // to close
		oldBackupFile *os.File   // to close
		oldMeta       *meta.Meta // to close
		newMeta       *meta.Meta // to close
		loop          = true
	)

	// state machine
	const (
		BackupOldMeta = iota
		BuildNewMeta
		Finish
		RollBack
	)

	statemachine := BackupOldMeta
	for loop {
		switch statemachine {
		case BackupOldMeta:
			{
				oldMetaPath = makeMetaPath()
				oldMetaFile, err = os.Open(oldMetaPath)
				if err != nil {
					log.Error(err.Error())
					statemachine = Finish
					continue
				}

				oldBackupPath = oldMetaPath + ".bak"
				oldBackupFile, err = os.OpenFile(oldBackupPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
					log.Error(err.Error())
					statemachine = Finish
					continue
				}

				if _, err = io.Copy(oldBackupFile, oldMetaFile); err != nil {
					log.Error(err.Error())
					statemachine = Finish
					continue
				}

				statemachine = BuildNewMeta // next state
			}
		case BuildNewMeta:
			{
				newMetaPath = oldMetaPath + ".tmp"
				os.Remove(newMetaPath)
				newMeta, err = meta.Open(newMetaPath, false)
				if err != nil {
					log.Error(err.Error())
					statemachine = Finish
					continue
				}

				oldMeta, err = meta.Open(oldMetaPath, true)
				if err != nil {
					log.Error(err.Error())
					statemachine = Finish
					continue
				}

				// 写入新SK
				tmpNew := crypt.EncryptSK(new)
				if err = newMeta.Add(bucketAuth, []byte("sk"), tmpNew); err != nil {
					log.Error(err.Error())
					statemachine = Finish
					continue
				}

				// 更新所有数据
				var (
					decOldKey, decOldVal []byte
					encNewKey, encNewVal []byte
					pairs                []meta.Pair
				)
				if pairs, err = oldMeta.ListAll(bucketDefault); err != nil {
					log.Error(err.Error())
					statemachine = Finish
					continue
				}
				for _, p := range pairs {
					if len(p.K) <= 0 || len(p.V) <= 0 {
						continue
					}
					if decOldKey, err = crypt.DecKey(old, p.K); err != nil {
						log.Error(err.Error())
						break
					}
					if decOldVal, err = crypt.DecPassword(old, p.V); err != nil {
						log.Error(err.Error())
						break
					}
					if encNewKey, err = crypt.EncKey(new, decOldKey); err != nil {
						log.Error(err.Error())
						break
					}
					if encNewVal, err = crypt.EncPassword(new, decOldVal); err != nil {
						log.Error(err.Error())
						break
					}
					if err = newMeta.Add(bucketDefault, encNewKey, encNewVal); err != nil {
						log.Error(err.Error())
						break
					}
				}
				if err != nil {
					statemachine = Finish
					continue
				}

				statemachine = Finish // next state
				continue
			}
		case Finish:
			{
				// 关闭资源
				if oldMeta != nil {
					oldMeta.Close()
					oldMeta = nil
				}
				if newMeta != nil {
					newMeta.Close()
					newMeta = nil
				}
				if oldBackupFile != nil {
					oldBackupFile.Close()
					oldBackupFile = nil
				}
				if oldMetaFile != nil {
					oldMetaFile.Close()
					oldMetaFile = nil
				}

				if err != nil {
					statemachine = RollBack
					continue
				}

				// 替换元数据
				if err = os.Rename(newMetaPath, oldMetaPath); err != nil {
					log.Error(err.Error())
					statemachine = RollBack
					continue
				}
				os.Remove(newMetaPath)
				loop = false
			}
		case RollBack:
			os.Remove(oldBackupPath)
			os.Remove(newMetaPath)
			loop = false
		default:
			err = fmt.Errorf("Unknown state %d", statemachine)
			loop = false
		}
	}
	return err
}

func makeMetaPath() string {
	return filepath.Join(g.MetaDir, MetaFile)
}
