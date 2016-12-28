package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

func main() {
	timer := time.NewTimer(time.Second * 5)
	timer.Stop()
	debug.SetTraceback("system")
	if len(os.Args) == 1 {
		panic("before timers")
	}
	for i := 0; i < 10000; i++ {
		time.AfterFunc(time.Duration(5*time.Second), func() {
			fmt.Println("hello")
		})
	}
	panic("after timers")
}
