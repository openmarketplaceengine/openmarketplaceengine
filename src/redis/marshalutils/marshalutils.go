package marshalutils

import (
	"encoding/json"
	"fmt"
)

func MapToStruct(m map[string]string, dest interface{}) error {

	bytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("error marshaling map: %w", err)
	}

	err = json.Unmarshal(bytes, dest)
	if err != nil {
		return fmt.Errorf("error unmarshaling map bytes to struct: %w", err)
	}

	return nil
}

func StructToMap(src interface{}) (map[string]interface{}, error) {

	bytes, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("error marshaling struct: %w", err)
	}

	var dest map[string]interface{}
	err = json.Unmarshal(bytes, &dest)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling struct bytes to map: %w", err)
	}

	return dest, nil
}
