package htp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Error struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type ValidationErrors struct {
	Errors  []Error     `json:"errors"`
	Example interface{} `json:"example"`
}

func (errs *ValidationErrors) Decode(b *bytes.Buffer) error {
	if err := json.NewDecoder(b).Decode(&errs); err != nil {
		return fmt.Errorf("decoding bytes error: %w", err)
	}

	return nil
}
