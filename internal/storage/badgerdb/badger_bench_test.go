package badgerdb

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func openTestStore(b *testing.B) *Store {
	dir := b.TempDir()
	st, err := Open(Options{Path: dir, Expiry: 1 * time.Hour})
	if err != nil {
		b.Fatalf("failed to open badger: %v", err)
	}
	b.Cleanup(func() {
		_ = st.Close()
		_ = os.RemoveAll(dir)
	})
	return st
}

func BenchmarkBadger_GetURL(b *testing.B) {
	st := openTestStore(b)
	st.Save("https://example.com", "abc1234", "example.com")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = st.GetURL("abc1234")
	}
}

func BenchmarkBadger_Save(b *testing.B) {
	st := openTestStore(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.Save("https://example.com", generateCode(i), "example.com")
	}
}

func generateCode(i int) string {
	return fmt.Sprintf("code%07d", i)
}
