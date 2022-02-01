package shortuid

import (
	"strings"

	"github.com/google/uuid"
)

var DefaultEncoder = &base57{newAlphabet(DefaultAlphabet)}

type Encoder interface {
	Encode(uuid.UUID) string
	Decode(string) (uuid.UUID, error)
}

func NewUID() string {
	return DefaultEncoder.Encode(uuid.New())
}

func NewWithNamespace(name string) string {
	var u uuid.UUID

	switch {
	case name == "":
		u = uuid.New()
	case strings.HasPrefix(name, "http"):
		u = uuid.NewSHA1(uuid.NameSpaceURL, []byte(name))
	default:
		u = uuid.NewSHA1(uuid.NameSpaceDNS, []byte(name))
	}

	return DefaultEncoder.Encode(u)
}
