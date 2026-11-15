package connector

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"sync/atomic"
	"testing"
)

func TestIsPublicAddr(t *testing.T) {
	blocked := []string{
		"127.0.0.1", "::1", "10.0.0.5", "172.16.3.4", "192.168.1.10", "169.254.169.254",
		"fe80::1", "fc00::1", "100.100.100.100", "0.0.0.0", "0.1.2.3", "::", "224.0.0.1",
		"::ffff:127.0.0.1", "::ffff:192.168.1.1", "64:ff9b::7f00:1",
	}
	for _, s := range blocked {
		if isPublicAddr(netip.MustParseAddr(s)) {
			t.Errorf("%s must be rejected", s)
		}
	}
	for _, s := range []string{"93.184.215.14", "2606:4700::1111", "::ffff:8.8.8.8"} {
		if !isPublicAddr(netip.MustParseAddr(s)) {
			t.Errorf("%s must be allowed", s)
		}
	}
}

func TestURLPreviewFetchRefusesLoopback(t *testing.T) {
	var hits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		_, _ = w.Write([]byte(`<html><head><title>internal</title></head></html>`))
	}))
	defer server.Close()
	ctx := context.Background()

	if meta := fetchPageMetadata(ctx, server.URL); len(meta) != 0 {
		t.Fatalf("metadata fetched from loopback: %v", meta)
	}
	if _, _, err := downloadURL(ctx, server.URL); !errors.Is(err, errNonPublicAddress) {
		t.Fatalf("downloadURL error = %v, want errNonPublicAddress", err)
	}
	if n := hits.Load(); n != 0 {
		t.Fatalf("loopback server received %d requests", n)
	}
}
