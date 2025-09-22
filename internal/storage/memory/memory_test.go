package memory

import (
	"reflect"
	"testing"
)

func TestMemStore_SaveAndGet(t *testing.T) {
	m := NewMemStore()
	url := "https://abcd.com/path"
	code := "xyz789"
	domain := "abcd.com"

	m.Save(url, code, domain)

	if c, ok := m.GetCode(url); !ok || c != code {
		t.Fatalf("expected code %q,got %q, ok=%v", code, c, ok)
	}
	if got := m.GetURL(code); got != url {
		t.Fatalf("expected url %q,got %q", url, got)
	}
}

func TestMemStore_SaveDuplicateUrls(t *testing.T) {
	m := NewMemStore()
	url := "https://abcd.com/x"
	code := "abc"
	url2 := "https://abcd.com/y"
	code2 := "def"

	m.Save(url, code, "abcd.com")
	m.Save(url, code, "abcd.com") // duplicate should not increase domain hits
	m.Save(url2, code2, "abcd.com")

	top := m.TopDomains(1)
	if len(top) != 1 {
		t.Fatalf("expected 1 top domain, got %d", len(top))
	}
	if top[0].Domain != "abcd.com" || top[0].Shortened != 2 {
		t.Fatalf("expected abcd.com with 2, got %+v", top[0])
	}
}

func TestMemStore_TopDomainsOrderingAndBounds(t *testing.T) {
	m := NewMemStore()

	// make hits: x:3, y:2, z:1
	m.Save("https://x.com/1", "x1", "x.com")
	m.Save("https://x.com/2", "x2", "x.com")
	m.Save("https://x.com/3", "x3", "x.com")
	m.Save("https://y.com/1", "y1", "y.com")
	m.Save("https://y.com/2", "y2", "y.com")
	m.Save("https://z.com/1", "z1", "z.com")

	got := m.TopDomains(5)
	expectedDomains := []string{"x.com", "y.com", "z.com"}
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
	for i, d := range expectedDomains {
		if got[i].Domain != d {
			t.Fatalf("at %d expected %s, got %s", i, d, got[i].Domain)
		}
	}

	// Request n=2
	got2 := m.TopDomains(2)
	if !reflect.DeepEqual([]string{got2[0].Domain, got2[1].Domain}, []string{"x.com", "y.com"}) {
		t.Fatalf("unexpected top2: %+v", got2)
	}
}
