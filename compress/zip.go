package compress

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"

	"github.com/xuzq3/glib/util"
)

type ZipCompresser struct {
	ctx      context.Context
	srcPath  string
	destPath string
	opts     *Options
}

func NewZipCompresser(ctx context.Context, srcPath string, destPath string, opts *Options) *ZipCompresser {
	if opts == nil {
		opts = NewOptions()
	}
	return &ZipCompresser{
		ctx:      ctx,
		srcPath:  srcPath,
		destPath: destPath,
		opts:     opts,
	}
}

func (c *ZipCompresser) Compress() error {
	srcPath := filepath.Clean(c.srcPath)
	dw, err := os.Create(c.destPath)
	if err != nil {
		return err
	}
	defer dw.Close()

	zw := zip.NewWriter(dw)
	defer zw.Close()

	fi, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	_, srcRelative := filepath.Split(srcPath)
	if fi.IsDir() {
		err = c.compressDir(srcPath, srcRelative, zw)
	} else {
		err = c.compressFile(srcPath, srcRelative, zw, fi)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *ZipCompresser) compressDir(srcPath string, srcRelative string, zw *zip.Writer) error {
	err := checkDone(c.ctx)
	if err != nil {
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
			err = c.compressDir(curPath, curRelative, zw)
			if err != nil {
				return err
			}
		} else {
			err = c.compressFile(curPath, curRelative, zw, fi)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ZipCompresser) compressFile(srcPath string, srcRelative string, zw *zip.Writer, fi os.FileInfo) error {
	err := checkDone(c.ctx)
	if err != nil {
		return err
	}

	hdr, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}
	hdr.Name = filepath.ToSlash(srcRelative)
	hdr.Method = zip.Deflate
	w, err := zw.CreateHeader(hdr)
	if err != nil {
		return err
	}
	r, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer r.Close()

	copier := util.NewCopier(c.ctx)
	copier.SetBlockSize(c.opts.BlockSize)
	copier.SetDelayTime(c.opts.DelayTime)
	_, err = copier.Copy(w, r)
	if err != nil {
		return err
	}

	return nil
}

type ZipDecompresser struct {
	ctx      context.Context
	srcPath  string
	destPath string
	opts     *Options
}

func NewZipDecompresser(ctx context.Context, srcPath string, destPath string, opts *Options) *ZipDecompresser {
	if opts == nil {
		opts = NewOptions()
	}
	return &ZipDecompresser{
		ctx:      ctx,
		srcPath:  srcPath,
		destPath: destPath,
		opts:     opts,
	}
}

func (d *ZipDecompresser) Decompress() error {
	reader, err := zip.OpenReader(d.srcPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(d.destPath, file.Name)
		if file.FileInfo().IsDir() {
			err = d.decompressDir(path, file)
		} else {
			err = d.decompressFile(path, file)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (d *ZipDecompresser) decompressDir(path string, file *zip.File) error {
	err := checkDone(d.ctx)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, file.Mode())
	if err != nil {
		return err
	}
	return nil
}

func (d *ZipDecompresser) decompressFile(path string, file *zip.File) error {
	err := checkDone(d.ctx)
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

	copier := util.NewCopier(d.ctx)
	copier.SetBlockSize(d.opts.BlockSize)
	copier.SetDelayTime(d.opts.DelayTime)
	_, err = copier.Copy(w, r)
	if err != nil {
		return err
	}

	return nil
}
