package storage

// Storage is an adapter interface, that defines the methods for our services
// storage logic.
type Storage interface {
	GetCode(url string) (string, bool)
	GetURL(code string) string
	Save(url, code, domain string)
	TopDomains(n int) map[string]int
}
