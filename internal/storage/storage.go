package storage

import (
	"github.com/parikshitg/urlshortner/internal/common"
)

// Storage is an adapter interface, that defines the methods for our services
// storage logic.
type Storage interface {
	// GetCode takes an url and gives the corresponding unique code.
	GetCode(url string) (string, bool)

	// GetURL takes a code and gives corresponding original url if exists.
	GetURL(code string) string

	// Save saves the url, code and domain hits in memstore.
	Save(url, code, domain string)

	// TopDomains returns the top n domains based on domain hits.
	TopDomains(n int) []common.TopN
}
