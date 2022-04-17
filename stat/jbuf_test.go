package stat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//-----------------------------------------------------------------------------

func TestJSONBuffer_Rune(t *testing.T) {
	s := "abcdœΩåßç√©∆˚🙂🐸🥝🇦🇶✅"
	var b JSONBuffer
	for _, r := range s {
		b.Rune(r)
	}
	require.Equal(t, s, b.UnsafeString())
}

//-----------------------------------------------------------------------------

func TestJSONBuffer_StrEsc(t *testing.T) {
	smap := map[string]string{
		"":       `""`,
		`""`:     `"\"\""`,
		"abcd":   `"abcd"`,
		"\n\t\r": `"\n\t\r"`,
		`<&>`:    `"\u003c\u0026\u003e"`,
		`"<&>"`:  `"\"\u003c\u0026\u003e\""`,
		"\\":     `"\\"`,
		"abcdœΩåßç√©∆˚🙂🐸🥝🇦🇶✅": `"abcdœΩåßç√©∆˚🙂🐸🥝🇦🇶✅"`,
	}
	var b JSONBuffer
	for k, v := range smap {
		b.Reset()
		b.String(k, true)
		require.Equal(t, v, b.UnsafeString())
	}
}

//-----------------------------------------------------------------------------

func TestJSONBuffer_EmptyArrayIndent(t *testing.T) {
	b := AcquireJSONBuffer(2)
	t.Cleanup(b.Release)
	b.ArrayStart()
	b.ArrayStart()
	b.ArrayStart()
	b.ArrayClose()
	b.ArrayClose()
	b.ArrayClose()
	want := "[\n  [\n    [\n    ]\n  ]\n]\n"
	require.Equal(t, want, string(b.End()))
}

func TestJSONBuffer_ArrayInline(t *testing.T) {
	var b JSONBuffer
	fillArray(&b, 2, true)
	b.Print()
}

func fillArray(b *JSONBuffer, indent int, inline bool) {
	now := time.Now()
	b.Reset()
	b.WithIndent(indent)
	b.Inline(inline)
	b.ArrayStart()
	b.Null()
	b.Comma()
	b.Int(1)
	b.Comma()
	b.True()
	b.Comma()
	b.False()
	b.Comma()
	_ = b.Value(now)
	b.Comma()
	b.UnixTime(now)
	b.Comma()
	_ = b.Value(nil)
	b.Comma()
	b.Comma()
	b.Comma()
	b.ArrayClose()
}
