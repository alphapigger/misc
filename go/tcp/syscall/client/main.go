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
	addr = flag.String("addr", "127.0.0.1:8080", "server addr")
)

func main() {
	flag.Parse()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}

	taddr, err := net.ResolveTCPAddr("tcp4", *addr)
	if err != nil {
		panic(err)
	}
	var ip [4]byte
	copy(ip[:4], taddr.IP.To4())
	if err := syscall.Connect(fd, &syscall.SockaddrInet4{
		Port: taddr.Port,
		Addr: ip,
	}); err != nil {
		panic(err)
	}

	go func() {
		_, err := syscall.Write(fd, []byte("hello, i'm client"))
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 128)
		n, err := syscall.Read(fd, buf)
		if err != nil {
			panic(err)
		}
		fmt.Println("read", string(buf[:n]))
		syscall.Close(fd)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	fmt.Println("caught signal", sig)
}
