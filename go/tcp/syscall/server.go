package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

var (
	bind = flag.String("bind", ":8080", "bind addr")
)

func main() {
	flag.Parse()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)

	addr, err := net.ResolveTCPAddr("tcp4", *bind)
	if err != nil {
		panic(err)
	}
	var ip [4]byte
	copy(ip[:4], addr.IP.To4())
	sa := syscall.SockaddrInet4{Port: addr.Port, Addr: ip}
	if err := syscall.Bind(fd, &sa); err != nil {
		panic(err)
	}

	if err = syscall.Listen(fd, 1); err != nil {
		panic(err)
	}

	for {
		nfd, nsa, err := syscall.Accept(fd)
		if err != nil {
			// TODO: check err
			panic(err)
		}

		sa4 := nsa.(*syscall.SockaddrInet4)
		ip := net.IPv4(sa4.Addr[0], sa4.Addr[1], sa4.Addr[2], sa4.Addr[3])
		naddr := fmt.Sprintf("%s:%d", ip, sa4.Port)
		log.Println("accpet new connection", naddr)

		go func() {
			defer syscall.Close(nfd)

			buf := make([]byte, 128)
			n, err := syscall.Read(nfd, buf)
			if err != nil {
				log.Printf("read from conn %s error: %v\n", naddr, err)
				return
			}
			log.Printf("read from conn %s: %s\n", naddr, buf[:n])

			// send message to client
			time.Sleep(15 * time.Second)
			if _, err := syscall.Write(nfd, []byte("hello, client")); err != nil {
				log.Printf("send message to conn %s, error: %v", naddr, err)
			}
		}()
	}
}
