package logx

import (
	"errors"
	"fmt"
	"os"
)

const maxInt = int(^uint(0) >> 1)

func ensureDirExist(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				return err
			}
		}
	} else {
		if !fi.IsDir() {
			return errors.New(fmt.Sprintf("%s is not a directory", dir))
		}
	}
	return nil
}
