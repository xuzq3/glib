package util

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(ctx context.Context, src, dst string) (int64, error) {
	r, err := os.Open(src)
	if err != nil {
		return -1, err
	}
	defer r.Close()

	dstDir := filepath.Dir(dst)
	err = os.MkdirAll(dstDir, 0755)
	if err != nil {
		return -1, err
	}

	w, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return -1, err
	}
	defer w.Close()

	n, err := io.Copy(w, r)
	if err != nil {
		return -1, nil
	}

	return n, nil
}
