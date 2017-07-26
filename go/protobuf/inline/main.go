package main

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/kirk91/misc/go/protobuf/inline/pb"
)

func main() {
	p, _ := Marshal(&pb.User{Name: "哈哈", Age: 23})
	au := &pb.AliasUser{}
	fmt.Println(au.Unmarshal(p), au.Name)

	// cannot inline, need to find the reason
	p, _ = Marshal(&pb.User{Name: "哈哈", Age: 23})
	p = MustMarshal(&pb.User{Age: 23})
	print(p)

	// a, b := 1, 2
	// max := Max(1, 2)
	// fmt.Println(a, b, max)
	buf := make([]byte, 100)
	n, err := (&pb.User{Name: "嘿嘿", Age: 30}).MarshalTo(buf)
	fmt.Println(n, err)

	// test nil struct
	a := &pb.A{Name: "a"}
	buf = make([]byte, 100)
	a.MarshalTo(buf)
	ua := new(pb.A)
	ua.Unmarshal(buf)
	fmt.Println(ua.B)
}

func Marshal(m proto.Message) ([]byte, error) {
	return proto.Marshal(m)
}

func print(p []byte) {
	fmt.Println("marshaled data length:", len(p))
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MustMarshal(m proto.Message) []byte {
	b, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}
	return b
}
