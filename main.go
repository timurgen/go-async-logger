package main

import (
	"logmonkey"
	"strconv"
)

func main() {
	var log = logmonkey.GetLogger("main")
	defer logmonkey.FlushAllLoggers()

	for i := 0; i < 1000; i++ {
		log.Info("logger record number :" + strconv.Itoa(i))
	}
}
