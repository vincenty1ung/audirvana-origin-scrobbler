package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("执行清理策略")
		err := clearFile("/Users/vincent/Developer/code/other/audirvana-origin-scrobbler/logs/audirvana-origin-scrobbler.log")
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 3600)
	}
}
