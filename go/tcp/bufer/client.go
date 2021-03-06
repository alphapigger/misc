package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.Println("begin dial...")
	raddr, err := net.ResolveTCPAddr("tcp4", ":8888")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp4", nil, raddr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	if err := conn.SetWriteBuffer(1024); err != nil {
		panic(err)
	}
	log.Println("dial ok")

	data := make([]byte, 65536)
	var total int
	for {
		n, err := conn.Write(data)
		if err != nil {
			total += n
			log.Printf("write %d bytes, error: %s\n", n, err)
			break
		}
		total += n
		log.Printf("write %d bytes this time, %d bytes in total\n", n, total)
	}
	log.Printf("write %d bytes in total\n", total)
	time.Sleep(time.Second * 10000)
}
