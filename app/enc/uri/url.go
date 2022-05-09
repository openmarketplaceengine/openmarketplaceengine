// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package uri

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	minPort = 1
	maxPort = 65535
)

var (
	ErrEmptyHost = errors.New("empty host")
	ErrCheckHost = errors.New("invalid host address")
	ErrParsePort = errors.New("invalid port number format")
	ErrEmptyPort = errors.New("empty network port")
	ErrPortRange = errors.New("network port is out of range [1...65535]")
)

//-----------------------------------------------------------------------------

func Parse(rawURL string) (*url.URL, error) {
	return ParseWithDefaultScheme(rawURL, "http")
}

//-----------------------------------------------------------------------------

func ParseWithDefaultScheme(rawURL string, scheme string) (*url.URL, error) {
	rawURL = defaultScheme(rawURL, scheme)

	u, err := url.Parse(rawURL)

	if err != nil {
		return nil, err
	}

	addr := strings.ToLower(u.Host)

	var host string

	host, _, err = SplitHostPort(addr)

	if err != nil {
		return nil, err
	}

	if err := checkHost(host); err != nil {
		return nil, err
	}

	u.Host = addr

	if len(u.Path) == 1 && u.Path[0] == '/' {
		u.Path = ""
	}

	return u, nil
}

// SplitHostPort splits network address of the form "host:port" into
// host and port. Unlike net.SplitHostPort(), it doesn't remove brackets
// from [IPv6] host.
func SplitHostPort(addr string) (host, port string, err error) {
	n := len(addr)

	host = addr

loop:
	for n > 0 {
		c := addr[n-1]
		switch c {
		case ':':
			if n == len(addr) {
				err = parseError(addr, ErrEmptyPort)
				return
			}
			host = addr[:n-1]
			port = addr[n:]
			break loop
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			n--
		default:
			break loop
		}
	}

	if len(host) == 0 {
		err = parseError(addr, ErrEmptyHost)
		return
	}

	if len(port) == 0 {
		return
	}

	var num int

	if num, err = strconv.Atoi(port); err != nil {
		err = parseError(addr, ErrParsePort)
	} else if num < minPort || num > maxPort {
		err = parseError(addr, ErrPortRange)
	}

	return //nolint
}

//-----------------------------------------------------------------------------

func parseError(URL string, err error) error {
	return &url.Error{Op: "error parsing URL", URL: URL, Err: err}
}

//-----------------------------------------------------------------------------

func defaultScheme(rawURL, scheme string) string {
	const schemeDelim = "://"
	// Force default http scheme, so net/url.Parse() doesn't
	// put both host and path into the (relative) path.
	if n := len(rawURL); n > 1 {
		if rawURL[0] == '/' && rawURL[1] == '/' {
			return scheme + ":" + rawURL
		}
		if n > 3 && strings.Contains(rawURL, schemeDelim) {
			return rawURL
		}
	}
	return scheme + schemeDelim + rawURL
}

//-----------------------------------------------------------------------------

var (
	hostRegexp = regexp.MustCompile(`^([a-zA-Z0-9-_]{1,63}\.)*([a-zA-Z0-9-]{1,63})$`)
	ipv4Regexp = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	ipv6Regexp = regexp.MustCompile(`^\[[a-fA-F0-9:]+]$`)
)

func checkHost(host string) error {
	if len(host) == 0 {
		return parseError("", ErrEmptyHost)
	}

	if hostRegexp.MatchString(host) {
		return nil
	}

	if ipv4Regexp.MatchString(host) || ipv6Regexp.MatchString(host) {
		return nil
	}

	return parseError(host, ErrCheckHost)
}
