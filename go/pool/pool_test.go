package pool

import (
	"testing"
)

func BenchmarkWriteString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WriteString("hello")
	}
}

func BenchmarkWriteStringWithPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WriteStringWithPool("hello")
	}
}
