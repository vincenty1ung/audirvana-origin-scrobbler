package main

import (
	"fmt"
	"os"
	"time"
)

func clearFile(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	_, err = f.Stat()
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	now := time.Now()
	/*buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("重置时间:%s", now.String()))
	iow := bufio.NewWriter(buf)*/

	n, err := f.WriteString(fmt.Sprintf("重置时间:%s", now.String()))
	if err != nil {
		return err
	}
	fmt.Println(n)
	return nil
}
