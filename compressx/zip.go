package compressx

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func CompressZip(srcPath string, destPath string) error {
	srcPath = filepath.Clean(srcPath)

	dw, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dw.Close()

	zw := zip.NewWriter(dw)
	defer zw.Close()
	// defer func() {
	// 	errf := zw.Close()
	// 	if err == nil && errf != nil {
	// 		err = errf
	// 	}
	// }()

	fi, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	_, srcRelative := filepath.Split(srcPath)
	if fi.IsDir() {
		err = compressZipDir(srcPath, srcRelative, zw)
	} else {
		err = compressZipFile(srcPath, srcRelative, zw, fi)
	}
	if err != nil {
		return err
	}

	return nil
}

func compressZipDir(srcPath string, srcRelative string, zw *zip.Writer) error {
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
			err = compressZipDir(curPath, curRelative, zw)
			if err != nil {
				return err
			}
		} else {
			err = compressZipFile(curPath, curRelative, zw, fi)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func compressZipFile(srcPath string, srcRelative string, zw *zip.Writer, fi os.FileInfo) error {
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

	fr, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer fr.Close()

	_, err = io.Copy(w, fr)
	if err != nil {
		return err
	}

	// oneSize := 1024 * 1024
	// buf := make([]byte, oneSize)
	// for {
	// 	n, err := fr.Read(buf)
	// 	if err != nil {
	// 		if io.EOF == err {
	// 			break
	// 		} else {
	// 			return err
	// 		}
	// 	}
	// 	_, err = w.Write(buf[0:n])
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func DecompressZip(srcpath string, destpath string) error {
	reader, err := zip.OpenReader(srcpath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destpath, file.Name)
		if file.FileInfo().IsDir() {
			err = decompressZipDir(path, file)
		} else {
			err = decompressZipFile(path, file)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func decompressZipDir(path string, file *zip.File) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func decompressZipFile(path string, file *zip.File) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	r, err := file.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(path)
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
