package protocol

import "testing"

func BenchmarkFree(b *testing.B) {
	go func() {
		for i := 0; i < b.N; i++ {
			sz := NewMessage(i % 65536)
			sz.Free()
		}
	}()
	for i := 0; i < b.N; i++ {
		sz := NewMessage(i % 65536)
		sz.Free()
	}
}
