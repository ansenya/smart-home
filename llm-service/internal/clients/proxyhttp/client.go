package proxyhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

// NewClient builds an *http.Client honoring proxyURL.
//
// http/https proxies route through http.Transport.Proxy (CONNECT method).
// socks5 / socks5h proxies route through a SOCKS5 dialer from x/net/proxy —
// Go's net/http does not understand socks5:// URLs in Transport.Proxy and
// would otherwise try to speak HTTP to a SOCKS server, which silently fails.
func NewClient(proxyURL string) (*http.Client, error) {
	if proxyURL == "" {
		return http.DefaultClient, nil
	}

	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}

	transport := &http.Transport{
		ForceAttemptHTTP2: true,
	}

	switch parsed.Scheme {
	case "http", "https":
		transport.Proxy = http.ProxyURL(parsed)

	case "socks5", "socks5h":
		var auth *proxy.Auth
		if parsed.User != nil {
			pw, _ := parsed.User.Password()
			auth = &proxy.Auth{User: parsed.User.Username(), Password: pw}
		}
		dialer, err := proxy.SOCKS5("tcp", parsed.Host, auth, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("socks5 dialer: %w", err)
		}
		ctxDialer, ok := dialer.(proxy.ContextDialer)
		if !ok {
			return nil, fmt.Errorf("socks5 dialer does not implement ContextDialer")
		}
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return ctxDialer.DialContext(ctx, network, addr)
		}

	default:
		return nil, fmt.Errorf("unsupported proxy scheme %q", parsed.Scheme)
	}

	return &http.Client{Transport: transport}, nil
}
