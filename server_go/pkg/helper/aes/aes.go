package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	tostring "server_go/pkg/helper/toString"
	"strconv"
	"strings"

	"errors"
)

// 生成一个固定长度的随机密钥
//
//	func generateKey() ([]byte, error) {
//		key := make([]byte, 16) // AES-128, 也可以选择 24 或 32 字节用于 AES-192 或 AES-256
//		if _, err := io.ReadFull(rand.Reader, key); err != nil {
//			return nil, err
//		}
//		return key, nil
//	}

// 根据自定义字符串生成随机aes密钥
func GenerateKey(customString string) []byte {
	hash := sha256.Sum256([]byte(customString))
	return hash[:]
}

// 使用 AES 加密数据
func Encrypt(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize] // 初始向量
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 使用 AES 解密数据
func Decrypt(key []byte, encodedString string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedString)

	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	ciphertext, err = pkcs7Unpad(ciphertext)
	return tostring.Strval(ciphertext), err
}

func StringToByte(keyStr string) ([]byte, error) {
	numStrs := strings.Fields(keyStr) // Use Fields to automatically handle whitespace
	var key []byte

	for _, numStr := range numStrs {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, err

		}
		key = append(key, byte(num))
	}

	return key, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS#7填充
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("data is empty")
	}
	padding := int(data[length-1])
	if padding > length || padding > 16 {
		return nil, errors.New("invalid padding size")
	}
	return data[:length-padding], nil
}
