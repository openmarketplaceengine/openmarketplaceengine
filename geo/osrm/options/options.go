package options

import (
	"fmt"
	"net/url"
)

func UrlEncode(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	var bytes []byte
	for k, v := range params {
		if len(bytes) > 0 {
			bytes = append(bytes, '&')
		}
		bytes = append(bytes, url.QueryEscape(k)...)
		bytes = append(bytes, '=')
		intSlice, ok := v.([]int)
		if ok {
			bytes = append(bytes, url.QueryEscape(concatInt(intSlice))...)
		} else {
			bytes = append(bytes, url.QueryEscape(fmt.Sprintf("%v", v))...)
		}
	}
	return string(bytes)
}

func concatInt(val []int) string {
	var buf []byte
	for i, v := range val {
		if i > 0 {
			buf = append(buf, ';')
		}
		buf = append(buf, fmt.Sprintf("%v", v)...)
	}
	return string(buf)
}
