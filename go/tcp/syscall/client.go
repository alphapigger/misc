package main

import (
	"flag"
	"log"
	"net"
	"syscall"
)

var (
	addr = flag.String("addr", "127.0.0.1:8080", "connect addr")
)

func main() {
	flag.Parse()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)

	taddr, err := net.ResolveTCPAddr("tcp4", *addr)
	if err != nil {
		panic(err)
	}
	var ip [4]byte
	copy(ip[:4], taddr.IP.To4())
	sa := &syscall.SockaddrInet4{
		Port: taddr.Port,
		Addr: ip,
	}
	if err := syscall.Connect(fd, sa); err != nil {
		panic(err)
	}

	if _, err := syscall.Write(fd, []byte("hello, server")); err != nil {
		panic(err)
	}
	buf := make([]byte, 128)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		panic(err)
	}
	log.Println("receive from server:", string(buf[:n]))
}
