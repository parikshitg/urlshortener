package memory

import (
	"sort"
	"sync"
	"time"

	"github.com/parikshitg/urlshortener/internal/common"
)

type Record struct {
	Domain      string
	Code        string
	OriginalUrl string
	CreatedAt   time.Time
	Expiry      time.Time
}

// MemStore is an in memory storage unit for our service.
type MemStore struct {
	// mu is ReadWrite mutex for shared access
	mu sync.RWMutex

	expiry time.Duration

	// urlToRecord is a map of url and its shortened code
	urlToRecord map[string]Record

	// domainHits is a map of domain and number of times that domain has been shortened
	domainHits map[string]int
}

// NewMemStore creates an instance of MemStore.
func NewMemStore(expiry time.Duration) *MemStore {
	return &MemStore{
		expiry:      expiry,
		urlToRecord: make(map[string]Record),
		domainHits:  make(map[string]int),
	}
}

// GetCode takes an url and gives the corresponding unique code.
func (m *MemStore) GetCode(url string) (string, bool) {
	if url == "" {
		return "", false
	}

	m.mu.RLock()
	record, ok := m.urlToRecord[url]
	m.mu.RUnlock()

	if !ok {
		return "", false
	}

	now := time.Now()
	if now.Before(record.Expiry) {
		return record.Code, true
	}
	// expired: delete and miss
	m.mu.Lock()
	delete(m.urlToRecord, url)
	m.mu.Unlock()
	return "", false
}

// GetURL takes a code and gives corresponding original url if exists.
func (m *MemStore) GetURL(code string) string {
	if code == "" {
		return ""
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	for url, c := range m.urlToRecord {
		if c.Code == code && time.Now().Before(c.Expiry) {
			return url
		}
	}
	return ""
}

// Save saves the url, code and domain hits in memstore.
func (m *MemStore) Save(url, code, domain string) {
	if url == "" || code == "" || domain == "" {
		return // Skip invalid entries
	}

	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, exists := m.urlToRecord[url]; exists {
		if now.After(existing.Expiry) {
			// overwrite expired
			m.urlToRecord[url] = Record{
				Domain:      domain,
				Code:        code,
				OriginalUrl: url,
				CreatedAt:   now,
				Expiry:      now.Add(m.expiry),
			}
			m.domainHits[domain]++
		}
		return
	}

	m.urlToRecord[url] = Record{
		Domain:      domain,
		Code:        code,
		OriginalUrl: url,
		CreatedAt:   now,
		Expiry:      now.Add(m.expiry),
	}
	m.domainHits[domain]++
}

// TopDomains returns the top n domains based on domain hits.
func (m *MemStore) TopDomains(n int) []common.TopN {
	if n <= 0 {
		return []common.TopN{}
	}

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

	res := make([]common.TopN, n)
	for i := range res {
		res[i] = common.TopN{
			Rank:      i + 1, // rank starts from 1
			Domain:    kvs[i].domain,
			Shortened: kvs[i].hits,
		}
	}

	return res
}

func (m *MemStore) Purge() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for url, r := range m.urlToRecord {
		if now.After(r.Expiry) {
			delete(m.urlToRecord, url)
		}
	}
}

// CodeExists checks if a shortcode already exists in the storage.
func (m *MemStore) CodeExists(code string) bool {
	if code == "" {
		return false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	for _, record := range m.urlToRecord {
		if record.Code == code && now.Before(record.Expiry) {
			return true
		}
	}
	return false
}
