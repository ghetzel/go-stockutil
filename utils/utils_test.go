package utils

import (
	"testing"
	"time"

	"github.com/ghetzel/testify/require"
)

type testEnum int

const (
	Value1 testEnum = iota
	Value2
	Value3
)

func (self testEnum) String() string {
	switch self {
	case Value1:
		return `value-1`
	case Value2:
		return `value-2`
	case Value3:
		return `value-3`
	default:
		return ``
	}
}

type testMarshal struct {
	Name      string `json:"name"`
	Count     int    `json:"count,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time `json:",omitempty"`
	Thing     testEnum  `json:"enum"`
}

func TestGenericMarshalJSON(t *testing.T) {
	var assert = require.New(t)

	data, err := GenericMarshalJSON(&testMarshal{
		Name:      `test`,
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Thing:     Value2,
	})

	assert.NoError(err)
	assert.Equal(
		`{"CreatedAt":"2009-11-10T23:00:00Z","enum":"value-2","name":"test"}`,
		string(data),
	)
}
