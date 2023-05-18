package md5pkg

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// Md5 获取MD5
func Md5(bodyBytes []byte) (string, error) {
	h := md5.New()
	_, err := h.Write(bodyBytes)
	if err != nil {
		return "", err
	}
	s := hex.EncodeToString(h.Sum(nil))
	return s, nil
}

// FileMd5 获取文件的MD5
func FileMd5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	s := hex.EncodeToString(h.Sum(nil))
	return s, nil
}
