package httpclient

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("User-Agent") != "SONDA/1.0" {
				t.Errorf(
					"unexpected User-Agent: %q",
					r.Header.Get("User-Agent"),
				)
			}

			if r.Header.Get("X-Sonda-Test") != "enabled" {
				t.Errorf(
					"custom header was not received",
				)
			}

			if r.Host != "virtual.example.test" {
				t.Errorf("unexpected Host header: %q", r.Host)
			}

			w.Header().Set("X-Sonda-Response", "true")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write([]byte("hello SONDA"))
		}),
	)

	defer server.Close()

	config := Config{
		Timeout: 5 * time.Second,
		Headers: http.Header{
			"X-Sonda-Test": []string{"enabled"},
			"Host":         []string{"virtual.example.test"},
		},
	}

	result, err := Get(server.URL, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			result.StatusCode,
		)
	}

	if string(result.Body) != "hello SONDA" {
		t.Fatalf(
			"unexpected body: %q",
			string(result.Body),
		)
	}

	if result.Headers.Get("X-Sonda-Response") != "true" {
		t.Fatal("response header was not captured")
	}
}

func TestGetIncludesBodyTransferInDuration(t *testing.T) {
	const delay = 100 * time.Millisecond

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.(http.Flusher).Flush()
			time.Sleep(delay)
			_, _ = w.Write([]byte("delayed body"))
		}),
	)
	defer server.Close()

	result, err := Get(server.URL, Config{Timeout: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Duration < delay/2 {
		t.Fatalf(
			"duration %s did not include the delayed response body",
			result.Duration,
		)
	}
}

func TestGetRejectsOversizedResponseBody(t *testing.T) {
	body := bytes.Repeat([]byte("a"), maxBodySize+1)

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(body)
		}),
	)
	defer server.Close()

	_, err := Get(server.URL, Config{Timeout: 5 * time.Second})
	if err == nil {
		t.Fatal("expected an error for an oversized response body")
	}

	if !strings.Contains(err.Error(), "response body exceeds") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetUsesProxy(t *testing.T) {
	const target = "http://example.test/search?q=admin"
	var proxyUsed atomic.Bool

	proxy := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxyUsed.Store(true)

			if r.URL.String() != target {
				t.Errorf("proxy received unexpected URL: %q", r.URL.String())
			}

			_, _ = w.Write([]byte("proxied response"))
		}),
	)
	defer proxy.Close()

	result, err := Get(target, Config{
		Timeout:  time.Second,
		ProxyURL: proxy.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !proxyUsed.Load() {
		t.Fatal("request did not pass through the configured proxy")
	}

	if string(result.Body) != "proxied response" {
		t.Fatalf("unexpected body: %q", result.Body)
	}
}
