package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const maxBodySize = 10 * 1024 * 1024

type Config struct {
	Timeout  time.Duration
	ProxyURL string
	Headers  http.Header
}

type Result struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
}

func Get(target string, config Config) (*Result, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("parse proxy URL: %w", err)
		}

		transport.Proxy = http.ProxyURL(proxyURL)
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "SONDA/1.0")

	for name, values := range config.Headers {
		if http.CanonicalHeaderKey(name) == "Host" {
			// net/http stores the outbound Host header in Request.Host.
			if len(values) > 0 {
				req.Host = values[0]
			}
			continue
		}

		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize+1))
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	duration := time.Since(start)

	if len(body) > maxBodySize {
		return nil, fmt.Errorf("response body exceeds %d-byte limit", maxBodySize)
	}

	return &Result{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Body:       body,
		Duration:   duration,
	}, nil
}
