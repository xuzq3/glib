package compressx

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func DoTarGzCompressPath(srcPath string, destPath string) error {
	srcPath = filepath.Clean(srcPath)

	fw, err := os.Create(destPath)
	if nil != err {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer func() {
		errf := tw.Close()
		if nil == err && nil != errf {
			err = errf
		}
	}()

	fi, err := os.Stat(srcPath)
	if nil != err {
		return err
	}

	_, srcRelative := filepath.Split(srcPath)
	if fi.IsDir() {
		err = doTarDir(srcPath, srcRelative, tw)
	} else {
		err = doTarFile(srcPath, srcRelative, tw, fi)
	}
	if nil != err {
		return err
	}

	return nil
}

func doTarDir(srcPath string, srcRelative string, tw *tar.Writer) error {
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
			err = doTarDir(curPath, curRelative, tw)
			if err != nil {
				return err
			}
		} else {
			err = doTarFile(curPath, curRelative, tw, fi)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func doTarFile(srcPath string, srcRelative string, tw *tar.Writer, fi os.FileInfo) error {
	hdr, err := tar.FileInfoHeader(fi, "")
	if nil != err {
		return err
	}

	hdr.Name = srcRelative
	err = tw.WriteHeader(hdr)
	if nil != err {
		return err
	}

	fr, err := os.Open(srcPath)
	if nil != err {
		return err
	}
	defer fr.Close()

	oneSize := 1024 * 1024
	buf := make([]byte, oneSize)
	for {
		n, err := fr.Read(buf)
		if nil != err {
			if io.EOF == err {
				break
			} else {
				return err
			}
		}
		_, err = tw.Write(buf[0:n])
		if nil != err {
			return err
		}
	}

	return nil
}

func DoZipCompressPath(srcPath string, destPath string) error {
	srcPath = filepath.Clean(srcPath)

	dw, err := os.Create(destPath)
	if nil != err {
		return err
	}
	defer dw.Close()

	zw := zip.NewWriter(dw)
	defer func() {
		errf := zw.Close()
		if nil == err && nil != errf {
			err = errf
		}
	}()

	fi, err := os.Stat(srcPath)
	if nil != err {
		return err
	}

	_, srcRelative := filepath.Split(srcPath)
	if fi.IsDir() {
		err = doZipCompressDir(srcPath, srcRelative, zw)
	} else {
		err = doZipCompressFile(srcPath, srcRelative, zw, fi)
	}
	if nil != err {
		return err
	}

	return nil
}

func doZipCompressDir(srcPath string, srcRelative string, zw *zip.Writer) error {
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
			err = doZipCompressDir(curPath, curRelative, zw)
			if err != nil {
				return err
			}
		} else {
			err = doZipCompressFile(curPath, curRelative, zw, fi)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func doZipCompressFile(srcPath string, srcRelative string, zw *zip.Writer, fi os.FileInfo) error {
	hdr, err := zip.FileInfoHeader(fi)
	if nil != err {
		return err
	}

	hdr.Name = srcRelative
	hdr.Method = zip.Deflate
	w, err := zw.CreateHeader(hdr)
	if nil != err {
		return err
	}

	fr, err := os.Open(srcPath)
	if nil != err {
		return err
	}
	defer fr.Close()

	oneSize := 1024 * 1024
	buf := make([]byte, oneSize)
	for {
		n, err := fr.Read(buf)
		if nil != err {
			if io.EOF == err {
				break
			} else {
				return err
			}
		}
		_, err = w.Write(buf[0:n])
		if nil != err {
			return err
		}
	}

	return nil
}
