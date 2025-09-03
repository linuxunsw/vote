package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/httprate"
	"github.com/linuxunsw/vote/backend/internal/config"
)

// Creates a new httprate ratelimiter, which uses an allowlist to determine whether
// to key by true IP or by "Real IP" headers (i.e authorised servers which handle many
// different client requests)
func newRateLimiter(opts config.RateLimitConfig, allowlist []string) func(http.Handler) http.Handler {
	return httprate.Limit(
		opts.RequestLimit,
		opts.WindowLength,
		httprate.WithKeyFuncs(allowlistKeyByRealIP(allowlist), httprate.KeyByEndpoint),
	)
}

// allowlistKeyByRealIP returns an httprate key function that will use the
// "real IP" headers (KeyByRealIP) only when the remote (true) IP is present
// in allowlist. allowlist is matched after canonicalisation.
// If allowlist is nil or empty, the returned key func always uses the true source IP.
func allowlistKeyByRealIP(allowlist []string) func(r *http.Request) (string, error) {
	// build canonicalised set for lookup
	allowed := make(map[string]struct{}, len(allowlist))
	for _, a := range allowlist {
		if a == "" {
			continue
		}
		allowed[canonicaliseIP(strings.TrimSpace(a))] = struct{}{}
	}

	return func(r *http.Request) (string, error) {
		// get the true remote IP (without port)
		remote := r.RemoteAddr
		if remote == "" {
			// fallback to remote
			return httprate.KeyByIP(r)
		}
		host, _, err := net.SplitHostPort(remote)
		if err != nil {
			// if SplitHostPort fails, try to use the whole string
			host = remote
		}
		canon := canonicaliseIP(host)

		// if allowlist is empty treat as "no trusted proxies": use true source IP
		if len(allowed) == 0 {
			return httprate.KeyByIP(r)
		}

		if _, ok := allowed[canon]; ok {
			// trusted: use RealIP-based key (this reads X-Forwarded-For, X-Real-IP, etc)
			return httprate.KeyByRealIP(r)
		}

		// not trusted: use remote addr
		return httprate.KeyByIP(r)
	}
}

// canonicaliseIP returns a form of ip suitable for comparison to other IPs.
// For IPv4 addresses, this is simply the whole string.
// For IPv6 addresses, this is the /64 prefix.
// If parsing fails, returns the original string.
func canonicaliseIP(ip string) string {
	// Fast path: detect if it contains ':' -> likely IPv6; '.' -> IPv4.
	if !strings.ContainsAny(ip, ":.") {
		return ip
	}
	// If it contains a dot but no colon it's IPv4 textual form.
	if strings.Contains(ip, ".") && !strings.Contains(ip, ":") {
		return ip
	}

	parsed := net.ParseIP(ip)
	if parsed == nil {
		return ip
	}

	// IPv4-mapped IPv6 should be normalised to IPv4 text form
	if ipv4 := parsed.To4(); ipv4 != nil {
		return ipv4.String()
	}

	// IPv6: mask to /64
	mask := net.CIDRMask(64, 128)
	return parsed.Mask(mask).String()
}
