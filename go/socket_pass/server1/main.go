package main

import (
	"log"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func createDomainSocket(name string) (fd int, err error) {
	fd, err = unix.Socket(unix.AF_UNIX, unix.SOCK_DGRAM, 0)
	if err != nil {
		return
	}

	sa := &unix.SockaddrUnix{Name: name}
	if err = unix.Bind(fd, sa); err != nil {
		unix.Close(fd)
	}
	return
}

func listen(port int) (fd int, err error) {
	fd, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.O_NONBLOCK, 0)
	if err != nil {
		return
	}

	sa := &unix.SockaddrInet4{
		Port: 8888,
		Addr: [4]byte{0, 0, 0, 0},
	}
	if err = unix.Bind(fd, sa); err != nil {
		unix.Close(fd)
		return
	}
	if err = unix.Listen(fd, 10); err != nil {
		unix.Close(fd)
		return
	}
	return
}

func main() {
	sock, err := createDomainSocket("@parent")
	if err != nil {
		panic(err)
	}
	defer unix.Close(sock)

	fd, err := listen(8888)
	if err != nil {
		panic(err)
	}
	defer unix.Close(fd)

	go func() {
		time.Sleep(time.Second * 10) // wait children started

		log.Println("start pass file descriptor to another process")
		rights := unix.UnixRights(fd)
		if err := unix.Sendmsg(sock, nil, rights, &unix.SockaddrUnix{Name: "@child"}, 0); err != nil {
			log.Fatalf("send fd %d to another process failed: %v\n", fd, err)
		}
		log.Printf("send fd %d to another process success\n", fd)
	}()

	for {
		connFd, _, err := unix.Accept(fd)
		if err != nil {
			errno, ok := err.(syscall.Errno)
			if ok && errno == syscall.EAGAIN {
				continue
			}
			log.Fatalf("conn: accept failed, %s\n", err)
		}

		log.Printf("conn: accept connection, fd is %d\n", connFd)
		go func(fd int) {
			defer unix.Close(fd)

			var buf [10]byte
			for {
				n, err := unix.Read(fd, buf[:])
				if err != nil {
					log.Printf("conn: read from fd %d failed, %v\n", fd, err)
					return
				}
				if n > 0 {
					log.Printf("conn: read from fd %d, data is: %s\n", fd, buf[:n])
				}
			}
		}(connFd)
	}
}
