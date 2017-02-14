package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	// strings.Reader
	n, err := strings.NewReader("hello, world").Read(make([]byte, 20))
	// 12, nil
	fmt.Println(n, err)

	r := bufio.NewReader(strings.NewReader("hello, world"))
	buf := make([]byte, 20)
	n, err = r.Read(buf)
	// 12, nil
	fmt.Println(n, err)
	n, err = r.Read(buf)
	// 0, EOF
	fmt.Println(n, err)

	r = bufio.NewReader(strings.NewReader("hello, world"))
	p, err := r.Peek(20)
	// 12, EOF
	fmt.Println(len(p), err)

	r = bufio.NewReaderSize(strings.NewReader("hello, world!!! how are you"), 16)
	p, err = r.Peek(20)
	// 16(min read buffer size), ErrBufferFull
	fmt.Println(len(p), err)

	// tcp.Conn
}
