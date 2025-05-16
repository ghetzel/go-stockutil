package maputil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestCoalesceTypedOneTierScalar(t *testing.T) {
	var assert = require.New(t)

	var errs []error

	var input = make(map[string]any)
	var output = make(map[string]any)

	input["id"] = "test"
	input["enabled"] = true
	input["float"] = 2.7

	if output, errs = CoalesceMapTyped(input, ".", ":"); len(errs) > 0 {
		for _, err := range errs {
			t.Errorf("%s\n", err)
		}
	}

	assert.Equal(`test`, output[`str:id`])
	assert.Equal(`true`, output[`bool:enabled`])
	assert.Equal(`2.7`, output[`float:float`])
}

func TestCoalesceTypedMultiTierScalar(t *testing.T) {
	var assert = require.New(t)

	var errs []error

	var input = make(map[string]any)
	var output = make(map[string]any)

	input["id"] = "top"
	input["nested"] = map[string]any{
		`data`:    true,
		`value`:   4.9,
		`awesome`: "very yes",
	}

	if output, errs = CoalesceMapTyped(input, "__", "|"); len(errs) > 0 {
		for _, err := range errs {
			t.Errorf("%s\n", err)
		}
	}

	assert.Equal(`top`, output[`str|id`])
	assert.Equal(`true`, output[`bool|nested__data`])
	assert.Equal(`4.9`, output[`float|nested__value`])
	assert.Equal(`very yes`, output[`str|nested__awesome`])
}

func TestCoalesceTypedTopLevelArray(t *testing.T) {
	var errs []error

	var input = make(map[string]any)
	var output = make(map[string]any)

	var numbers = make([]any, 0)
	numbers = append(numbers, 1)
	numbers = append(numbers, 2)
	numbers = append(numbers, 3)

	input["numbers"] = numbers

	if output, errs = CoalesceMapTyped(input, ".", ":"); len(errs) > 0 {
		for _, err := range errs {
			t.Errorf("%s\n", err)
		}
	}

	if v, ok := output["int:numbers.0"]; !ok || v != "1" {
		t.Errorf("Incorrect value '%s' for key %s", v, "int:numbers.0")
	}

	if v, ok := output["int:numbers.1"]; !ok || v != "2" {
		t.Errorf("Incorrect value '%s' for key %s", v, "int:numbers.1")
	}

	if v, ok := output["int:numbers.2"]; !ok || v != "3" {
		t.Errorf("Incorrect value '%s' for key %s", v, "int:numbers.2")
	}
}

func TestCoalesceTypedArrayWithNestedMap(t *testing.T) {
	var errs []error

	var input = make(map[string]any)
	var output = make(map[string]any)

	var numbers = make([]any, 0)
	numbers = append(numbers, map[string]any{
		"name":  "test",
		"count": 2,
	})

	numbers = append(numbers, map[string]any{
		"name":  "test2",
		"count": 4,
	})

	numbers = append(numbers, map[string]any{
		"name":  "test3",
		"count": 8,
	})

	input["numbers"] = numbers

	if output, errs = CoalesceMapTyped(input, ".", ":"); len(errs) > 0 {
		for _, err := range errs {
			t.Errorf("%s\n", err)
		}
	}

	if v, ok := output["str:numbers.0.name"]; !ok || v != "test" {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.0.name")
	}

	if v, ok := output["int:numbers.0.count"]; !ok || v != "2" {
		t.Errorf("Incorrect value '%s' for key %s", v, "int:numbers.0.count")
	}

	if v, ok := output["str:numbers.1.name"]; !ok || v != "test2" {
		t.Errorf("Incorrect value '%s' for key %s", v, "str:numbers.1.name")
	}

	if v, ok := output["int:numbers.1.count"]; !ok || v != "4" {
		t.Errorf("Incorrect value '%s' for key %s", v, "int:numbers.1.count")
	}

	if v, ok := output["str:numbers.2.name"]; !ok || v != "test3" {
		t.Errorf("Incorrect value '%s' for key %s", v, "str:numbers.2.name")
	}

	if v, ok := output["int:numbers.2.count"]; !ok || v != "8" {
		t.Errorf("Incorrect value '%s' for key %s", v, "int:numbers.2.count")
	}
}
