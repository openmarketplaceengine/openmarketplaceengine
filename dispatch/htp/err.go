package htp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/validate"
)

type ValidationErrors struct {
	Errors  []validate.Error `json:"errors"`
	Example interface{}      `json:"example"`
}

func (errs *ValidationErrors) Decode(b *bytes.Buffer) error {
	if err := json.NewDecoder(b).Decode(&errs); err != nil {
		return fmt.Errorf("decoding bytes error: %w", err)
	}

	return nil
}
