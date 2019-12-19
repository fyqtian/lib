package utils

import (
	"bytes"
	"errors"
	"fmt"
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
