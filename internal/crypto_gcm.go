package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func GcmEncrypt(key, text []byte) ([]byte, error) {
	key32 := makeKeyLength32(key)
	c, err := aes.NewCipher(key32)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	result := gcm.Seal(nonce, nonce, text, nil)
	return result, nil
}

func Base64EncodeCipherText(cipherText []byte) string {
	return base64.StdEncoding.EncodeToString(cipherText)
}

func Base64DecodeCipherString(msg string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(msg)
	return cipherText, err
}

func GcmDecrypt(key, cipherText []byte) ([]byte, error) {
	key32 := makeKeyLength32(key)
	c, err := aes.NewCipher(key32)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("invalid cipherText")
	}

	nonce := cipherText[:nonceSize]

	cipherData := cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
