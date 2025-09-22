package storage

type Record struct {
	URL  string
	Code string
}

type Storage interface {
	GetCode(url string) (string, bool)
	GetURL(code string) string
	Save(url, code, domain string)
	TopDomains(n int) map[string]int
}
