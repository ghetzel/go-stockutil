package maputil

import (
	"fmt"
	"testing"

	"github.com/ghetzel/go-stockutil/sliceutil"

	"github.com/ghetzel/testify/require"
)

func TestDiffuseTypedOneTierScalar(t *testing.T) {
	var errs []error

	var input = make(map[string]any)
	var output = make(map[string]any)

	input["str:id"] = "test"
	input["name"] = "default-string"
	input["bool:enabled"] = "true"
	input["float:float"] = "2.7"

	if output, errs = DiffuseMapTyped(input, ".", ":"); len(errs) > 0 {
		for _, err := range errs {
			t.Errorf("%s\n", err)
		}
	}

	if v, ok := output["id"]; !ok || v != "test" {
		t.Errorf("Incorrect value '%s' for key %s", v, "id")
	}

	if v, ok := output["name"]; !ok || v != "default-string" {
		t.Errorf("Incorrect value '%s' for key %s", v, "default-string")
	}

	if v, ok := output["enabled"]; !ok || v != true {
		t.Errorf("Incorrect value '%s' for key %s", v, "enabled")
	}

	if v, ok := output["float"]; !ok || v != 2.7 {
		t.Errorf("Incorrect value '%s' for key %s", v, "float")
	}
}

func TestDiffuseTypedOneTierComplex(t *testing.T) {
	var assert = require.New(t)

	var errs []error

	var input = make(map[string]any)
	var output = make(map[string]any)

	input["str:array"] = []string{"first", "third", "fifth"}
	input["array2"] = []string{"first", "third", "fifth"}
	input["int:numary.0"] = "9"
	input["int:numary.1"] = "7"
	input["int:numary.2"] = "3"
	input["int:things.one"] = "1"
	input["int:things.two"] = "2"
	input["int:things.three"] = "3"

	if output, errs = DiffuseMapTyped(input, ".", ":"); len(errs) > 0 {
		for _, err := range errs {
			assert.NoError(err)
		}
	}

	//  test string array
	assert.Contains(output, `array`)
	assert.Len(output[`array`], 3)

	for i, v := range output["array"].([]string) {
		assert.Equal(v, input["str:array"].([]string)[i])
	}

	assert.Contains(output, `array2`)
	assert.Len(output[`array2`], 3)

	for i, v := range output["array2"].([]string) {
		assert.Equal(v, input["array2"].([]string)[i])
	}

	//  test int array
	assert.Contains(output, `numary`)
	assert.Len(output[`numary`], 3)
	assert.ElementsMatch(output["numary"], []int64{9, 7, 3})

	//  test string-int map
	assert.Contains(output, `things`)

	for k, v := range output["things"].(map[string]any) {
		switch k {
		case `one`:
			if v.(int64) != 1 {
				t.Errorf("Expected things['one'] = 1, got %v", v)
			}
		case `two`:
			if v.(int64) != 2 {
				t.Errorf("Expected things['two'] = 2, got %v", v)
			}
		case `three`:
			if v.(int64) != 3 {
				t.Errorf("Expected things['three'] = 3, got %v", v)
			}
		}
	}
}

func TestDiffuseTypedMultiTierScalar(t *testing.T) {
	var assert = require.New(t)
	var errs []error

	var input = make(map[string]any)
	var output = make(map[string]any)

	input["int:items.0"] = 54
	input["int:items.1"] = 77
	input["int:items.2"] = 82

	output, errs = DiffuseMapTyped(input, ".", ":")
	assert.Len(errs, 0)

	assert.ElementsMatch(output["items"], []int64{54, 77, 82})
}

func TestDiffuseTypedMultiTierComplex(t *testing.T) {
	var assert = require.New(t)
	var errs []error

	var input = make(map[string]any)
	var output = make(map[string]any)

	input["str:items.0.name"] = "First"
	input["int:items.0.age"] = 54
	input["str:items.1.name"] = "Second"
	input["int:items.1.age"] = 77
	input["str:items.2.name"] = "Third"
	input["int:items.2.age"] = 82

	output, errs = DiffuseMapTyped(input, ".", ":")
	assert.Len(errs, 0)

	assert.Len(output["items"], 3)

	if i_items, ok := output["items"]; ok {
		var items = i_items.([]any)

		for item_id, obj := range items {
			for k, v := range M(obj).MapNative() {
				switch k {
				case `name`:
					assert.Equal(v, input[fmt.Sprintf("str:items.%d.%s", item_id, k)])
				case `age`:
					assert.EqualValues(v, input[fmt.Sprintf("int:items.%d.%s", item_id, k)])
				}
			}
		}
	} else {
		t.Errorf("Key 'items' is missing from output: %v", output)
	}
}

func TestDiffuseTypedMultiTierMixed(t *testing.T) {
	var assert = require.New(t)
	var errs []error

	var input = make(map[string]any)
	var output = make(map[string]any)

	input["items.0.tags"] = []string{"base", "other"}
	input["items.1.tags"] = []string{"thing", "still-other", "more-other"}
	input["items.2.tags"] = []string{"last"}

	output, errs = DiffuseMapTyped(input, ".", ":")
	assert.Len(errs, 0)

	if i_items, ok := output["items"]; ok {
		var items = i_items.([]any)

		if len(items) != 3 {
			t.Errorf("Key 'items' should be an array with 3 elements, got %v", i_items)
		}

		for item_id, obj := range items {
			for k, v := range M(obj).MapNative() {
				var vAry = sliceutil.Stringify(v)

				if inValue, ok := input[fmt.Sprintf("items.%d.%s", item_id, k)]; !ok {
					t.Errorf("Key %s Incorrect, expected %s, got %s", fmt.Sprintf("items.%d.%s", item_id, k), inValue, v)
				} else {
					var inValueAry = inValue.([]string)

					for i, vAryV := range vAry {

						if vAryV != inValueAry[i] {
							t.Errorf("Key %s[%d] Incorrect, expected %s, got %s", fmt.Sprintf("items.%d.%s", item_id, k), i, inValueAry[i], vAryV)
						}
					}
				}
			}
		}
	} else {
		t.Errorf("Key 'items' is missing from output: %v", output)
	}
}
