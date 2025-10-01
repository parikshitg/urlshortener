package memory

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkMem_GetURL(b *testing.B) {
	store := NewMemStore(1 * time.Hour)
	store.Save("https://example.com", "abc1234", "example.com")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.GetURL("abc1234")
	}
}

func BenchmarkMem_Save(b *testing.B) {
	store := NewMemStore(1 * time.Hour)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		code := fmt.Sprintf("code%07d", i)
		store.Save("https://example.com", code, "example.com")
	}
}
