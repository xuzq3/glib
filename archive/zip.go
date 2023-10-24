package archive

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
)

const (
	ZipStore   = zip.Store
	ZipDeflate = zip.Deflate
)

type ZipOption struct {
	Method uint16
}

var DefaultZipOption = ZipOption{
	Method: ZipStore,
}

type ZipWriter struct {
	fw   *os.File
	zw   *zip.Writer
	path string
	opt  ZipOption
}

func NewZipWriter(path string, opt ZipOption) (*ZipWriter, error) {
	fw, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	zw := zip.NewWriter(fw)

	return &ZipWriter{
		fw:   fw,
		zw:   zw,
		path: path,
		opt:  opt,
	}, nil
}

func (p *ZipWriter) Close() {
	_ = p.zw.Close()
	_ = p.fw.Close()
}

func (p *ZipWriter) add(ctx context.Context, hdr *zip.FileHeader, reader io.Reader) error {
	if err := checkDone(ctx); err != nil {
		return err
	}

	w, err := p.zw.CreateHeader(hdr)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, reader)
	if err != nil {
		return err
	}
	return nil
}

func (p *ZipWriter) AddFileByReader(ctx context.Context, name string, fi os.FileInfo, reader io.Reader) error {
	hdr, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}
	hdr.Name = filepath.ToSlash(name)
	hdr.Method = p.opt.Method

	return p.add(ctx, hdr, reader)
}

func (p *ZipWriter) AddFile(ctx context.Context, name string, file string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}

	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()

	err = p.AddFileByReader(ctx, name, fi, r)
	if err != nil {
		return err
	}
	return nil
}

func (p *ZipWriter) AddDir(ctx context.Context, name string, dir string) error {
	if err := checkDone(ctx); err != nil {
		return err
	}

	df, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer func() { _ = df.Close() }()

	fis, err := df.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		childPath := filepath.Join(dir, fi.Name())
		childName := filepath.Join(name, fi.Name())
		if fi.IsDir() {
			err = p.AddDir(ctx, childName, childPath)
			if err != nil {
				return err
			}
		} else {
			err = p.AddFile(ctx, childName, childPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//func (p *ZipWriter) AddFileOnlyName(name string, reader io.Reader) error {
//	hdr := &zip.FileHeader{
//		Name:     filepath.ToSlash(name),
//		Method:   p.opt.Method,
//		Modified: time.Now(),
//	}
//	return p.add(hdr, reader)
//}

func CompressZip(ctx context.Context, srcPath string, destPath string) error {
	writer, err := NewZipWriter(destPath, DefaultZipOption)
	if err != nil {
		return err
	}
	defer writer.Close()

	srcPath = filepath.Clean(srcPath)
	fi, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	_, name := filepath.Split(srcPath)
	if fi.IsDir() {
		err = writer.AddDir(ctx, name, srcPath)
	} else {
		err = writer.AddFile(ctx, name, srcPath)
	}
	if err != nil {
		return err
	}
	return nil
}

//
//func compressZipDir(ctx context.Context, srcPath string, srcRelative string, zw *ZipWriter) error {
//	if err := checkDone(ctx); err != nil {
//		return err
//	}
//
//	dir, err := os.Open(srcPath)
//	if err != nil {
//		return err
//	}
//	defer dir.Close()
//
//	fis, err := dir.Readdir(0)
//	if err != nil {
//		return err
//	}
//
//	for _, fi := range fis {
//		curPath := filepath.Join(srcPath, fi.Name())
//		curRelative := filepath.Join(srcRelative, fi.Name())
//		if fi.IsDir() {
//			err = compressZipDir(ctx, curPath, curRelative, zw)
//			if err != nil {
//				return err
//			}
//		} else {
//			err = compressZipFile(ctx, curPath, curRelative, zw, fi)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//func compressZipFile(ctx context.Context, srcPath string, srcRelative string, zw *ZipWriter, fi os.FileInfo) error {
//	if err := checkDone(ctx); err != nil {
//		return err
//	}
//
//	r, err := os.Open(srcPath)
//	if err != nil {
//		return err
//	}
//	defer r.Close()
//
//	err = zw.AddFileByReader(fi, srcRelative, r)
//	if err != nil {
//		return err
//	}
//	return nil
//}

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
