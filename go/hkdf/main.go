package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

func main() {
	hash := sha256.New

	master := []byte("hello")

	hf1 := hkdf.New(hash, master, nil, nil)

	hf2 := hkdf.New(hash, master, nil, nil)

	// Generate the required keys
	keys1 := make([][]byte, 3)
	keys2 := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		keys1[i] = make([]byte, 100)
		n, err := io.ReadFull(hf1, keys1[i])
		if n != len(keys1[i]) || err != nil {
			fmt.Println("error:", err)
			return
		}

		keys2[i] = make([]byte, 100)
		n, err = io.ReadFull(hf2, keys2[i])
		if n != len(keys2[i]) || err != nil {
			fmt.Println("error:", err)
			return
		}
	}
	// Keys should contain 192 bit random keys
	// for i := 1; i <= len(keys); i++ {
	// fmt.Printf("Key #%d: %v\n", i, !bytes.Equal(keys[i-1], make([]byte, 24)))
	// }
	for i := 1; i <= 3; i++ {
		fmt.Printf("Key #%d: %v\n", i, bytes.Equal(keys1[i-1], keys2[i-1]))
	}
}
