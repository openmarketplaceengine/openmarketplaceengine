package shortuid

import (
	"testing"

	"github.com/google/uuid"
)

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUid()
	}
}

func BenchmarkEncoding(b *testing.B) {
	u := uuid.New()
	for i := 0; i < b.N; i++ {
		DefaultEncoder.Encode(u)
	}
}

func BenchmarkDecoding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DefaultEncoder.Decode("c3eeb3e6-e577-4de2-b5bb-08371196b453")
	}
}
