package badgerdb

import (
	"os"
	"testing"
	"time"
)

func withStore(t *testing.T, expiry time.Duration, fn func(*Store)) {
	t.Helper()
	dir := t.TempDir()
	st, err := Open(Options{Path: dir, Expiry: expiry})
	if err != nil {
		t.Fatalf("failed to open badger: %v", err)
	}
	t.Cleanup(func() {
		_ = st.Close()
		_ = os.RemoveAll(dir)
	})
	fn(st)
}

func TestBadger_SaveGetResolve(t *testing.T) {
	withStore(t, 1*time.Hour, func(st *Store) {
		url := "https://example.com"
		code := "abc123"
		domain := "example.com"

		st.Save(url, code, domain)

		if got, ok := st.GetCode(url); !ok || got != code {
			t.Fatalf("GetCode: want %q ok=true, got %q ok=%v", code, got, ok)
		}
		if got := st.GetURL(code); got != url {
			t.Fatalf("GetURL: want %q, got %q", url, got)
		}
	})
}

func TestBadger_Expiry(t *testing.T) {
	withStore(t, 1*time.Second, func(st *Store) {
		url := "https://example.com"
		code := "abc123"
		domain := "example.com"
		st.Save(url, code, domain)
		// Initially present
		if _, ok := st.GetCode(url); !ok {
			t.Fatalf("expected code to exist")
		}
		if st.GetURL(code) == "" {
			t.Fatalf("expected url to exist")
		}
		// Wait for TTL
		time.Sleep(1200 * time.Millisecond)
		if _, ok := st.GetCode(url); ok {
			t.Fatalf("expected code to expire")
		}
		if st.GetURL(code) != "" {
			t.Fatalf("expected url to expire")
		}
	})
}

func TestBadger_TopDomains(t *testing.T) {
	withStore(t, 1*time.Hour, func(st *Store) {
		st.Save("https://a.com", "a1", "a.com")
		st.Save("https://a.com/x", "a2", "a.com")
		st.Save("https://b.com", "b1", "b.com")

		got := st.TopDomains(2)
		if len(got) != 2 {
			t.Fatalf("expected 2 results, got %d", len(got))
		}
		if got[0].Domain != "a.com" || got[0].Shortened != 2 {
			t.Fatalf("unexpected top[0]: %+v", got[0])
		}
		if got[1].Domain != "b.com" || got[1].Shortened != 1 {
			t.Fatalf("unexpected top[1]: %+v", got[1])
		}
	})
}

func TestBadger_CodeExists(t *testing.T) {
	withStore(t, 1*time.Hour, func(st *Store) {
		if st.CodeExists("nope") {
			t.Fatalf("expected false for non-existent code")
		}
		st.Save("https://x.com", "xy1", "x.com")
		if !st.CodeExists("xy1") {
			t.Fatalf("expected true after save")
		}
	})
}

func TestBadger_GCDoesNotPanic(t *testing.T) {
	withStore(t, 500*time.Millisecond, func(st *Store) {
		st.Save("https://gc.com", "gc1", "gc.com")
		time.Sleep(600 * time.Millisecond)
		st.Purge() // run GC; should not panic
	})
}
