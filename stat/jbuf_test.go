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

//-----------------------------------------------------------------------------
// Trim & White Space
//-----------------------------------------------------------------------------

func TestTrimRight(t *testing.T) {
	data := map[string]string{
		"":            "",
		" ":           "",
		" \t\n\r ":    "",
		"xyz \t\n\r ": "xyz",
		"abc":         "abc",
	}
	for k, v := range data {
		have := string(trimRight([]byte(k)))
		require.Equal(t, v, have)
	}
}

func TestJSONBuffer_AppendJSON_Trim(t *testing.T) {
	var b JSONBuffer
	b.ArrayStart()
	_ = b.AppendJSON([]byte("12345\n\t\r   "))
	b.ArrayClose()
	require.Equal(t, "[12345]", b.UnsafeString())
}

const tspaces = " \t\r\n abcde \n\r\t "

func checkSpace(c byte) bool {
	switch c {
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	}
}

func BenchmarkSpaceSwitch(b *testing.B) {
	b.ReportAllocs()
	counter := 0
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(tspaces); j++ {
			if checkSpace(tspaces[j]) {
				counter++
			}
		}
	}
	if counter == 0 {
		b.Fatal("counter is zero")
	}
}

func BenchmarkSpaceTable(b *testing.B) {
	b.ReportAllocs()
	counter := 0
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(tspaces); j++ {
			if space[tspaces[j]] {
				counter++
			}
		}
	}
	if counter == 0 {
		b.Fatal("counter is zero")
	}
}

func BenchmarkSpaceDetect(b *testing.B) {
	b.Run("Switch", BenchmarkSpaceSwitch)
	b.Run("Table", BenchmarkSpaceTable)
}
