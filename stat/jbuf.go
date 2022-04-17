package stat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"
)

const (
	MaxIndent = 128
	defBufLen = 1024
)

type jbtoken int

const (
	jbStart jbtoken = iota
	jbObject
	jbObjEnd
	jbObjKey
	jbArray
	jbAryEnd
	jbValue
	jbComma
)

type JSONBuffer struct {
	buf []byte  // buffer
	tok jbtoken // current token
	pos int     // current token pos
	lev int     // nested object/array level
	ind int     // indent level
	inl int     // inline level
}

type JSONAppender interface {
	AppendJSON(b []byte) ([]byte, error)
}

type JSONWriter interface {
	WriteJSON(b *JSONBuffer) error
}

//-----------------------------------------------------------------------------

func (b *JSONBuffer) Len() int {
	return len(b.buf)
}

func (b *JSONBuffer) Cap() int {
	return cap(b.buf)
}

func (b *JSONBuffer) Buf() []byte {
	return b.buf
}

func (b *JSONBuffer) Cut(pos int) {
	b.buf = b.buf[:pos]
}

func (b *JSONBuffer) End() []byte {
	if b.ind > 0 {
		b.NewLine()
	}
	return b.buf
}

func (b *JSONBuffer) UnsafeString() string {
	return *(*string)(unsafe.Pointer(&b.buf))
}

//-----------------------------------------------------------------------------
// Pool Management
//-----------------------------------------------------------------------------

var _jsonBufferPool = sync.Pool{New: func() interface{} {
	return &JSONBuffer{
		buf: make([]byte, 0, defBufLen),
	}
}}

func AcquireJSONBuffer(indent int) *JSONBuffer {
	buf, _ := _jsonBufferPool.Get().(*JSONBuffer)
	return buf.WithIndent(indent)
}

func (b *JSONBuffer) Reset() {
	if len(b.buf) > 0 {
		b.buf = b.buf[:0]
	}
	b.tok = jbStart
	b.pos = 0
	b.lev = 0
	b.ind = 0
	b.inl = 0
}

func (b *JSONBuffer) Release() {
	b.Reset()
	_jsonBufferPool.Put(b)
}

//-----------------------------------------------------------------------------
// Indent
//-----------------------------------------------------------------------------

func (b *JSONBuffer) WithIndent(indent int) *JSONBuffer {
	switch {
	case indent <= 0:
		b.ind = 0
	case indent > MaxIndent:
		b.ind = MaxIndent
	default:
		b.ind = indent
	}
	return b
}

// Indent returns current indent level.
func (b *JSONBuffer) Indent() int {
	if b.inl == 0 {
		return b.ind * b.lev
	}
	return 0
}

// ConfigIndent returns JSONBuffer indent value.
func (b *JSONBuffer) ConfigIndent() int {
	return b.ind
}

func (b *JSONBuffer) ShouldIndent() bool {
	return b.ind > 0 && b.inl == 0
}

func (b *JSONBuffer) IgnoreIndent() bool {
	return b.ind == 0 || b.inl > 0
}

//-----------------------------------------------------------------------------
// Inline
//-----------------------------------------------------------------------------

func (b *JSONBuffer) Inline(inline bool) {
	if inline {
		b.inl++
	} else if b.inl > 0 {
		b.inl--
	}
}

//-----------------------------------------------------------------------------
// Token
//-----------------------------------------------------------------------------

func (b *JSONBuffer) setval() {
	b.tok = jbValue
	b.pos = len(b.buf)
}

func (b *JSONBuffer) setpos(tok jbtoken) {
	b.tok = tok
	b.pos = len(b.buf)
}

func (b *JSONBuffer) srctok() jbtoken { //nolint
	for i := len(b.buf); i > 0; i-- {
		b.pos = i
		switch b.buf[i-1] {
		case ' ', '\n', '\r', '\t':
			// ignore
		case '}':
			return jbObjEnd
		case ']':
			return jbAryEnd
		case ',':
			return jbComma
		case '{':
			return jbObject
		case '[':
			return jbArray
		case ':':
			return jbObjKey
		default:
			return jbValue
		}
	}
	b.pos = 0
	return jbStart
}

//-----------------------------------------------------------------------------
// Basic Types
//-----------------------------------------------------------------------------

func (b *JSONBuffer) Null() {
	b.buf = append(b.buf, 'n', 'u', 'l', 'l')
	b.setval()
}

//-----------------------------------------------------------------------------

func (b *JSONBuffer) True() {
	b.buf = append(b.buf, 't', 'r', 'u', 'e')
	b.setval()
}

//-----------------------------------------------------------------------------

func (b *JSONBuffer) False() {
	b.buf = append(b.buf, 'f', 'a', 'l', 's', 'e')
	b.setval()
}

//-----------------------------------------------------------------------------

func (b *JSONBuffer) Bool(v bool) {
	if v {
		b.True()
		return
	}
	b.False()
}

//-----------------------------------------------------------------------------
// JSON String
//-----------------------------------------------------------------------------

func (b *JSONBuffer) String(s string, safeHTML bool) {
	b.buf = append(b.buf, '"')
	b.buf = escape(b.buf, s, safeHTML)
	b.buf = append(b.buf, '"')
	b.setval()
}

func (b *JSONBuffer) Stringf(format string, args ...interface{}) {
	b.String(fmt.Sprintf(format, args...), true)
}

func (b *JSONBuffer) BinaryString(s []byte, safeHTML bool) {
	if s == nil {
		b.Null()
		return
	}
	b.buf = append(b.buf, '"')
	v := *(*string)(unsafe.Pointer(&b))
	b.buf = escape(b.buf, v, safeHTML)
	b.buf = append(b.buf, '"')
	b.setval()
}

func (b *JSONBuffer) RawString(s string) {
	b.buf = append(b.buf, '"')
	b.buf = append(b.buf, s...)
	b.buf = append(b.buf, '"')
	b.setval()
}

//-----------------------------------------------------------------------------
// Append
//-----------------------------------------------------------------------------

func (b *JSONBuffer) AppendFrom(a JSONAppender) error {
	return b.AppendFunc(a.AppendJSON)
}

func (b *JSONBuffer) AppendFunc(f func(b []byte) ([]byte, error)) error {
	buf, err := f(b.buf)
	if err == nil && len(buf) > len(b.buf) {
		b.buf = buf
		b.srctok()
	}
	return err
}

func (b *JSONBuffer) AppendJSON(src []byte) error {
	if src == nil {
		b.Null()
		return nil
	}
	var pfx, ind string
	if b.ShouldIndent() {
		pfx = whiteSpace(b.Indent())
		ind = whiteSpace(b.ConfigIndent())
	}
	dst := AcquireByteBuffer(len(src) + ((len(pfx) + len(ind)) * 32))
	defer ReleaseByteBuffer(dst)
	err := json.Indent(dst, src, pfx, ind)
	if err == nil {
		b.buf = append(b.buf, dst.Bytes()...)
		b.srctok()
	}
	return err
}

//-----------------------------------------------------------------------------
// Encode
//-----------------------------------------------------------------------------

func (b *JSONBuffer) EncodeJSON(v interface{}) error {
	if v == nil {
		b.Null()
		return nil
	}
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.AppendJSON(buf)
}

//-----------------------------------------------------------------------------

func (b *JSONBuffer) WriteFrom(w JSONWriter) error {
	err := w.WriteJSON(b)
	if err == nil {
		b.srctok()
	}
	return err
}

//-----------------------------------------------------------------------------
// interface{}
//-----------------------------------------------------------------------------

func (b *JSONBuffer) Value(i interface{}) error {
	switch v := i.(type) {
	case nil:
		b.Null()
	case JSONWriter:
		return b.WriteFrom(v)
	case JSONAppender:
		return b.AppendFrom(v)
	case string:
		b.String(v, true)
	case int64:
		b.Int64(v)
	case int:
		b.Int(v)
	case uint64:
		b.Uint64(v)
	case uint:
		b.Uint(v)
	case time.Time:
		b.Time(v)
	case time.Duration:
		b.Seconds(v)
	case json.RawMessage:
		return b.AppendJSON(v)
	case bool:
		b.Bool(v)
	default:
		return b.EncodeJSON(i)
	}
	return nil
}

//-----------------------------------------------------------------------------
// start / close
//-----------------------------------------------------------------------------

func (b *JSONBuffer) start(c byte, start jbtoken) {
	b.lev++
	b.tok = start
	if b.IgnoreIndent() {
		b.buf = append(b.buf, c)
		b.pos = len(b.buf)
		return
	}
	b.buf = append(b.buf, c, '\n')
	b.pos = len(b.buf) - 1
	b.indent()
}

func (b *JSONBuffer) close(c byte, start, endtok jbtoken) {
	ind := b.ShouldIndent()
	switch b.tok {
	case start:
		if ind {
			b.Cut(b.pos)
		}
	case jbComma:
		b.Cut(b.pos - 1)
	}
	b.lev--
	if ind {
		b.eolind()
	}
	b.buf = append(b.buf, c)
	b.pos = len(b.buf)
	b.tok = endtok
}

//-----------------------------------------------------------------------------
// Object
//-----------------------------------------------------------------------------

func (b *JSONBuffer) ObjectStart() {
	b.start('{', jbObject)
}

func (b *JSONBuffer) ObjectClose() {
	b.close('}', jbObject, jbObjEnd)
}

func (b *JSONBuffer) EmptyObject() {
	b.buf = append(b.buf, '{', '}')
	b.setpos(jbObjEnd)
}

//-----------------------------------------------------------------------------
// Object Key
//-----------------------------------------------------------------------------

func (b *JSONBuffer) KeyEsc(key string) {
	b.Key(key, true)
}

func (b *JSONBuffer) KeyRaw(key string) {
	b.Key(key, false)
}

func (b *JSONBuffer) Key(key string, esc bool) {
	b.tok = jbObjKey
	b.buf = append(b.buf, '"')
	if esc {
		b.buf = escape(b.buf, key, true)
	} else {
		b.buf = append(b.buf, key...)
	}
	if b.ind > 0 {
		b.buf = append(b.buf, '"', ':', ' ')
		b.pos = len(b.buf) - 1
	} else {
		b.buf = append(b.buf, '"', ':')
		b.pos = len(b.buf)
	}
}

//-----------------------------------------------------------------------------
// Array
//-----------------------------------------------------------------------------

func (b *JSONBuffer) ArrayStart() {
	b.start('[', jbArray)
}

func (b *JSONBuffer) ArrayClose() {
	b.close(']', jbArray, jbAryEnd)
}

func (b *JSONBuffer) EmptyArray() {
	b.buf = append(b.buf, '[', ']')
	b.setpos(jbAryEnd)
}

//-----------------------------------------------------------------------------
// Comma
//-----------------------------------------------------------------------------

func (b *JSONBuffer) Comma() {
	switch {
	case b.tok == jbComma:
		return
	case b.ind == 0:
		b.buf = append(b.buf, ',')
		b.pos = len(b.buf)
	case b.inl > 0:
		b.buf = append(b.buf, ',', ' ')
		b.pos = len(b.buf) - 1
	default:
		b.buf = append(b.buf, ',', '\n')
		b.pos = len(b.buf) - 1
		b.indent()
	}
	b.tok = jbComma
}

//-----------------------------------------------------------------------------
// Char / Rune
//-----------------------------------------------------------------------------

func (b *JSONBuffer) Char(c byte) {
	b.buf = append(b.buf, c)
}

func (b *JSONBuffer) Char2(c1, c2 byte) {
	b.buf = append(b.buf, c1, c2)
}

func (b *JSONBuffer) Chars(chars ...byte) {
	b.buf = append(b.buf, chars...)
}

func (b *JSONBuffer) Rune(r rune) {
	if r < utf8.RuneSelf {
		b.buf = append(b.buf, byte(r))
		return
	}
	b.runeMust(r)
}

func (b *JSONBuffer) runeMust(r rune) {
	n := len(b.buf)
	b.buf = append(b.buf, ' ', ' ', ' ', ' ')
	n += utf8.EncodeRune(b.buf[n:], r)
	b.buf = b.buf[:n]
}

//-----------------------------------------------------------------------------
// Formatting
//-----------------------------------------------------------------------------

func (b *JSONBuffer) NewLine() {
	if n := len(b.buf); n > 0 && b.buf[n-1] != '\n' {
		b.buf = append(b.buf, '\n')
	}
}

func (b *JSONBuffer) Space(n int) {
	const space = "                                                            "
	const splen = len(space)
	if n <= splen {
		b.buf = append(b.buf, space[:n]...)
		return
	}
	for n >= splen {
		b.buf = append(b.buf, space...)
		n -= splen
	}
	if n > 0 {
		b.buf = append(b.buf, space[:n]...)
	}
}

func (b *JSONBuffer) indent() {
	b.Space(b.ind * b.lev)
}

func (b *JSONBuffer) eolind() {
	b.buf = append(b.buf, '\n')
	b.Space(b.ind * b.lev)
}

//-----------------------------------------------------------------------------
// Int
//-----------------------------------------------------------------------------

func (b *JSONBuffer) Int64(i int64) {
	b.buf = strconv.AppendInt(b.buf, i, 10)
	b.setval()
}

func (b *JSONBuffer) Int32(i int32) {
	b.Int64(int64(i))
}

func (b *JSONBuffer) Int16(i int16) {
	b.Int64(int64(i))
}

func (b *JSONBuffer) Int8(i int8) {
	b.Int64(int64(i))
}

func (b *JSONBuffer) Int(i int) {
	b.Int64(int64(i))
}

//-----------------------------------------------------------------------------
// Uint
//-----------------------------------------------------------------------------

func (b *JSONBuffer) Uint64(i uint64) {
	b.buf = strconv.AppendUint(b.buf, i, 10)
	b.setval()
}

func (b *JSONBuffer) Uint32(i uint32) {
	b.Uint64(uint64(i))
}

func (b *JSONBuffer) Uint16(i uint16) {
	b.Uint64(uint64(i))
}

func (b *JSONBuffer) Uint8(i uint8) {
	b.Uint64(uint64(i))
}

func (b *JSONBuffer) Uint(i uint) {
	b.Uint64(uint64(i))
}

func (b *JSONBuffer) Uintptr(i uintptr) {
	b.Uint64(uint64(i))
}

//-----------------------------------------------------------------------------
// Time
//-----------------------------------------------------------------------------

func (b *JSONBuffer) Time(t time.Time) {
	b.RFC3339Time(t)
}

func (b *JSONBuffer) UnixTime(t time.Time) {
	b.Int64(t.Unix())
}

func (b *JSONBuffer) RFC3339Time(t time.Time) {
	b.buf = append(b.buf, '"')
	b.buf = t.AppendFormat(b.buf, time.RFC3339Nano)
	b.buf = append(b.buf, '"')
	b.setval()
}

func (b *JSONBuffer) Seconds(d time.Duration) {
	d /= time.Second
	b.Int64(int64(d))
}

//-----------------------------------------------------------------------------
// Output
//-----------------------------------------------------------------------------

func (b *JSONBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.End())
	return int64(n), err
}

func (b *JSONBuffer) Print() {
	b.NewLine()
	_, err := b.WriteTo(os.Stdout)
	if err != nil {
		println("JSONBuffer write to stdout failed:", err.Error())
	}
}

//-----------------------------------------------------------------------------

type JSONUnixTime struct {
	time.Time
}

func (t JSONUnixTime) WriteJSON(b *JSONBuffer) error {
	b.Int64(t.Unix())
	return nil
}

//-----------------------------------------------------------------------------
// Helpers
//-----------------------------------------------------------------------------

func whiteSpace(n int) string {
	const space = "                                                            "
	const splen = len(space)
	if n <= splen {
		return space[:n]
	}
	buf := make([]byte, 0, n)
	for n >= splen {
		buf = append(buf, space...)
		n -= splen
	}
	if n > 0 {
		buf = append(buf, space[:n]...)
	}
	return *(*string)(unsafe.Pointer(&buf))
}

//-----------------------------------------------------------------------------
// String Escape
//-----------------------------------------------------------------------------

func escape(b []byte, s string, safeHTML bool) []byte {
	const hex = "0123456789abcdef"
	var esc = escASCII
	if safeHTML {
		esc = escHTML
	}
	start := 0
	for i := 0; i < len(s); {
		if c := s[i]; c < utf8.RuneSelf {
			e := esc[c]
			if e == 0 {
				i++
				continue
			}
			if start < i {
				b = append(b, s[start:i]...)
			}
			if e == 1 {
				b = append(b, '\\', 'u', '0', '0', hex[c>>4], hex[c&0xF])
			} else {
				b = append(b, '\\', e)
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				b = append(b, s[start:i]...)
			}
			b = append(b, '\\', 'u', 'f', 'f', 'f', 'd')
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				b = append(b, s[start:i]...)
			}
			b = append(b, '\\', 'u', '2', '0', '2', hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		b = append(b, s[start:]...)
	}
	return b
}

//-----------------------------------------------------------------------------
// bytes.Buffer Pool
//-----------------------------------------------------------------------------

var _byteBufferPool = sync.Pool{New: func() interface{} {
	return new(bytes.Buffer)
}}

func AcquireByteBuffer(size int) *bytes.Buffer {
	buf, _ := _byteBufferPool.Get().(*bytes.Buffer)
	if size > 0 && size > buf.Cap() {
		buf.Grow(size)
	}
	return buf
}

func ReleaseByteBuffer(buf *bytes.Buffer) {
	buf.Reset()
	_byteBufferPool.Put(buf)
}

//-----------------------------------------------------------------------------
// JSON Field Helpers
//-----------------------------------------------------------------------------

func CheckJSONField(field string, errorPrefix string) error {
	if len(errorPrefix) == 0 {
		errorPrefix = "JSON field"
	}
	n := len(field)
	if n == 0 {
		return fmt.Errorf("%s is empty", errorPrefix)
	}
	for i := 0; i < n; i++ {
		if c := field[i]; ascii[c] == 0 {
			r, _ := utf8.DecodeRuneInString(field[i:])
			return fmt.Errorf("%s has invalid character %q at index %d", errorPrefix, r, i+1)
		}
	}
	return nil
}

//-----------------------------------------------------------------------------

func MustJSONField(s string, errorPrefix string) {
	if err := CheckJSONField(s, errorPrefix); err != nil {
		panic(err)
	}
}

//-----------------------------------------------------------------------------
// Tables
//-----------------------------------------------------------------------------

var ascii = [256]uint8{
	'0': 1, '1': 1, '2': 1, '3': 1, '4': 1, '5': 1, '6': 1, '7': 1, '8': 1, '9': 1,
	'A': 1, 'B': 1, 'C': 1, 'D': 1, 'E': 1, 'F': 1, 'G': 1, 'H': 1, 'I': 1, 'J': 1,
	'K': 1, 'L': 1, 'M': 1, 'N': 1, 'O': 1, 'P': 1, 'Q': 1, 'R': 1, 'S': 1, 'T': 1,
	'U': 1, 'V': 1, 'W': 1, 'X': 1, 'Y': 1, 'Z': 1,
	'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1, 'f': 1, 'g': 1, 'h': 1, 'i': 1, 'j': 1,
	'k': 1, 'l': 1, 'm': 1, 'n': 1, 'o': 1, 'p': 1, 'q': 1, 'r': 1, 's': 1, 't': 1,
	'u': 1, 'v': 1, 'w': 1, 'x': 1, 'y': 1, 'z': 1,
	'-': 1, '_': 1, '@': 1, '*': 1,
}

var escASCII = [128]uint8{
	0: 1, 1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1,
	7:  1,   // \a
	8:  1,   // \b
	9:  't', // \t
	10: 'n', // \n
	11: 1,   // \v
	12: 1,   // \f
	13: 'r', // \r
	14: 1, 15: 1, 16: 1, 17: 1, 18: 1, 19: 1,
	20: 1, 21: 1, 22: 1, 23: 1, 24: 1, 25: 1,
	26: 1, 27: 1, 28: 1, 29: 1, 30: 1, 31: 1,
	'"':  '"',
	'\\': '\\',
}

var escHTML = [128]uint8{
	0: 1, 1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1,
	7:  1,   // \a
	8:  1,   // \b
	9:  't', // \t
	10: 'n', // \n
	11: 1,   // \v
	12: 1,   // \f
	13: 'r', // \r
	14: 1, 15: 1, 16: 1, 17: 1, 18: 1, 19: 1,
	20: 1, 21: 1, 22: 1, 23: 1, 24: 1, 25: 1,
	26: 1, 27: 1, 28: 1, 29: 1, 30: 1, 31: 1,
	'"':  '"',
	'\\': '\\',
	'<':  1,
	'>':  1,
	'&':  1,
}
