package badgerdb

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/parikshitg/urlshortener/internal/common"

	"github.com/dgraph-io/badger/v4"
)

type Store struct {
	db     *badger.DB
	expiry time.Duration
}

type Options struct {
	Path   string
	Expiry time.Duration
}

func Open(opts Options) (*Store, error) {
	bo := badger.DefaultOptions(opts.Path)
	db, err := badger.Open(bo)
	if err != nil {
		return nil, err
	}
	return &Store{db: db, expiry: opts.Expiry}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Keys
func keyCode(code string) []byte   { return []byte("code:" + code) }
func keyURL(url string) []byte     { return []byte("url:" + url) }
func keyHits(domain string) []byte { return []byte("domain_hits:" + domain) }

func (s *Store) CodeExists(code string) bool {
	if code == "" {
		return false
	}
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(keyCode(code))
		return err
	})
	return err == nil
}

func (s *Store) GetCode(url string) (string, bool) {
	if url == "" {
		return "", false
	}
	var code string
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(keyURL(url))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			code = string(val)
			return nil
		})
	})
	if err != nil {
		return "", false
	}
	return code, true
}

func (s *Store) GetURL(code string) string {
	if code == "" {
		return ""
	}
	var url string
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(keyCode(code))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			url = string(val)
			return nil
		})
	})
	if err != nil {
		return ""
	}
	return url
}

func (s *Store) Save(url, code, domain string) {
	if url == "" || code == "" || domain == "" {
		return
	}
	_ = s.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(keyURL(url), []byte(code)).WithTTL(s.expiry)
		if err := txn.SetEntry(e); err != nil {
			return err
		}
		e2 := badger.NewEntry(keyCode(code), []byte(url)).WithTTL(s.expiry)
		if err := txn.SetEntry(e2); err != nil {
			return err
		}
		// increment domain hits (no TTL)
		k := keyHits(domain)
		var count uint64
		if item, err := txn.Get(k); err == nil {
			_ = item.Value(func(val []byte) error {
				count = binary.BigEndian.Uint64(val)
				return nil
			})
		}
		count++
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, count)
		return txn.Set(k, buf)
	})
}

func (s *Store) TopDomains(n int) []common.TopN {
	if n <= 0 {
		return nil
	}
	type kv struct {
		domain string
		hits   uint64
	}
	var list []kv
	_ = s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("domain_hits:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			domain := string(bytes.TrimPrefix(item.Key(), prefix))
			var hits uint64
			_ = item.Value(func(val []byte) error {
				hits = binary.BigEndian.Uint64(val)
				return nil
			})
			list = append(list, kv{domain: domain, hits: hits})
		}
		return nil
	})
	// sort by hits desc
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[j].hits > list[i].hits {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
	if n > len(list) {
		n = len(list)
	}
	res := make([]common.TopN, n)
	for i := 0; i < n; i++ {
		res[i] = common.TopN{Rank: i + 1, Domain: list[i].domain, Shortened: int(list[i].hits)}
	}
	return res
}

func (s *Store) Purge() {
	// Badger handles TTL expiry; run value log GC opportunistically
	for i := 0; i < 2; i++ {
		if err := s.db.RunValueLogGC(0.5); err != nil {
			break
		}
	}
}
