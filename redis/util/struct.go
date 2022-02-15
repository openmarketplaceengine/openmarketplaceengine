package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// StructToMap purpose is to transform struct to store in redis HSET.
// This makes efficient updating one struct field, while not touching the others.
// If we store struct as json string - updating one field would require:
// GET(redis) -> json.unmarshal -> update.field -> json.marshal -> SET()redis
//
//Limitations: struct fields must be values.
//Nested structs will be json marshalled and expected to be not mutable.
func StructToMap(item interface{}) (m map[string]interface{}, err error) {
	if item == nil {
		err = errors.New("argument is nil")
		return
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		err = fmt.Errorf("values expected, %q is a field pointer", v.Name())
		return
	}
	m = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Field(i).Name
		field := reflectValue.Field(i).Interface()

		if reflect.TypeOf(field).Kind() == reflect.Ptr {
			m = nil
			err = fmt.Errorf("values expected, %q is a field pointer", v.Name())
			return
		}

		if v.Field(i).Type.Kind() == reflect.Struct {
			bytes, innerErr := json.Marshal(field)
			if innerErr != nil {
				err = errors.Wrap(innerErr, fmt.Sprintf("unable to json marshal field %v", field))
				return
			}
			m[fieldName] = string(bytes)
		} else {
			m[fieldName] = fmt.Sprintf("%v", field)
		}
	}
	return m, nil
}

func MapToStruct(m map[string]string, item interface{}) error {
	itemType := reflect.TypeOf(item)
	if itemType.Elem().Kind() == reflect.Ptr {
		return fmt.Errorf("pointer expected, got item=%v", item)
	}

	itemElem := reflect.ValueOf(item).Elem()

	for k, v := range m {
		field := itemElem.FieldByName(k)
		fieldKind := field.Kind()
		value := reflect.ValueOf(v)

		switch fieldKind {
		case reflect.Bool:
			b, innerErr := strconv.ParseBool(value.String())
			if innerErr != nil {
				return errors.Wrap(innerErr, fmt.Sprintf("unable to ParseBool value=%v", value.String()))
			}
			field.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, innerErr := strconv.ParseInt(value.String(), 10, 64)
			if innerErr != nil {
				return errors.Wrap(innerErr, fmt.Sprintf("unable to ParseInt value=%v", value.String()))
			}
			field.SetInt(i)
		case reflect.Float32, reflect.Float64:
			f, innerErr := strconv.ParseFloat(value.String(), 64)
			if innerErr != nil {
				return errors.Wrap(innerErr, fmt.Sprintf("unable to ParseFloa value=%v", value.String()))
			}
			field.SetFloat(f)
		case reflect.String:
			field.SetString(reflect.ValueOf(v).String())
		case reflect.Struct:
			t := reflect.TypeOf(field.Interface())
			nested := reflect.New(t).Interface()
			innerErr := json.Unmarshal([]byte(v), nested)
			if innerErr != nil {
				return errors.Wrap(innerErr, fmt.Sprintf("unable unmarshal nested struct value=%v", v))
			}
			field.Set(reflect.ValueOf(nested).Elem())
		case reflect.Slice:
			s := value.String()
			s = strings.ReplaceAll(s, "[", "")
			s = strings.ReplaceAll(s, "]", "")
			split := strings.Split(s, " ")
			field.Set(reflect.ValueOf(split))

		default:
			return fmt.Errorf("unsupported type=%v", fieldKind)
		}
	}
	return nil
}

func TransformMap(in map[string]interface{}) (out map[string]string) {
	out = make(map[string]string)
	for k, v := range in {
		out[k] = fmt.Sprintf("%v", v)
	}
	return
}
