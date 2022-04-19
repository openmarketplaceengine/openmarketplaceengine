// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package htp

import (
	"context"
	"net/http"
	"strconv"
	"sync"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv/htp/hdr"
)

type (
	Header = http.Header
)

type Response struct {
	res http.ResponseWriter
	Req *http.Request
	Ctx context.Context
	hdr Header
}

//-----------------------------------------------------------------------------
// Response Pool
//-----------------------------------------------------------------------------

var rpool = sync.Pool{New: func() interface{} {
	return new(Response)
}}

var rzero Response

func GetRes(res http.ResponseWriter, req *http.Request) *Response {
	r, _ := rpool.Get().(*Response)
	return r.Init(res, req)
}

func (r *Response) Init(res http.ResponseWriter, req *http.Request) *Response {
	r.res = res
	r.Req = req
	r.Ctx = req.Context()
	r.hdr = res.Header()
	return r
}

func (r *Response) Release() {
	*r = rzero
	rpool.Put(r)
}

//-----------------------------------------------------------------------------
// Context
//-----------------------------------------------------------------------------

func (r *Response) Done() bool {
	select {
	case <-r.Ctx.Done():
		return true
	default:
		return false
	}
}

//-----------------------------------------------------------------------------
// Headers
//-----------------------------------------------------------------------------

func (r *Response) Header() Header {
	return r.hdr
}

func (r *Response) WriteHeader(statusCode int) {
	r.res.WriteHeader(statusCode)
}

func (r *Response) SetContentLength(n int64) {
	r.hdr.Set(hdr.ContentLength, strconv.FormatInt(n, 10))
}

//-----------------------------------------------------------------------------
// Output
//-----------------------------------------------------------------------------

func (r *Response) Write(p []byte) (int, error) {
	return r.res.Write(p)
}

func (r *Response) Flush() {
	if fl, ok := r.res.(http.Flusher); ok {
		fl.Flush()
	}
}

//-----------------------------------------------------------------------------
// Content-Type
//-----------------------------------------------------------------------------

func (r *Response) SetContentType(ctype string) {
	r.hdr.Set(hdr.ContentType, ctype)
}

func (r *Response) SetText() {
	r.SetContentType(hdr.TextPlain)
}

func (r *Response) SetJSON() {
	r.SetContentType(hdr.AppJSON)
}
