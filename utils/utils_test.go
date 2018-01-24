package utils

import (
	"bytes"
	"testing"

	"github.com/Hurricanezwf/gopass/utils"
)

func TestEncrypt(t *testing.T) {
	key := []byte("hello")
	ek := utils.Encrypt(key)

	dk, err := utils.Decrypt(ek)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if bytes.Equal(key, dk) == false {
		t.Fatalf("key(%x) != decryptKey(%x)", key, dk)
	}
}
