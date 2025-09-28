package storage

import (
	"github.com/parikshitg/urlshortener/internal/common"
)

// Storage is an adapter interface, that defines the methods for our services
// storage logic.
type Storage interface {
	// CodeExists checks if a shortcode already exists in the storage.
	CodeExists(code string) bool

	// GetCode takes an url and gives the corresponding unique code.
	GetCode(url string) (string, bool)

	// GetURL takes a code and gives corresponding original url if exists.
	GetURL(code string) string

	// Save saves the url, code and domain hits in memstore.
	Save(url, code, domain string)

	// TopDomains returns the top n domains based on domain hits.
	TopDomains(n int) []common.TopN

	// Purge deletes the expired records
	Purge()
}
