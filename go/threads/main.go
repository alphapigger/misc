package main

import (
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
)

func main() {
	runtime.GOMAXPROCS(2)
	data := make([]byte, 128*1024*1024)
	for i := 0; i < 200; i++ {
		go func(n int) {
			for {
				// too long syscall will cause allocating excussive threads
				err := ioutil.WriteFile("testxxx"+strconv.Itoa(n), []byte(data), os.ModePerm)
				if err != nil {
					println(err)
					break
				}
			}
		}(i)
	}
	select {}
}
