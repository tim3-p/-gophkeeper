package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"io"

	"github.com/tim3-p/gophkeeper/internal/common"
)

// MakeKey makes interanl key from any string
func MakeKey(s string) common.Key {
	return common.Key(sha256.Sum256([]byte(s)))
}

// EncryptString is the same as Encrypt but for strings
func EncryptString(key common.Key, clearText string) (string, error) {
	buf, err := Encrypt(key, []byte(clearText))
	return hex.EncodeToString(buf), err
}

// DecryptString is the same as Encrypt but for strings
func DecryptString(key common.Key, cypherText string) (string, error) {
	tmp, err := hex.DecodeString(cypherText)
	if err != nil {
		return "", err
	}
	buf, err := Decrypt(key, tmp)
	return string(buf), err
}

// Encrypt encrypts the cleartext with the key given
func Encrypt(key common.Key, clearText []byte) ([]byte, error) {
	c, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nonce, nonce, clearText, nil)
	return cipherText, nil
}

// Decrypt decrypts the ciphertext with the key given
func Decrypt(key common.Key, cipherText []byte) ([]byte, error) {
	c, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("cipher text is too short")
	}

	nonce, text := cipherText[:nonceSize], cipherText[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, text, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
