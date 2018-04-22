package crypt

import (
	"crypto/md5"
	"fmt"

	"github.com/Hurricanezwf/gopass/utils"
)

func EncKey(sk, key []byte) ([]byte, error) {
	return utils.Encrypt(sk, key, utils.TypeXORBase64)
}

func DecKey(sk, key []byte) ([]byte, error) {
	return utils.Decrypt(sk, key)
}

func EncPassword(sk, val []byte) ([]byte, error) {
	return utils.Encrypt(sk, val, utils.TypeAES)
}

func DecPassword(sk, val []byte) ([]byte, error) {
	return utils.Decrypt(sk, val)
}

func EncryptSK(sk []byte) []byte {
	return []byte(fmt.Sprintf("%x", md5.Sum(sk)))
}
