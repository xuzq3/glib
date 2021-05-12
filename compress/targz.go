package compress

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/xuzq3/glib/util"
)

type TargzCompresser struct {
	ctx      context.Context
	srcPath  string
	destPath string
	opts     *Options
}

func NewTargzCompresser(ctx context.Context, srcPath string, destPath string, opts *Options) *TargzCompresser {
	if opts == nil {
		opts = NewOptions()
	}
	return &TargzCompresser{
		ctx:      ctx,
		srcPath:  srcPath,
		destPath: destPath,
		opts:     opts,
	}
}

func (c *TargzCompresser) Compress() error {
	srcPath := filepath.Clean(c.srcPath)
	fw, err := os.Create(c.destPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	fi, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	_, srcRelative := filepath.Split(srcPath)
	if fi.IsDir() {
		err = c.compressDir(srcPath, srcRelative, tw)
	} else {
		err = c.compressFile(srcPath, srcRelative, tw, fi)
	}
	if err != nil {
		return err
	}

	return nil
}

func (c *TargzCompresser) compressDir(srcPath string, srcRelative string, tw *tar.Writer) error {
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
			err = c.compressDir(curPath, curRelative, tw)
			if err != nil {
				return err
			}
		} else {
			err = c.compressFile(curPath, curRelative, tw, fi)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (c *TargzCompresser) compressFile(srcPath string, srcRelative string, tw *tar.Writer, fi os.FileInfo) error {
	err := checkDone(c.ctx)
	if err != nil {
		return err
	}

	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}

	hdr.Name = filepath.ToSlash(srcRelative)
	err = tw.WriteHeader(hdr)
	if err != nil {
		return err
	}

	fr, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer fr.Close()

	copier := util.NewCopier(c.ctx)
	copier.SetBlockSize(c.opts.BlockSize)
	copier.SetDelayTime(c.opts.DelayTime)
	_, err = copier.Copy(tw, fr)
	if err != nil {
		return err
	}

	return nil
}

type TargzDecompresser struct {
	ctx      context.Context
	srcPath  string
	destPath string
	opts     *Options
}

func NewTargzDecompresser(ctx context.Context, srcPath string, destPath string, opts *Options) *TargzDecompresser {
	if opts == nil {
		opts = NewOptions()
	}
	return &TargzDecompresser{
		ctx:      ctx,
		srcPath:  srcPath,
		destPath: destPath,
		opts:     opts,
	}
}

func (d *TargzDecompresser) Decompress() error {
	fr, err := os.Open(d.srcPath)
	if err != nil {
		return err
	}
	defer fr.Close()

	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		path := filepath.Join(d.destPath, hdr.Name)
		if hdr.FileInfo().IsDir() {
			err = d.decompressDir(path, hdr, tr)
		} else {
			err = d.decompressFile(path, hdr, tr)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *TargzDecompresser) decompressDir(path string, hdr *tar.Header, tr *tar.Reader) error {
	err := checkDone(d.ctx)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, os.FileMode(hdr.Mode))
	if err != nil {
		return err
	}
	return nil
}

func (d *TargzDecompresser) decompressFile(path string, hdr *tar.Header, tr *tar.Reader) error {
	err := checkDone(d.ctx)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	w, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
	if err != nil {
		return err
	}
	defer w.Close()

	copier := util.NewCopier(d.ctx)
	copier.SetBlockSize(d.opts.BlockSize)
	copier.SetDelayTime(d.opts.DelayTime)
	_, err = copier.Copy(w, tr)
	if err != nil {
		return err
	}

	return nil
}
