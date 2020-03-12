package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//偷懒
func RetryForever(f func() error, interval time.Duration) error {
	return Retry(f, 99999999999, interval)
}

func Retry(f func() error, attempts int, interval time.Duration) error {
	err := f()
	if err == nil {
		return nil
	}
	var num = 0
	for {
		num++
		if attempts <= num {
			break
		}
		time.Sleep(interval)
		if err = f(); err == nil {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("retry fail after %d try", attempts))
}

func CombineString(strs ...string) string {
	var buffer bytes.Buffer
	for _, v := range strs {
		buffer.WriteString(v)
	}
	return buffer.String()
}

func Abs(f string) (p string) {
	p, _ = filepath.Abs(f)
	return
}

//unexcepted path when go test or working with ide
func ExecPath() string {
	file, _ := exec.LookPath(os.Args[0])
	//得到全路径，比如在windows下E:\\golang\\test\\a.exe
	return filepath.Dir(Abs(file))
}
