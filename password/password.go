package password

import (
	"CloudSync/log"
	"crypto/md5"
	"fmt"
	"io"
	"os"
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
	if key, err = encKey(sk, key); err != nil {
		return err
	}
	if password, err = encPassword(sk, password); err != nil {
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
	if key, err = encKey(sk, key); err != nil {
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
	if key, err = encKey(sk, key); err != nil {
		return err
	}
	if newPass, err = encPassword(sk, newPass); err != nil {
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
	if key, err = encKey(sk, key); err != nil {
		return nil, err
	}
	if m, err = meta.Open(makeMetaPath(), false); err != nil {
		return nil, err
	}
	defer m.Close()

	if p, err = m.Get(bucketDefault, key); err != nil {
		return nil, err
	}
	if p, err = decPassword(sk, p); err != nil {
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
		if k, err = decKey(sk, k); err != nil {
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

				oldBackupPath = filepath.Join(oldMetaPath, ".bak")
				oldBackupFile, err = os.OpenFile(oldBackupPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
				newMetaPath = filepath.Join(oldMetaPath, ".tmp")
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
				tmpNew := []byte(fmt.Sprintf("%x", md5.Sum(new)))
				if err = newMeta.Add(bucketAuth, []byte("sk"), tmpNew); err != nil {
					log.Error(err.Error())
					statemachine = Finish
					continue
				}

				// 更新所有数据
				var pairs []meta.Pair
				if pairs, err = oldMeta.ListAll(bucketDefault); err != nil {
					log.Error(err.Error())
					statemachine = Finish
					continue
				}
				for _, p := range pairs {
					if p.K, err = decKey(old, p.K); err != nil {
						log.Error(err.Error())
						break
					}
					if p.V, err = decPassword(old, p.V); err != nil {
						log.Error(err.Error())
						break
					}
					if p.K, err = encKey(new, p.K); err != nil {
						log.Error(err.Error())
						break
					}
					if p.V, err = encPassword(new, p.V); err != nil {
						log.Error(err.Error())
						break
					}
					if err = newMeta.Add(bucketDefault, p.K, p.V); err != nil {
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

func encKey(sk, key []byte) ([]byte, error) {
	return utils.Encrypt(sk, key, utils.TypeXORBase64)
}

func decKey(sk, key []byte) ([]byte, error) {
	return utils.Decrypt(sk, key)
}

func encPassword(sk, val []byte) ([]byte, error) {
	return utils.Encrypt(sk, val, utils.TypeCBC)
}

func decPassword(sk, val []byte) ([]byte, error) {
	return utils.Decrypt(sk, val)
}
