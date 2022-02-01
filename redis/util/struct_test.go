package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	StringValue1 string       `json:"stringValue1"`
	StringValue2 string       `json:"stringValue2"`
	IntValue1    int          `json:"intValue1"`
	IntValue2    int          `json:"intValue2"`
	FloatValue   float64      `json:"floatValue"`
	BoolValue    bool         `json:"boolValue"`
	Nested       nestedStruct `json:"nested"`
}

type nestedStruct struct {
	IntValue    int     `json:"nestedIntValue"`
	BoolValue   bool    `json:"nestedBoolValue"`
	StringValue string  `json:"nestedStringValue"`
	FloatValue  float64 `json:"nestedFloatValue"`
}

type testStructPointer struct {
	JustValue   string        `json:"justValue"`
	JustPointer *nestedStruct `json:"justPointer"`
}

func TestStructToMap(t *testing.T) {
	t.Run("testStructToMap", func(t *testing.T) {
		testStructToMap(t)
	})

	t.Run("testStructToMapNestedPointer", func(t *testing.T) {
		testStructToMapNestedPointer(t)
	})

	t.Run("testStructToMapPointer", func(t *testing.T) {
		testStructToMapPointer(t)
	})

	t.Run("testMapToStruct", func(t *testing.T) {
		testMapToStruct(t)
	})
}

func testStructToMap(t *testing.T) {
	example := testStruct{
		StringValue1: "Hello World",
		StringValue2: "",
		IntValue1:    123,
		IntValue2:    0,
		FloatValue:   1.2,
		BoolValue:    true,
		Nested: nestedStruct{
			IntValue:    123,
			BoolValue:   false,
			StringValue: "abc",
			FloatValue:  1.2,
		},
	}

	m, err := StructToMap(example)
	require.NoError(t, err)
	require.NotNil(t, m)
	assert.Equal(t, m["StringValue1"], example.StringValue1)
	assert.Equal(t, m["StringValue2"], example.StringValue2)
	assert.Equal(t, m["IntValue1"], fmt.Sprintf("%v", example.IntValue1))
	assert.Equal(t, m["IntValue2"], fmt.Sprintf("%v", example.IntValue2))
	assert.Equal(t, m["FloatValue"], fmt.Sprintf("%v", example.FloatValue))
	assert.Equal(t, m["BoolValue"], fmt.Sprintf("%v", example.BoolValue))
	assert.Equal(t, m["Nested"], "{\"nestedIntValue\":123,\"nestedBoolValue\":false,\"nestedStringValue\":\"abc\",\"nestedFloatValue\":1.2}")
	assert.Equal(t, m["IntValue1"], fmt.Sprintf("%v", example.IntValue1))
}

func testStructToMapNestedPointer(t *testing.T) {
	example := testStructPointer{
		JustValue: "Hello World",
		JustPointer: &nestedStruct{
			IntValue:    0,
			BoolValue:   false,
			StringValue: "",
			FloatValue:  0,
		},
	}

	m, err := StructToMap(example)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is a field pointer")
	require.Nil(t, m)
}

func testStructToMapPointer(t *testing.T) {
	m, err := StructToMap(&testStruct{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "is a field pointer")
	require.Nil(t, m)
}

func testMapToStruct(t *testing.T) {
	exampleIn := testStruct{
		StringValue1: "Hello World",
		StringValue2: "",
		IntValue1:    123,
		IntValue2:    0,
		FloatValue:   1.2,
		BoolValue:    true,
		Nested: nestedStruct{
			IntValue:    123,
			BoolValue:   false,
			StringValue: "abc",
			FloatValue:  1.2,
		},
	}

	m, err := StructToMap(exampleIn)
	require.NoError(t, err)
	require.NotNil(t, m)

	var exampleOut testStruct

	err = MapToStruct(TransformMap(m), &exampleOut)
	require.NoError(t, err)

	assert.Equal(t, exampleIn, exampleOut)
}
