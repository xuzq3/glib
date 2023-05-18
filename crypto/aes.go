package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type AesCipher struct {
	key       []byte
	block     cipher.Block
	blockSize int
}

func NewAesCipher(key []byte) (*AesCipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	a := &AesCipher{
		key:       key,
		block:     block,
		blockSize: block.BlockSize(),
	}
	return a, nil
}

func (a *AesCipher) EncryptToString(data string) (string, error) {
	// 转成字节数组
	origData := []byte(data)
	// 补全码
	origData = pkcs7Padding(origData, a.block.BlockSize())
	// 创建数组
	encrypted := make([]byte, len(origData))
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(a.block, a.key[:a.blockSize])
	// 加密
	blockMode.CryptBlocks(encrypted, origData)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (a *AesCipher) Encrypt(plainText []byte) ([]byte, error) {
	encrypted, err := a.EncryptToString(string(plainText))
	if err != nil {
		return nil, err
	}
	return []byte(encrypted), nil
}

func (a *AesCipher) DecryptToString(encrypted string) (string, error) {
	encryptedByte, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	if len(encryptedByte)%a.blockSize != 0 {
		return "", errors.New("encrypted wrong")
	}
	// 创建数组
	decrypted := make([]byte, len(encryptedByte))
	// 解密
	blockMode := cipher.NewCBCDecrypter(a.block, a.key[:a.blockSize])
	blockMode.CryptBlocks(decrypted, encryptedByte)
	// 去补全码
	decrypted, err = pkcs7UnPadding(decrypted)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func (a *AesCipher) Decrypt(cipherText []byte) ([]byte, error) {
	decrypted, err := a.DecryptToString(string(cipherText))
	if err != nil {
		return nil, err
	}
	return []byte(decrypted), nil
}

// 补码
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// 去码
func pkcs7UnPadding(data []byte) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("UnPadding failed")
		}
	}()
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid UnPadding string")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}
