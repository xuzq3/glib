package iox

import (
	"context"
	"io"
	"time"
)

const defaultBlockSize = 4 * 1024

type Copier struct {
	ctx       context.Context
	blockSize int64
	sleepTime time.Duration
}

func NewCopier(ctx context.Context) *Copier {
	return &Copier{
		ctx:       ctx,
		blockSize: defaultBlockSize,
		sleepTime: 0,
	}
}

func (c *Copier) Copy(dst io.Writer, src io.Reader) (int64, error) {
	var written int64 = 0
	var err error
	if c.blockSize <= 0 {
		written, err = io.Copy(dst, src)
		if err != nil {
			return written, err
		}
	} else {
		var n int
		buf := make([]byte, c.blockSize)
		for {
			n, err = src.Read(buf)
			if err != nil {
				if io.EOF == err {
					break
				} else {
					return written, err
				}
			}

			n, err = dst.Write(buf[0:n])
			if err != nil {
				return written, err
			}

			written += int64(n)

			if c.sleepTime > 0 {
			sleep:
				select {
				case <-c.ctx.Done():
					return written, c.ctx.Err()
				case <-time.After(c.sleepTime):
					break sleep
				}
			}
		}
	}
	return written, nil
}

func (c *Copier) SetBlockSize(size int64) *Copier {
	c.blockSize = size
	return c
}

func (c *Copier) SetSleepTime(d time.Duration) *Copier {
	c.sleepTime = d
	return c
}
