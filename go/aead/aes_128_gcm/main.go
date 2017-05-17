package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"reflect"
)

var (
	key        = []byte("1234567812345678")
	nonce      = make([]byte, 12)
	plaintext  = []byte("hello, world")
	ciphertext []byte
)

func init() {
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
}

func encrypt() {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	enc, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	ciphertext = enc.Seal(nil, nonce, plaintext, nil)
}

func decrypt() {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	dec, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	plaintext1, err := dec.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("decrypt error: ", err)
		return
	}
	if !reflect.DeepEqual(plaintext, plaintext1) {
		fmt.Println("raw plian text", plaintext, "decrypt plain text", plaintext1)
	}
	fmt.Println("decrypt success")
}

func main() {
	encrypt()
	decrypt()
}
