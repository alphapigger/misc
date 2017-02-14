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
	port = flag.Int("port", 8080, "listen port")
)

func main() {
	flag.Parse()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)

	// NOTE: if socket doesn't bind to addr, kernel will use random port(ip: 0.0.0.0) to fill it.
	if err = syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: *port,
	}); err != nil {
		panic(err)
	}

	if err = syscall.Listen(fd, 1); err != nil {
		panic(err)
	}

	for {
		fmt.Println("prepare to accpet new connection")
		nfd, sa, err := syscall.Accept(fd)
		if err != nil {
			fmt.Println("accpet error:", err)
			panic(err)
		}
		sa4 := sa.(*syscall.SockaddrInet4)
		rip := net.IPv4(sa4.Addr[0], sa4.Addr[1], sa4.Addr[2], sa4.Addr[3])
		raddr := fmt.Sprintf("%s:%d", rip.String(), sa4.Port)
		log.Printf("new connection, fd: %d, remote addr: %s\n", nfd, raddr)

		go func() {
			buf := make([]byte, 128)
			n, err := syscall.Read(nfd, buf)
			if err != nil {
				panic(err)
			}
			fmt.Println("read", string(buf[:n]))

			time.Sleep(time.Second * 15)
			log.Printf("prepare to send message to %s\n", raddr)
			_, err = syscall.Write(nfd, []byte("hello, i'm server"))
			if err != nil {
				log.Println("send message failed", err)
			} else {
				log.Println("send message success")
			}
			syscall.Close(nfd)
		}()
	}
}
