package archive

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
)

func CompressTarGzip(ctx context.Context, srcPath string, destPath string) error {
	srcPath = filepath.Clean(srcPath)
	fw, err := os.Create(destPath)
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
		err = compressTarGzipDir(ctx, srcPath, srcRelative, tw)
	} else {
		err = compressTarGzipFile(ctx, srcPath, srcRelative, tw, fi)
	}
	if err != nil {
		return err
	}

	return nil
}

func compressTarGzipDir(ctx context.Context, srcPath string, srcRelative string, tw *tar.Writer) error {
	err := checkDone(ctx)
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
			err = compressTarGzipDir(ctx, curPath, curRelative, tw)
			if err != nil {
				return err
			}
		} else {
			err = compressTarGzipFile(ctx, curPath, curRelative, tw, fi)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func compressTarGzipFile(ctx context.Context, srcPath string, srcRelative string, tw *tar.Writer, fi os.FileInfo) error {
	err := checkDone(ctx)
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

	_, err = io.Copy(tw, fr)
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
