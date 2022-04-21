// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package hdr

import (
	"net/http"
	"strconv"
)

type Status int

//goland:noinspection GoUnusedConst
const (
	StatusContinue           Status = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols Status = 101 // RFC 7231, 6.2.2
	StatusProcessing         Status = 102 // RFC 2518, 10.1
)

//goland:noinspection GoUnusedConst
const (
	StatusOK                   Status = 200 // RFC 7231, 6.3.1
	StatusCreated              Status = 201 // RFC 7231, 6.3.2
	StatusAccepted             Status = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo Status = 203 // RFC 7231, 6.3.4
	StatusNoContent            Status = 204 // RFC 7231, 6.3.5
	StatusResetContent         Status = 205 // RFC 7231, 6.3.6
	StatusPartialContent       Status = 206 // RFC 7233, 4.1
	StatusMultiStatus          Status = 207 // RFC 4918, 11.1
	StatusAlreadyReported      Status = 208 // RFC 5842, 7.1
	StatusIMUsed               Status = 226 // RFC 3229, 10.4.1
)

//goland:noinspection GoUnusedConst
const (
	StatusMultipleChoices   Status = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  Status = 301 // RFC 7231, 6.4.2
	StatusFound             Status = 302 // RFC 7231, 6.4.3
	StatusSeeOther          Status = 303 // RFC 7231, 6.4.4
	StatusNotModified       Status = 304 // RFC 7232, 4.1
	StatusUseProxy          Status = 305 // RFC 7231, 6.4.5
	_                       Status = 306 // RFC 7231, 6.4.6 (Unused)
	StatusTemporaryRedirect Status = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect Status = 308 // RFC 7538, 3
)

//goland:noinspection GoUnusedConst
const (
	StatusBadRequest                   Status = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 Status = 401 // RFC 7235, 3.1
	StatusPaymentRequired              Status = 402 // RFC 7231, 6.5.2
	StatusForbidden                    Status = 403 // RFC 7231, 6.5.3
	StatusNotFound                     Status = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             Status = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                Status = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            Status = 407 // RFC 7235, 3.2
	StatusRequestTimeout               Status = 408 // RFC 7231, 6.5.7
	StatusConflict                     Status = 409 // RFC 7231, 6.5.8
	StatusGone                         Status = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               Status = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           Status = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        Status = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            Status = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         Status = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable Status = 416 // RFC 7233, 4.4
	StatusExpectationFailed            Status = 417 // RFC 7231, 6.5.14
	StatusTeapot                       Status = 418 // RFC 7168, 2.3.3
	StatusUnprocessableEntity          Status = 422 // RFC 4918, 11.2
	StatusLocked                       Status = 423 // RFC 4918, 11.3
	StatusFailedDependency             Status = 424 // RFC 4918, 11.4
	StatusUpgradeRequired              Status = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         Status = 428 // RFC 6585, 3
	StatusTooManyRequests              Status = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  Status = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   Status = 451 // RFC 7725, 3
)

//goland:noinspection GoUnusedConst
const (
	StatusInternalServerError           Status = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                Status = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    Status = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            Status = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                Status = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       Status = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         Status = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           Status = 507 // RFC 4918, 11.5
	StatusLoopDetected                  Status = 508 // RFC 5842, 7.2
	StatusNotExtended                   Status = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired Status = 511 // RFC 6585, 6
)

//-----------------------------------------------------------------------------

func (s Status) String() string {
	text := http.StatusText(int(s))
	if len(text) > 0 {
		return text
	}
	return strconv.Itoa(int(s))
}
