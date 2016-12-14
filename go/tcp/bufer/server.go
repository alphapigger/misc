package main

import (
	"log"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp4", ":8888")
	if err != nil {
		panic(err)
	}
	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			panic(err)
		}
		if err := conn.SetReadBuffer(1024); err != nil {
			panic(err)
		}
		go func() {
			defer conn.Close()
			time.Sleep(time.Second * 10)
			for {
				time.Sleep(5 * time.Second)
				buf := make([]byte, 65536)
				log.Println("start to read from conn")
				n, err := conn.Read(buf)
				if err != nil {
					log.Printf("conn read %d bytes, error: %s", n, err)
					if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
						continue
					}
				}

				log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
			}
		}()
	}
}
