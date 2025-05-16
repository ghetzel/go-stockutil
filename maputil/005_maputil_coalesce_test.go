package maputil

import (
	"testing"
)

func TestCoalesceOneTierScalar(t *testing.T) {
	var err error

	var input = make(map[string]any)
	var output = make(map[string]any)

	input["id"] = "test"
	input["enabled"] = true
	input["float"] = 2.7

	if output, err = CoalesceMap(input, "."); err != nil {
		t.Errorf("%s\n", err)
	}

	if v, ok := output["id"]; !ok || v != "test" {
		t.Errorf("Incorrect value '%s' for key %s", v, "id")
	}

	if v, ok := output["enabled"]; !ok || v != true {
		t.Errorf("Incorrect value '%s' for key %s", v, "enabled")
	}

	if v, ok := output["float"]; !ok || v != 2.7 {
		t.Errorf("Incorrect value '%s' for key %s", v, "float")
	}
}

func TestCoalesceMultiTierScalar(t *testing.T) {
	var err error

	var input = make(map[string]any)
	var output = make(map[string]any)

	input["id"] = "top"
	input["nested"] = make(map[string]any)
	input["nested"].(map[string]any)["data"] = true
	input["nested"].(map[string]any)["value"] = 4.9
	input["nested"].(map[string]any)["awesome"] = "very yes"

	if output, err = CoalesceMap(input, "."); err != nil {
		t.Errorf("%s\n", err)
	}

	if v, ok := output["id"]; !ok || v != "top" {
		t.Errorf("Incorrect value '%s' for key %s", v, "id")
	}

	if v, ok := output["nested.data"]; !ok || v != true {
		t.Errorf("Incorrect value '%s' for key %s", v, "nested.data")
	}

	if v, ok := output["nested.value"]; !ok || v != 4.9 {
		t.Errorf("Incorrect value '%s' for key %s", v, "nested.value")
	}

	if v, ok := output["nested.awesome"]; !ok || v != "very yes" {
		t.Errorf("Incorrect value '%s' for key %s", v, "nested.awesome")
	}
}

func TestCoalesceTopLevelArray(t *testing.T) {
	var err error

	var input = make(map[string]any)
	var output = make(map[string]any)

	var numbers = make([]any, 0)
	numbers = append(numbers, 1)
	numbers = append(numbers, 2)
	numbers = append(numbers, 3)

	input["numbers"] = numbers

	if output, err = CoalesceMap(input, "."); err != nil {
		t.Errorf("%s\n", err)
	}

	if v, ok := output["numbers.0"]; !ok || v != 1 {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.0")
	}

	if v, ok := output["numbers.1"]; !ok || v != 2 {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.1")
	}

	if v, ok := output["numbers.2"]; !ok || v != 3 {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.2")
	}
}

func TestCoalesceArrayWithNestedMap(t *testing.T) {
	var err error

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

	if output, err = CoalesceMap(input, "."); err != nil {
		t.Errorf("%s\n", err)
	}

	if v, ok := output["numbers.0.name"]; !ok || v != "test" {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.0.name")
	}

	if v, ok := output["numbers.0.count"]; !ok || v != 2 {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.0.count")
	}

	if v, ok := output["numbers.1.name"]; !ok || v != "test2" {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.1.name")
	}

	if v, ok := output["numbers.1.count"]; !ok || v != 4 {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.1.count")
	}

	if v, ok := output["numbers.2.name"]; !ok || v != "test3" {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.2.name")
	}

	if v, ok := output["numbers.2.count"]; !ok || v != 8 {
		t.Errorf("Incorrect value '%s' for key %s", v, "numbers.2.count")
	}
}
