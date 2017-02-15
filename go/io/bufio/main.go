package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func stringsReader() {
	sr := strings.NewReader("hello, world")
	// 12, nil
	n, err := sr.Read(make([]byte, 20))
	log.Printf("raw strings.Reader: read data length is: %d, err is: %v\n", n, err)
	// 0, EOF
	n, err = sr.Read(make([]byte, 20))
	log.Printf("raw strings.Reader: read data length is: %d, err is: %v\n", n, err)

	r := bufio.NewReader(strings.NewReader("hello, world"))
	// 12, nil
	n, err = r.Read(make([]byte, 20))
	log.Printf("bufio strings.Reader: read data lenght is: %d, err is: %v\n", n, err)
	// 0, EOF
	n, err = r.Read(make([]byte, 20))
	log.Printf("bufio strings.Reader: read data length is: %d, err is: %v\n", n, err)

	r = bufio.NewReader(strings.NewReader("hello, world"))
	// 12, EOF
	p, err := r.Peek(20)
	log.Printf("bufio strings.Reader: peek data length is: %d, err is: %v\n", len(p), err)

	r = bufio.NewReaderSize(strings.NewReader("hello, world!"), 17)
	// 13, nil
	p, err = r.Peek(13)
	log.Printf("bufio strings.Reader: peek data length is: %d, err is: %v\n", len(p), err)

	r = bufio.NewReaderSize(strings.NewReader("hello, world! how are you"), 16)
	// 16 (min read buffer size), ErrBufferFull
	p, err = r.Peek(20)
	log.Printf("bufio strings.Reader: peek data length is: %d, err is: %v\n", len(p), err)
}

func tcp() {
	lis, err := net.ListenTCP("tcp4", new(net.TCPAddr))
	if err != nil {
		panic(err)
	}
	var isClosed bool
	defer func() {
		isClosed = true
		lis.Close()
	}()

	go func() {
		for {
			conn, err := lis.AcceptTCP()
			if err != nil {
				if !isClosed {
					panic(err)
				} else {
					return
				}
			}
			log.Printf("tcp server: accept new conn: %s\n", conn.RemoteAddr())
			// for demo: use blocked processing
			buf := make([]byte, 13)
			// "hello, world", nil
			n, err := conn.Read(buf)
			log.Printf("tcp server: read data %s from %s, length is: %d, err is: %v\n", buf[:n], conn.RemoteAddr(), n, err)
			buf = make([]byte, 20)
			time.Sleep(time.Millisecond * 10)
			// "how are youhaha", nil
			n, err = conn.Read(buf)
			log.Printf("tcp server: read data %s from %s, length is: %d, err is: %v\n", buf[:n], conn.RemoteAddr(), n, err)
			// 0, EOF
			n, err = conn.Read(buf)
			log.Printf("tcp server: read data %s from %s, length is: %d, err is: %v\n", buf[:n], conn.RemoteAddr(), n, err)
			if err == io.EOF {
				conn.Write([]byte("hahah"))
				conn.Close()
			}
		}
	}()

	conn, err := net.Dial(lis.Addr().Network(), lis.Addr().String())
	if err != nil {
		panic(err)
	}
	n, err := conn.Write([]byte("hello, world!"))
	log.Printf("tcp client: send data('hello, world!') to %s, length is: %d, err is: %v\n", conn.RemoteAddr(), n, err)
	n, err = conn.Write([]byte("how are you"))
	log.Printf("tcp client: send data('how are you') to %s, length is: %d, err is: %v\n", conn.RemoteAddr(), n, err)
	time.Sleep(time.Millisecond * 1)
	conn.Write([]byte("haha"))
	conn.Close()
	time.Sleep(time.Second * 1)
}

func tcpBufio() {
	lis, err := net.ListenTCP("tcp4", new(net.TCPAddr))
	if err != nil {
		panic(err)
	}
	var isClosed bool
	defer func() {
		isClosed = true
		lis.Close()
	}()

	go func() {
		for {
			conn, err := lis.AcceptTCP()
			if err != nil {
				if !isClosed {
					panic(err)
				} else {
					return
				}
			}
			r := bufio.NewReaderSize(conn, 30)

			b, err := r.Peek(17)
			// hello, world!!!!!, nil
			log.Printf("bufio tcp(server): receive data %s from %s, err: %v\n", b, conn.RemoteAddr(), err)

			b, err = r.Peek(28)
			// hello, world!!!!!👌hihaha, nil
			log.Printf("bufio tcp(server): receive data %s from %s, err: %v\n", b, conn.RemoteAddr(), err)

			// reset buffer
			r.Read(make([]byte, 27))

			b, err = r.Peek(35)
			// "hehe", ErrBufferFull
			log.Printf("bufio tcp(server): receive data %s from %s, err: %v\n", b, conn.RemoteAddr(), err)

			b, err = r.Peek(10)
			// "hehe", EOF
			log.Printf("bufio tcp(server): receive data %s from %s, err: %v\n", b, conn.RemoteAddr(), err)

			conn.Close()
		}
	}()

	conn, err := net.Dial(lis.Addr().Network(), lis.Addr().String())
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriterSize(conn, 64)
	n, err := w.Write([]byte("hello, world!"))
	log.Printf("bufio tcp(client): send data('hello, world!') to %s, length is: %d, err is: %v\n", conn.RemoteAddr(), n, err)
	n, err = w.Write([]byte("!!!!👌"))
	log.Printf("bufio tcp(client): send data(' hello, elese!') to %s, length is: %d, err is: %v\n", conn.RemoteAddr(), n, err)
	err = w.Flush()
	log.Printf("bufio tcp(client): flush err: %v\n", err)
	w.Write([]byte("hi"))
	w.Flush()

	n, err = w.Write([]byte("haha"))
	log.Printf("bufio tcp(client): send data('haha') to %s, length is: %d, err is: %v\n", conn.RemoteAddr(), n, err)
	w.Flush()

	n, err = w.Write([]byte("hehe"))
	w.Flush()
	time.Sleep(time.Second * 1)
	conn.Close()
	time.Sleep(time.Second)
}

func main() {
	fmt.Println("=================== strings.Reader ==================")
	stringsReader()
	fmt.Println("=====================================================\n\n")

	fmt.Println("=================== tcp ReadWriter ==================")
	tcp()
	tcpBufio()
	fmt.Println("=====================================================\n\n")
}
