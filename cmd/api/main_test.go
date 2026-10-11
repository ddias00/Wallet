package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHealth(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusServiceUnavailable} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/health" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(status)
			}))
			defer server.Close()
			if err := checkHealth(server.URL + "/health"); (err != nil) != (status != http.StatusOK) {
				t.Fatalf("status %d: unexpected error %v", status, err)
			}
		})
	}
}

func TestCheckHealthUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	server.Close()
	if err := checkHealth(server.URL + "/health"); err == nil {
		t.Fatal("expected an error when the server is unavailable")
	}
}
