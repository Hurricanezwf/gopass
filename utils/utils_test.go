package utils

import (
	"bytes"
	"testing"
)

func TestCBCEncrypt(t *testing.T) {
	key := []byte("zwf")
	toEncrypt := []byte("hello")
	ek, err := Encrypt(key, toEncrypt, TypeCBC)
	if err != nil {
		t.Fatalf("%v", err)
	}

	dk, err := Decrypt(key, ek)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if bytes.Equal(toEncrypt, dk) == false {
		t.Fatalf("origion(%x) != decrypt(%x)", toEncrypt, dk)
	}
}

func TestXORBase64Encrypt(t *testing.T) {
	key := []byte("zwf")
	toDecrypt := []byte("hello")
	ek, err := Encrypt(key, toDecrypt, TypeXORBase64)
	if err != nil {
		t.Fatalf("%v", err)
	}

	dk, err := Decrypt(key, ek)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if bytes.Equal(toDecrypt, dk) == false {
		t.Fatalf("origion(%s) != decrypt(%s)", string(toDecrypt), string(dk))
	}
}
