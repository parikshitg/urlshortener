package ratelimiter

import (
	"testing"
	"time"
)

func TestAllowed_FirstHitConsumesToken(t *testing.T) {
	store := NewRateStore(5, 200*time.Millisecond)
	ip := "1.2.3.4"

	if !store.Allowed(ip) {
		t.Fatalf("expected first request to be allowed")
	}

	// Next 4 calls should be allowed, 6th should block within window
	for i := 0; i < 4; i++ {
		if !store.Allowed(ip) {
			t.Fatalf("expected request %d to be allowed", i+2)
		}
	}
	if store.Allowed(ip) {
		t.Fatalf("expected request to be blocked after tokens exhausted")
	}
}

func TestAllowed_ResetsAfterExpiry(t *testing.T) {
	store := NewRateStore(3, 50*time.Millisecond)
	ip := "5.6.7.8"

	// Exhaust
	for i := 0; i < 3; i++ {
		if !store.Allowed(ip) {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}
	if store.Allowed(ip) {
		t.Fatalf("expected to be blocked after exhaustion")
	}

	// Wait for expiry and confirm it resets
	time.Sleep(60 * time.Millisecond)
	if !store.Allowed(ip) {
		t.Fatalf("expected request to be allowed after window reset")
	}
}

func TestAllowed_IndependentIPs(t *testing.T) {
	store := NewRateStore(2, 200*time.Millisecond)
	ipA := "10.0.0.1"
	ipB := "10.0.0.2"

	// Consume all for A
	if !store.Allowed(ipA) || !store.Allowed(ipA) {
		t.Fatalf("expected first two requests for A to be allowed")
	}
	if store.Allowed(ipA) {
		t.Fatalf("expected third request for A to be blocked")
	}

	// B should be unaffected
	if !store.Allowed(ipB) {
		t.Fatalf("expected first request for B to be allowed")
	}
	if !store.Allowed(ipB) {
		t.Fatalf("expected second request for B to be allowed")
	}
}

func TestPurge_RemovesExpiredOnly(t *testing.T) {
	store := NewRateStore(1, 30*time.Millisecond)
	ipActive := "192.168.0.1"
	ipExpired := "192.168.0.2"

	// Create entries
	_ = store.Allowed(ipActive) // consume
	_ = store.Allowed(ipExpired)

	// Wait long enough to expire both, but refresh one
	time.Sleep(40 * time.Millisecond)
	// This should recreate/refresh ipActive
	if !store.Allowed(ipActive) {
		t.Fatalf("expected active ip to be allowed after refresh")
	}
	// Do not touch ipExpired so it remains expired

	store.Purge()

	// ipExpired should have been deleted; a call should recreate it and allow
	if !store.Allowed(ipExpired) {
		t.Fatalf("expected expired ip to be recreated and allowed after purge")
	}
}
