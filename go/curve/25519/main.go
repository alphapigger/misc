package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"
)

func main() {
	var (
		pub1, priv1, shared1 [32]byte
		pub2, priv2, shared2 [32]byte
	)

	if _, err := io.ReadFull(rand.Reader, priv1[:]); err != nil {
		panic(err)
	}
	if _, err := io.ReadFull(rand.Reader, priv2[:]); err != nil {
		panic(err)
	}

	priv1[0] &= 248
	priv1[31] &= 127
	priv1[31] |= 64

	priv2[0] &= 248
	priv2[31] &= 127
	priv2[31] |= 64

	curve25519.ScalarBaseMult(&pub1, &priv1)
	curve25519.ScalarBaseMult(&pub2, &priv2)

	curve25519.ScalarMult(&shared1, &priv1, &pub2)
	curve25519.ScalarMult(&shared2, &priv2, &pub1)
	fmt.Println("shared1: ", hex.EncodeToString(shared1[:]))
	fmt.Println("shared2: ", hex.EncodeToString(shared2[:]))
	if shared1 != shared2 {
		fmt.Println("shared1 not equal to shared2")
	}
}
