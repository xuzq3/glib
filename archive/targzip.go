package archive

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
)

type TarGzipWriter struct {
	fw   *os.File
	gw   *gzip.Writer
	tw   *tar.Writer
	path string
}

func NewTarGzipWriter(path string) (*TarGzipWriter, error) {
	fw, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	gw := gzip.NewWriter(fw)
	tw := tar.NewWriter(gw)

	return &TarGzipWriter{
		fw:   fw,
		gw:   gw,
		tw:   tw,
		path: path,
	}, nil
}

func (p *TarGzipWriter) Close() {
	_ = p.tw.Close()
	_ = p.gw.Close()
	_ = p.fw.Close()
}

func (p *TarGzipWriter) add(ctx context.Context, hdr *tar.Header, reader io.Reader) error {
	if err := checkDone(ctx); err != nil {
		return err
	}

	err := p.tw.WriteHeader(hdr)
	if err != nil {
		return err
	}

	_, err = io.Copy(p.tw, reader)
	if err != nil {
		return err
	}
	return nil
}

func (p *TarGzipWriter) AddFileByReader(ctx context.Context, name string, fi os.FileInfo, reader io.Reader) error {
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}

	hdr.Name = filepath.ToSlash(name)
	return p.add(ctx, hdr, reader)
}

func (p *TarGzipWriter) AddFile(ctx context.Context, name string, file string) error {
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

func (p *TarGzipWriter) AddDir(ctx context.Context, name string, dir string) error {
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

//func (p *TarGzipWriter) AddFileOnlyName(name string, reader io.Reader) error {
//	hdr := &tar.Header{
//		Name:    filepath.ToSlash(name),
//		ModTime: time.Now(),
//	}
//	return p.add(hdr, reader)
//}

func CompressTarGzip(ctx context.Context, srcPath string, destPath string) error {
	writer, err := NewTarGzipWriter(destPath)
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
		err = writer.AddDir(ctx, srcRelative, srcPath)
	} else {
		err = writer.AddFile(ctx, srcRelative, srcPath)
	}
	if err != nil {
		return err
	}

	return nil
}

func DecompressTarGzip(ctx context.Context, srcPath string, destPath string) error {
	fr, err := os.Open(srcPath)
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

		path := filepath.Join(destPath, hdr.Name)
		if hdr.FileInfo().IsDir() {
			err = decompressTarGzipDir(ctx, path, hdr, tr)
		} else {
			err = decompressTarGzipFile(ctx, path, hdr, tr)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func decompressTarGzipDir(ctx context.Context, path string, hdr *tar.Header, tr *tar.Reader) error {
	err := checkDone(ctx)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, os.FileMode(hdr.Mode))
	if err != nil {
		return err
	}
	return nil
}

func decompressTarGzipFile(ctx context.Context, path string, hdr *tar.Header, tr *tar.Reader) error {
	err := checkDone(ctx)
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

	_, err = io.Copy(w, tr)
	if err != nil {
		return err
	}

	return nil
}
