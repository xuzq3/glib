package compressx

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func CompressTargz(srcPath string, destPath string) error {
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
	// defer func() {
	// 	errf := tw.Close()
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
		err = compressTargzDir(srcPath, srcRelative, tw)
	} else {
		err = compressTargzFile(srcPath, srcRelative, tw, fi)
	}
	if err != nil {
		return err
	}

	return nil
}

func compressTargzDir(srcPath string, srcRelative string, tw *tar.Writer) error {
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
			err = compressTargzDir(curPath, curRelative, tw)
			if err != nil {
				return err
			}
		} else {
			err = compressTargzFile(curPath, curRelative, tw, fi)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func compressTargzFile(srcPath string, srcRelative string, tw *tar.Writer, fi os.FileInfo) error {
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
	// 	_, err = tw.Write(buf[0:n])
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func DecompressTargz(srcpath string, destpath string) error {
	fr, err := os.Open(srcpath)
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

		path := filepath.Join(destpath, hdr.Name)
		if hdr.FileInfo().IsDir() {
			err = decompressTargzDir(path, tr)
		} else {
			err = decompressTargzFile(path, tr)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func decompressTargzDir(path string, tr *tar.Reader) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func decompressTargzFile(path string, tr *tar.Reader) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	w, err := os.Create(path)
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
