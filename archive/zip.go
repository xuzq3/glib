package archive

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	ZipStore   = zip.Store
	ZipDeflate = zip.Deflate
)

type ZipOption struct {
	Method uint16
}

type OptionFunc func(option *ZipOption) *ZipOption

func WithZipMethod(method uint16) OptionFunc {
	return func(o *ZipOption) *ZipOption {
		o.Method = method
		return o
	}
}

type ZipWriter struct {
	dw   *os.File
	zw   *zip.Writer
	path string
	opt  *ZipOption
}

func NewZipWriter(path string, opts ...OptionFunc) (*ZipWriter, error) {
	dw, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	zw := zip.NewWriter(dw)

	opt := &ZipOption{}
	for _, o := range opts {
		o(opt)
	}

	return &ZipWriter{
		dw:   dw,
		zw:   zw,
		path: path,
		opt:  opt,
	}, nil
}

func (c *ZipWriter) Close() {
	_ = c.zw.Close()
	_ = c.dw.Close()
}

func (c *ZipWriter) add(hdr *zip.FileHeader, reader io.Reader) error {
	w, err := c.zw.CreateHeader(hdr)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, reader)
	if err != nil {
		return err
	}
	return nil
}

func (c *ZipWriter) AddFile(fi os.FileInfo, name string, reader io.Reader) error {
	hdr, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}
	hdr.Name = filepath.ToSlash(name)
	hdr.Method = c.opt.Method

	return c.add(hdr, reader)
}

func (c *ZipWriter) AddFileOnlyName(name string, reader io.Reader) error {
	hdr := &zip.FileHeader{
		Name:     name,
		Method:   c.opt.Method,
		Modified: time.Now(),
	}
	return c.add(hdr, reader)
}

func CompressZip(ctx context.Context, srcPath string, destPath string) error {
	writer, err := NewZipWriter(destPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	srcPath = filepath.Clean(srcPath)
	fi, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	_, srcRelative := filepath.Split(srcPath)
	if fi.IsDir() {
		err = compressZipDir(ctx, srcPath, srcRelative, writer)
	} else {
		err = compressZipFile(ctx, srcPath, srcRelative, writer, fi)
	}
	if err != nil {
		return err
	}
	return nil
}

func compressZipDir(ctx context.Context, srcPath string, srcRelative string, zw *ZipWriter) error {
	if err := checkDone(ctx); err != nil {
		return err
	}

	dir, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		curPath := filepath.Join(srcPath, fi.Name())
		curRelative := filepath.Join(srcRelative, fi.Name())
		if fi.IsDir() {
			err = compressZipDir(ctx, curPath, curRelative, zw)
			if err != nil {
				return err
			}
		} else {
			err = compressZipFile(ctx, curPath, curRelative, zw, fi)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func compressZipFile(ctx context.Context, srcPath string, srcRelative string, zw *ZipWriter, fi os.FileInfo) error {
	if err := checkDone(ctx); err != nil {
		return err
	}

	r, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer r.Close()

	err = zw.AddFile(fi, srcRelative, r)
	if err != nil {
		return err
	}
	return nil
}

func DecompressZip(ctx context.Context, srcPath string, destPath string) error {
	reader, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destPath, file.Name)
		if file.FileInfo().IsDir() {
			err = decompressZipDir(ctx, path, file)
		} else {
			err = decompressZipFile(ctx, path, file)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func decompressZipDir(ctx context.Context, path string, file *zip.File) error {
	err := checkDone(ctx)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, file.Mode())
	if err != nil {
		return err
	}
	return nil
}

func decompressZipFile(ctx context.Context, path string, file *zip.File) error {
	err := checkDone(ctx)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	r, err := file.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	return nil
}
