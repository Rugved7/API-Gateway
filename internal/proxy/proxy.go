// Package proxy handles the reverse proxy logic and http route forwarding
package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Proxy struct {
	timeout time.Duration
}

func New(timeout time.Duration) *Proxy {
	return &Proxy{
		timeout: timeout,
	}
}

func (p *Proxy) Handler(target string) (http.Handler, error) {
	upstreamURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(upstreamURL)

	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: p.timeout,
	}

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Standard forwarded headers
		if req.Header.Get("X-Forwarded-For") == "" {
			req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		}
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", "http")
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, _ error) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}

	return proxy, nil
}
