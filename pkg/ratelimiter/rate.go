package ratelimiter

import (
	"sync"
	"time"
)

type Rate struct {
	IP        string
	Tokens    int
	ExpiresAt time.Time
}

type RateStore struct {
	mu           sync.Mutex
	store        map[string]*Rate
	maxAvailable int
	expiry       time.Duration
}

func NewRateStore(maxAvailable int, expiry time.Duration) *RateStore {
	return &RateStore{
		store:        make(map[string]*Rate),
		maxAvailable: maxAvailable,
		expiry:       expiry,
	}
}

func (r *RateStore) Allowed(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	rate, ok := r.store[ip]

	if !ok || now.After(rate.ExpiresAt) {
		r.store[ip] = &Rate{
			IP:        ip,
			Tokens:    r.maxAvailable - 1,
			ExpiresAt: now.Add(r.expiry),
		}
		return true
	}

	if rate.Tokens == 0 {
		return false
	}

	rate.Tokens--
	return true
}

func (r *RateStore) Purge() {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for ip, rate := range r.store {
		if rate.ExpiresAt.Before(now) {
			delete(r.store, ip)
		}
	}
}
