// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package htp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"unsafe"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv/htp/hdr"
)

type (
	Header = http.Header
	Status = hdr.Status
)

type Response struct {
	res  http.ResponseWriter
	Req  *http.Request
	Ctx  context.Context
	hdr  Header
	mime string // response content type
	code int    // response HTTP code
	clen int64  // response content length
}

//-----------------------------------------------------------------------------
// Response Pool
//-----------------------------------------------------------------------------

var rpool = sync.Pool{New: func() interface{} {
	return new(Response)
}}

var rzero Response

//-----------------------------------------------------------------------------

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
	if r.code == 0 {
		r.code = statusCode
		r.res.WriteHeader(statusCode)
	}
}

//-----------------------------------------------------------------------------

func (r *Response) Set(header string, value string) {
	r.hdr.Set(header, value)
}

func (r *Response) SetContentType(mime string) {
	r.mime = mime
	r.hdr.Set(hdr.ContentType, mime)
}

func (r *Response) SetContentLength(n int64) {
	r.clen = n
	r.hdr.Set(hdr.ContentLength, strconv.FormatInt(n, 10))
}

//-----------------------------------------------------------------------------
// Output
//-----------------------------------------------------------------------------

func (r *Response) Write(p []byte) (int, error) {
	return r.res.Write(p)
}

func (r *Response) WriteString(s string) (int, error) {
	return r.Write(unsafeBytes(s))
}

func (r *Response) Flush() {
	if fl, ok := r.res.(http.Flusher); ok {
		fl.Flush()
	}
}

//-----------------------------------------------------------------------------
// Content-Type
//-----------------------------------------------------------------------------

func (r *Response) SetText() {
	r.SetContentType(hdr.TextPlain)
}

func (r *Response) SetJSON() {
	r.SetContentType(hdr.AppJSON)
}

//-----------------------------------------------------------------------------
// Send
//-----------------------------------------------------------------------------

func (r *Response) SendBytes(b []byte) error {
	r.SetContentLength(int64(len(b)))
	_, err := r.res.Write(b)
	return err
}

//-----------------------------------------------------------------------------
// Errors
//-----------------------------------------------------------------------------

func (r *Response) Error(code Status, text ...string) bool {
	if r.code != 0 {
		return false
	}
	if len(text) == 0 {
		text = append(text, code.String())
	}
	r.errorHeader(code, strLen(text))
	var err error
	for i := 0; i < len(text) && err == nil; i++ {
		_, err = r.WriteString(text[i])
	}
	return err == nil
}

//-----------------------------------------------------------------------------

func (r *Response) Errorf(code Status, format string, args ...interface{}) bool {
	if r.code != 0 {
		return false
	}
	if len(args) == 0 {
		return r.Error(code, format)
	}
	r.errorHeader(code, 0)
	_, err := fmt.Fprintf(r, format, args...)
	return err == nil
}

//-----------------------------------------------------------------------------

func (r *Response) errorHeader(code Status, clen int) {
	r.SetText()
	r.Set(hdr.XContentTypeOptions, hdr.Nosniff)
	if clen > 0 {
		r.SetContentLength(int64(clen))
	}
	r.WriteHeader(int(code))
}

//-----------------------------------------------------------------------------
// Common Errors
//-----------------------------------------------------------------------------

func (r *Response) ServerError() bool {
	return r.Error(hdr.StatusInternalServerError)
}

//-----------------------------------------------------------------------------
// Unsafe Bytes
//-----------------------------------------------------------------------------

type slice struct {
	s string
	c int
}

func unsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&slice{s, len(s)}))
}

//-----------------------------------------------------------------------------
// Helpers
//-----------------------------------------------------------------------------

func strLen(strings []string) (n int) {
	for i := 0; i < len(strings); i++ {
		n += len(strings[i])
	}
	return
}
