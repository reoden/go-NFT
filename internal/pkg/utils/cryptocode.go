package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

func Encrypt(plaintext string) (string, error) {
	aesKey := os.Getenv("AES_KEY")
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string) (string, error) {
	aesKey := os.Getenv("AES_KEY")
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce := data[:gcm.NonceSize()]
	encrypted := data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

type AESCryptoString string

func (a *AESCryptoString) Value() (driver.Value, error) {
	if *a == "" {
		return "", nil
	}
	return Encrypt(string(*a))
}

func (a *AESCryptoString) Scan(value interface{}) error {
	if value == nil {
		*a = ""
		return nil
	}

	var encrypted string

	switch v := value.(type) {
	case string:
		encrypted = v
	case []byte:
		encrypted = string(v)
	default:
		return fmt.Errorf("AesString: unsupported type %T", value)
	}

	if encrypted == "" {
		*a = ""
		return nil
	}
	plain, err := Decrypt(encrypted)
	if err != nil {
		return err
	}

	*a = AESCryptoString(plain)
	return nil
}
