package main

import (
	"logmonkey"
	"strconv"
	"time"
)

func main() {
	var log = logmonkey.GetLogger("main")

	for i := 0; i < 100; i++ {
		log.Info("logger record number :" + strconv.Itoa(i))
	}
	time.Sleep(1000*time.Millisecond)
}
