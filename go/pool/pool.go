package pool

import (
	"fmt"
	"os"
	"sync"
)

type buf struct {
	bytes []byte
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return &buf{
			bytes: make([]byte, 0, 1024),
		}
	},
}

func WriteString(s string) {
	b := make([]byte, 0, 1024)
	b = append(b, s...)

	f, _ := os.Open("/dev/null")
	defer f.Close()

	fmt.Fprint(f, b)
}

func WriteStringWithPool(s string) {
	b := bufPool.Get().(*buf)
	b.bytes = append(b.bytes, s...)

	f, _ := os.Open("/dev/null")
	defer f.Close()

	fmt.Fprint(f, b)

	// put to pool
	b.bytes = b.bytes[:0]
	bufPool.Put(b)
}
