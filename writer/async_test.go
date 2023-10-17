package klog

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestAsyncWriter(t *testing.T) {
	aw := NewAsyncWriter(os.Stdout, 16)
	for i := 0; i < 10; i++ {
		_, err := aw.Write([]byte(fmt.Sprintf("test %d\n", i)))
		if err != nil {
			t.Fatal(err)
		}
	}
	aw.Close()
	time.Sleep(time.Second)
}
