package main

import (
	"log"
	"syscall"

	"golang.org/x/sys/unix"
)

func createAndBindDomainSocket(name string) (fd int, err error) {
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

func main() {
	sock, err := createAndBindDomainSocket("@child")
	if err != nil {
		panic(err)
	}
	defer unix.Close(sock)

	b := make([]byte, unix.CmsgSpace(4))
	_, _, _, _, err = unix.Recvmsg(sock, nil, b, 0)
	if err != nil {
		panic(err)
	}

	// parse socket control message
	cmsgs, err := unix.ParseSocketControlMessage(b)
	if err != nil {
		panic(err)
	}
	fds, err := unix.ParseUnixRights(&cmsgs[0])
	if err != nil {
		panic(err)
	}
	fd := fds[0]
	log.Printf("got socket fd %d\n", fd)

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
