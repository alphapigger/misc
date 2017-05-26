package mutex

import (
	"testing"
)

func BenchmarkMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Mutex()
		}
	})
}

func BenchmarkRWMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RWMutex()
		}
	})
}
