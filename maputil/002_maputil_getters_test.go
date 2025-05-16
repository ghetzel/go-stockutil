package maputil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

type mapTestSubstruct struct {
	Name  string
	Value any
}

type mapTestStruct struct {
	Name    string
	NestedP *mapTestSubstruct
	Nested  mapTestSubstruct
}

func TestGetNil(t *testing.T) {
	var assert = require.New(t)

	var input = make(map[string]any)
	var level1 = make(map[string]any)

	level1["nilvalue"] = nil

	input["test"] = level1

	assert.Nil(DeepGet(input, []string{"test", "nilvalue"}, "nope"))
	assert.Nil(DeepGet(input, []string{"test", "nilvalue"}, nil))
}

func TestDeepGetScalar(t *testing.T) {
	var assert = require.New(t)

	var input = make(map[string]any)

	input = DeepSet(input, []string{"deeply", "nested", "value"}, 1.4).(map[string]any)

	assert.NotNil(DeepGet(input, []string{"deeply", "nested", "value"}, nil))
	assert.Equal(true, DeepGet(input, []string{"deeply", "nested", "value2"}, true))
	assert.Equal(`fallback`, DeepGet(input, []string{"deeply", "nested", "value2"}, "fallback"))
}

func TestDeepGetArrayElement(t *testing.T) {
	var input = make(map[string]any)

	input = DeepSet(input, []string{"tags", "0"}, "base").(map[string]any)
	input = DeepSet(input, []string{"tags", "1"}, "other").(map[string]any)

	if v := DeepGet(input, []string{"tags", "0"}, nil); v != "base" {
		t.Errorf("%s\n", v)
	}

	if v := DeepGet(input, []string{"tags", "1"}, nil); v != "other" {
		t.Errorf("%s\n", v)
	}
}

func TestDeepGetMapKeyInArray(t *testing.T) {
	var assert = require.New(t)

	var input = make(map[string]any)

	input = DeepSet(input, []string{"devices", "0", "name"}, "lo").(map[string]any)
	input = DeepSet(input, []string{"devices", "1", "name"}, "eth0").(map[string]any)

	assert.Equal(`lo`, DeepGet(input, []string{"devices", "0", "name"}, nil))
	assert.Equal(`eth0`, DeepGet(input, []string{"devices", "1", "name"}, nil))
}

func TestDeepGetMapKeyInDeepArray(t *testing.T) {
	var input = make(map[string]any)

	input = DeepSet(input, []string{"devices", "0", "switch", "peers", "0"}, "0.0.0.0").(map[string]any)
	input = DeepSet(input, []string{"devices", "0", "switch", "peers", "1"}, "0.0.1.1").(map[string]any)
	input = DeepSet(input, []string{"devices", "1", "switch", "peers", "0"}, "1.1.0.0").(map[string]any)
	input = DeepSet(input, []string{"devices", "1", "switch", "peers", "1"}, "1.1.1.1").(map[string]any)

	if v := DeepGet(input, []string{"devices", "0", "switch", "peers", "0"}, nil); v != "0.0.0.0" {
		t.Errorf("%s\n", v)
	}

	if v := DeepGet(input, []string{"devices", "0", "switch", "peers", "1"}, nil); v != "0.0.1.1" {
		t.Errorf("%s\n", v)
	}

	if v := DeepGet(input, []string{"devices", "1", "switch", "peers", "0"}, nil); v != "1.1.0.0" {
		t.Errorf("%s\n", v)
	}

	if v := DeepGet(input, []string{"devices", "1", "switch", "peers", "1"}, nil); v != "1.1.1.1" {
		t.Errorf("%s\n", v)
	}
}

func TestDeepGetBool(t *testing.T) {
	var assert = require.New(t)
	var input any

	input = make(map[string]any)

	input = DeepSet(input, []string{"deeply", "nested", "value"}, true)
	input = DeepSet(input, []string{"deeply", "nested", "thing"}, "nope")

	assert.True(DeepGetBool(input, []string{"deeply", "nested", "value"}))
	assert.False(DeepGetBool(input, []string{"deeply", "nested", "other"}))
	assert.False(DeepGetBool(input, []string{"deeply", "nested", "nope"}))
}

func TestDeepGetMapInMap(t *testing.T) {
	var assert = require.New(t)

	var in = map[string]any{
		`ok`: true,
		`always`: map[string]any{
			`finishing`: map[string]any{
				`each_others`: `sentences`,
			},
		},
	}

	assert.Equal(`sentences`, DeepGet(in, []string{`always`, `finishing`, `each_others`}))
	assert.Nil(DeepGet(in, []string{`always`, `finishing`, `each_others`, `sandwiches`}))
}

func TestDeepStructs(t *testing.T) {
	var assert = require.New(t)

	var in = &mapTestStruct{
		Name: `toplevel`,
		NestedP: &mapTestSubstruct{
			Name:  `one-ptr`,
			Value: true,
		},
		Nested: mapTestSubstruct{
			Name:  `one-value`,
			Value: 3.14,
		},
	}

	assert.Equal(`toplevel`, DeepGet(in, []string{`Name`}))
	assert.Equal(`one-ptr`, DeepGet(in, []string{`NestedP`, `Name`}))
	assert.Equal(true, DeepGet(in, []string{`NestedP`, `Value`}))
	assert.Equal(`one-value`, DeepGet(in, []string{`Nested`, `Name`}))
	assert.Equal(float64(3.14), DeepGet(in, []string{`Nested`, `Value`}))
}

func TestDeepGetNestedArrayElements(t *testing.T) {
	var assert = require.New(t)

	var input = map[string]any{
		`interfaces`: []string{
			`lo0`, `en1`, `wlan0`,
		},
	}

	assert.EqualValues([]any{
		`lo0`, `en1`, `wlan0`,
	}, DeepGet(input, []string{`interfaces`, `*`}))
}

func TestDeepGetNestedArrayOfMaps(t *testing.T) {
	var assert = require.New(t)

	var input = map[string]any{
		`interfaces`: []map[string]any{
			{
				`name`: `lo0`,
				`type`: `loopback`,
			}, {
				`name`: `en1`,
				`type`: `ethernet`,
			}, {
				`name`:     `wlan0`,
				`type`:     `ethernet`,
				`wireless`: true,
			},
		},
	}

	assert.EqualValues([]any{
		`loopback`, `ethernet`, `ethernet`,
	}, DeepGet(input, []string{`interfaces`, `*`, `type`}))

	assert.EqualValues([]any{
		false, false, true,
	}, DeepGet(input, []string{`interfaces`, `*`, `wireless`}, false))
}
