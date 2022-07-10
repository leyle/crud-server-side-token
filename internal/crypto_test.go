package internal

import (
	"testing"
)

var key = []byte("v}#0CEYuG%M8#c77HkUp")
var msg = []byte("eventservice")

func encrypt() []byte {
	cipherText, err := GcmEncrypt(key, msg)
	if err != nil {
		panic(err)
	}

	return cipherText
}

func TestGCMEncrypt(t *testing.T) {
	cipherText := encrypt()

	t.Log(cipherText)
	t.Log(string(cipherText))

	t.Log(Base64EncodeCipherText(cipherText))
}

func TestGcmDecrypt(t *testing.T) {
	cipherText := encrypt()
	t.Log(Base64EncodeCipherText(cipherText))

	text, err := GcmDecrypt(key, cipherText)
	if err != nil {
		t.Error(err)
	}

	t.Log(text)
	t.Log(string(text))
}

func TestGCMDecryptSRCBase64(t *testing.T) {
	// raw := "B2h+HK3at6Yo0pkWbsTgMRx6gy7es9u70v9hVrvqOa/q85pk9jixgA=="
	raw := "BLu20Mtix1bfceKx+TrCLM7oFEtAQK5JVAiC4q+kwyZX0eysBH8JwA=="
	decodeBytes, err := Base64DecodeCipherString(raw)
	if err != nil {
		t.Error(err)
	}

	text, err := GcmDecrypt(key, decodeBytes)
	if err != nil {
		t.Error(err)
	}
	t.Log(text)
	t.Log(string(text))
}

func TestEncryptAndDecrypt(t *testing.T) {
	cipherText := encrypt()
	b64Str := Base64EncodeCipherText(cipherText)

	t.Log(b64Str)

	srcCipherText, err := Base64DecodeCipherString(b64Str)
	if err != nil {
		panic(err)
	}
	text, err := GcmDecrypt(key, srcCipherText)
	if err != nil {
		panic(err)
	}

	t.Log(string(text))
}
