package qfcrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// PBKDF2 KDF计算key
func PBKDF2(masterKey []byte, salt []byte) []byte {

	hashFunc := func() hash.Hash {
		// 采用hmac_sha256加密
		return sha256.New()
	}
	secret := pbkdf2.Key(masterKey, salt, 10000, 32, hashFunc)

	return secret
}

// GenerateKey 生成无规则key的函数，用于为加密产生一个master key
func GenerateKey(origin []byte) (string, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	key := base64.StdEncoding.EncodeToString(PBKDF2(origin, iv))
	return key, nil
}

// GenerateIV iv生成
func GenerateIV(len int) ([]byte, error) {
	iv := make([]byte, len)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	return iv, nil
}

// -------------------此部分接口在原始版本中涉及
const (
	BYTES_16 = 16
	BYTES_24 = 24
	BYTES_32 = 32
)

func ZorePadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	return append(ciphertext, padtext...)
}

func GetFormatKey(key []byte) []byte {
	if len(key) <= BYTES_16 {
		key = ZorePadding(key, BYTES_16)
	} else if len(key) <= BYTES_24 {
		key = ZorePadding(key, BYTES_24)
	} else {
		key = ZorePadding(key, BYTES_32)
	}
	return key
}

// -------------------此部分接口在原始版本中涉及
