// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package uri

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//-----------------------------------------------------------------------------

func TestParse(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err bool
	}{
		// Error out on missing host:
		{in: "", err: true},
		{in: "/", err: true},
		{in: "//", err: true},

		// Test schemes:
		{in: "http://example.com", out: "http://example.com"},
		{in: "HTTP://x.example.com", out: "http://x.example.com"},
		{in: "http://localhost", out: "http://localhost"},
		{in: "http://user.local", out: "http://user.local"},
		{in: "http://kubernetes-service", out: "http://kubernetes-service"},
		{in: "https://example.com", out: "https://example.com"},
		{in: "HTTPS://example.com", out: "https://example.com"},
		{in: "ssh://example.com:22", out: "ssh://example.com:22"},
		{in: "jabber://example.com:5222", out: "jabber://example.com:5222"},

		// Leading double slashes (any scheme) defaults to http:
		{in: "//example.com", out: "http://example.com"},

		// Empty scheme defaults to http:
		{in: "localhost", out: "http://localhost"},
		{in: "LOCALHOST", out: "http://localhost"},
		{in: "localhost:80", out: "http://localhost:80"},
		{in: "localhost:8080", out: "http://localhost:8080"},
		{in: "user.local", out: "http://user.local"},
		{in: "user.local:80", out: "http://user.local:80"},
		{in: "user.local:8080", out: "http://user.local:8080"},
		{in: "kubernetes-service", out: "http://kubernetes-service"},
		{in: "kubernetes-service:80", out: "http://kubernetes-service:80"},
		{in: "kubernetes-service:8080", out: "http://kubernetes-service:8080"},
		{in: "127.0.0.1", out: "http://127.0.0.1"},
		{in: "127.0.0.1:80", out: "http://127.0.0.1:80"},
		{in: "127.0.0.1:8080", out: "http://127.0.0.1:8080"},
		{in: "[2001:db8:a0b:12f0::1]", out: "http://[2001:db8:a0b:12f0::1]"},
		{in: "[2001:db8:a0b:12f0::80]", out: "http://[2001:db8:a0b:12f0::80]"},

		// Keep the port even on matching scheme:
		{in: "http://localhost:80", out: "http://localhost:80"},
		{in: "http://localhost:8080", out: "http://localhost:8080"},
		{in: "http://x.example.io:8080", out: "http://x.example.io:8080"},
		{in: "[2001:db8:a0b:12f0::80]:80", out: "http://[2001:db8:a0b:12f0::80]:80"},
		{in: "[2001:db8:a0b:12f0::1]:8080", out: "http://[2001:db8:a0b:12f0::1]:8080"},

		// Test domains, subdomains etc.:
		{in: "example.com", out: "http://example.com"},
		{in: "1.example.com", out: "http://1.example.com"},
		{in: "1.example.io", out: "http://1.example.io"},
		{in: "subsub.sub.example.com", out: "http://subsub.sub.example.com"},
		{in: "subdomain_test.example.com", out: "http://subdomain_test.example.com"},

		// Test userinfo:
		{in: "user@example.com", out: "http://user@example.com"},
		{in: "user:passwd@example.com", out: "http://user:passwd@example.com"},
		{in: "https://user:passwd@subsub.sub.example.com", out: "https://user:passwd@subsub.sub.example.com"},

		// Lowercase scheme and host by default. Let net/url normalize URL by default:
		{in: "hTTp://subSUB.sub.EXAMPLE.COM/x//////y///foo.mp3?c=z&a=x&b=y#t=20", out: "http://subsub.sub.example.com/x//////y///foo.mp3?c=z&a=x&b=y#t=20"},

		// Some obviously wrong data:
		{in: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==", err: true},
		{in: "javascript:evilFunction()", err: true},
		{in: "otherscheme:garbage", err: true},
		{in: "<funnnytag>", err: true},

		{in: "http://www.google.com", out: "http://www.google.com"},
		{in: "https://www.google.com", out: "https://www.google.com"},
		{in: "HTTP://WWW.GOOGLE.COM", out: "http://www.google.com"},
		{in: "HTTPS://WWW.google.COM", out: "https://www.google.com"},
		{in: "http:/www.google.com", err: true},
		{in: "http:///www.google.com", err: true},
		{in: "javascript:void(0)", err: true},
		{in: "<script>", err: true},
		{in: "http:/www.google.com", err: true},
	}

	for i := range tests {
		tt := &tests[i]
		url, err := Parse(tt.in)
		if tt.err {
			require.Error(t, err)
			continue
		}
		if url.String() != tt.out {
			t.Errorf(`"%s": got "%s", want "%v"`, tt.in, url, tt.out)
		}
		// If the above defaulted to HTTP, let's test HTTPS too.
		if !strings.HasPrefix(strings.ToLower(tt.in), "http://") && strings.HasPrefix(tt.out, "http://") {
			url, err = ParseWithDefaultScheme(tt.in, "https")
			require.NoError(t, err)
			if !strings.HasPrefix(url.String(), "https://") {
				t.Errorf("%q: expected %q with https:// prefix, got %q", tt.in, tt.out, url.String())
			}
		}
	}
}

//-----------------------------------------------------------------------------

func TestSplitHostPort(t *testing.T) {
	tests := [...]struct {
		addr string
		host string
		port string
		fail error
	}{
		{"", "", "", ErrEmptyHost},
		{":", "", "", ErrEmptyHost},
		{"xxx:", "", "", ErrEmptyPort},
		{":8080", "", "", ErrEmptyHost},
		{"localhost", "localhost", "", nil},
		{"localhost:8080", "localhost", "8080", nil},
		{"[fe80::1437:feab:1dcb:c767]", "[fe80::1437:feab:1dcb:c767]", "", nil},
		{"localhost:1", "localhost", "1", nil},
		{"localhost:0", "", "", ErrPortRange},
		{"localhost:65536", "", "", ErrPortRange},
		{"127.0.0.1", "127.0.0.1", "", nil},
	}
	for i := range tests {
		d := &tests[i]
		host, port, err := SplitHostPort(d.addr)
		if d.fail != nil {
			require.Equal(t, d.fail, errors.Unwrap(err))
			continue
		}
		require.Equal(t, d.host, host, "host mismatch")
		require.Equal(t, d.port, port, "port mismatch")
	}
}
