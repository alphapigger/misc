package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	addr = flag.String("addr", ":8080", "listen address")
)

func main() {
	var stop bool
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		stop = true
		fmt.Println("caught signal:", sig)
	}()

	l, err := net.Listen("tcp4", *addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("listening on:", l.Addr())
	for {
		if stop {
			return
		}
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("accept new conn", conn.LocalAddr(), conn.RemoteAddr())
		}
	}
}
