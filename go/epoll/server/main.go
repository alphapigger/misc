package main

import (
	"log"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {
	fd, err := unix.Socket(unix.AF_INET, unix.O_NONBLOCK|unix.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalf("create socket fd failed: %v", err)
	}
	defer unix.Close(fd)

	unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	addr := &unix.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{0, 0, 0, 0},
	}

	if err := unix.Bind(fd, addr); err != nil {
		log.Fatalf("bind fd %d to addr %s:%d failed: %v", fd, addr.Addr, addr.Port, err)
	}
	if err := unix.Listen(fd, 4); err != nil {
		log.Fatalf("listen failed: %v", err)
	}
	log.Println("listening on", ":8000")

	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		log.Fatalf("epoll: create failed: %v", err)
	}
	defer unix.Close(epfd)

	ev := &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(fd),
	}
	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, fd, ev); err != nil {
		panic(err)
	}

	events := make([]unix.EpollEvent, 1024)
	for {
		n, err := unix.EpollWait(epfd, events, -1)
		if err != nil {
			errno, ok := err.(syscall.Errno)
			if ok && errno.Temporary() {
				continue
			}
			log.Fatalf("epoll: wait error: %v", err)
		}

		for i := 0; i < n; i++ {
			ev := events[i]
			if ev.Events&unix.EPOLLERR != 0 || ev.Events&unix.EPOLLRDHUP != 0 || ev.Events&unix.EPOLLHUP != 0 {
				log.Printf("epoll: fd %d error\n", ev.Fd)
				unix.EpollCtl(epfd, unix.EPOLL_CTL_DEL, int(ev.Fd), nil)
				unix.Close(int(ev.Fd))
				continue
			}

			if int(ev.Fd) == fd {
				connFd, _, err := unix.Accept(fd)
				if err != nil {
					log.Fatalf("conn: could not accept: %s", err)
				}

				log.Println("conn: accept new connection, fd is:", connFd)
				unix.SetNonblock(connFd, true)
				ev := &unix.EpollEvent{
					Events: unix.EPOLLIN | unix.EPOLLET | unix.EPOLLHUP | unix.EPOLLRDHUP,
					Fd:     int32(connFd),
				}
				if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, connFd, ev); err != nil {
					log.Printf("epoll: add fd %d with events %v failed: %v", connFd, ev, err)
					unix.Close(connFd)
				}
			} else {
				go func(fd int) {
					var buf [10]byte
					for {
						n, e := unix.Read(fd, buf[:])
						if e != nil {
							errno, ok := e.(syscall.Errno)
							if ok && errno == syscall.EAGAIN {
								return
							}
							log.Printf("conn: read error: %v\n", e)
							unix.EpollCtl(epfd, unix.EPOLL_CTL_DEL, fd, nil)
							unix.Close(fd)
							return
						}
						if n > 0 {
							log.Printf("conn: read from fd %d, data is: %s\n", fd, buf[:n])
						}
					}
				}(int(ev.Fd))
			}
		}
	}
}
