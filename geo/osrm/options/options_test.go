package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlank(t *testing.T) {
	var params map[string]interface{}
	assert.Equal(t, "", UrlEncode(params))
}

func TestOne(t *testing.T) {
	params := map[string]interface{}{"k1": "v1"}
	assert.Equal(t, "k1=v1", UrlEncode(params))
}

func TestBool(t *testing.T) {
	params := map[string]interface{}{"k1": true}
	assert.Equal(t, "k1=true", UrlEncode(params))
}

func TestTwo(t *testing.T) {
	params := map[string]interface{}{"k1": "v1", "k2": "v2"}
	encoded := UrlEncode(params)
	assert.Contains(t, encoded, "k1=v1")
	assert.Contains(t, encoded, "&")
	assert.Contains(t, encoded, "k2=v2")
}

func TestArray(t *testing.T) {
	params := map[string]interface{}{"k1": []int{0, 1, 2}}
	assert.Equal(t, "k1=0%3B1%3B2", UrlEncode(params))
}

func TestArrayPlus(t *testing.T) {
	params := map[string]interface{}{"k1": []int{0, 1, 2}, "k2": "v2"}
	encoded := UrlEncode(params)
	assert.Contains(t, encoded, "k1=0%3B1%3B2")
	assert.Contains(t, encoded, "&")
	assert.Contains(t, encoded, "k2=v2")
}
