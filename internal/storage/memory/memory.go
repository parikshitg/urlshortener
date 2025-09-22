package memory

import (
	"sort"
	"sync"
)

// MemStore is an in memory storage unit for our service.
type MemStore struct {
	// mu is ReadWrite mutex for shared access
	mu sync.RWMutex

	// urlToCode is a map of url and its shortened code
	urlToCode map[string]string

	// domainHits is a map of domain and number of times that domain has been shortened
	domainHits map[string]int
}

// NewMemStore creates an instance of MemStore.
func NewMemStore() *MemStore {
	return &MemStore{
		urlToCode:  make(map[string]string),
		domainHits: make(map[string]int),
	}
}

// GetCode takes an url and gives the corresponding unique code.
func (m *MemStore) GetCode(url string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	code, ok := m.urlToCode[url]
	return code, ok
}

// GetURL takes a code and gives corresponding original url if exists.
func (m *MemStore) GetURL(code string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for url, c := range m.urlToCode {
		if c == code {
			return url
		}
	}
	return ""
}

// Save saves the url, code and domain hits in memstore.
func (m *MemStore) Save(url, code, domain string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.urlToCode[url]; !exists {
		m.urlToCode[url] = code
		m.domainHits[domain]++
	}
}

// TopDomains returns the top n domains based on domain hits.
func (m *MemStore) TopDomains(n int) map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	type kv struct {
		domain string
		hits   int
	}

	var kvs []kv
	for domain, hits := range m.domainHits {
		kvs = append(kvs, kv{
			domain: domain,
			hits:   hits,
		})
	}

	// Sort slice by value in descending order
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].hits > kvs[j].hits
	})

	// Handle case where n > len(kvs)
	if n > len(kvs) {
		n = len(kvs)
	}

	res := make(map[string]int)
	for i := 0; i < n; i++ {
		res[kvs[i].domain] = kvs[i].hits
	}

	return res
}
