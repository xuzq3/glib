package archive

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
)

func CompressZip(ctx context.Context, srcPath string, destPath string) error {
	srcPath = filepath.Clean(srcPath)
	dw, err := os.Create(destPath)
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
		err = compressZipDir(ctx, srcPath, srcRelative, zw)
	} else {
		err = compressZipFile(ctx, srcPath, srcRelative, zw, fi)
	}
	if err != nil {
		return err
	}
	return nil
}

func compressZipDir(ctx context.Context, srcPath string, srcRelative string, zw *zip.Writer) error {
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

func compressZipFile(ctx context.Context, srcPath string, srcRelative string, zw *zip.Writer, fi os.FileInfo) error {
	err := checkDone(ctx)
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

	_, err = io.Copy(w, r)
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
