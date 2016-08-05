package maputil

import (
	"strings"
	"testing"
	"time"
)

type MyTestThing struct {
	Name  string
	Other int
}

func TestMapJoin(t *testing.T) {
	input := map[string]interface{}{
		`key1`: `value1`,
		`key2`: true,
		`key3`: 3,
	}

	output := Join(input, `=`, `&`)

	if output == `` {
		t.Error("Output should not be empty")
	}

	if !strings.Contains(output, `key1=value1`) {
		t.Errorf("Output should contain '%s'", `key1=value1`)
	}

	if !strings.Contains(output, `key2=true`) {
		t.Errorf("Output should contain '%s'", `key2=true`)
	}

	if !strings.Contains(output, `key3=3`) {
		t.Errorf("Output should contain '%s'", `key3=3`)
	}
}

func TestStringKeys(t *testing.T) {
	i1 := map[string]interface{}{
		`1`: 1,
		`2`: true,
		`3`: `three`,
	}

	i2 := map[string]bool{
		`1`: true,
		`2`: false,
		`3`: true,
	}

	i3 := map[string]MyTestThing{
		`1`: MyTestThing{},
		`2`: MyTestThing{},
		`3`: MyTestThing{},
	}

	output := []string{`1`, `2`, `3`}

	for i, v := range StringKeys(i1) {
		if key := output[i]; key != v {
			t.Errorf("map[string]interface{} key %d; expected: %q, got: %q", i, key, v)
		}
	}

	for i, v := range StringKeys(i2) {
		if key := output[i]; key != v {
			t.Errorf("map[string]bool key %d; expected: %q, got: %q", i, key, v)
		}
	}

	for i, v := range StringKeys(i3) {
		if key := output[i]; key != v {
			t.Errorf("map[string]MyTestThing key %d; expected: %q, got: %q", i, key, v)
		}
	}

	if keys := StringKeys(true); len(keys) != 0 {
		t.Errorf("StringKeys(true); expected: len=0, got: len=%d", len(keys))
	}

	if keys := StringKeys(4); len(keys) != 0 {
		t.Errorf("StringKeys(4); expected: len=0, got: len=%d", len(keys))
	}

	if keys := StringKeys([]int{1, 2, 3}); len(keys) != 0 {
		t.Errorf("StringKeys([]int{1,2,3}); expected: len=0, got: len=%d", len(keys))
	}
}

func TestMapSplit(t *testing.T) {
	input := `key1=value1&key2=true&key3=3`

	output := Split(input, `=`, `&`)

	if len(output) == 0 {
		t.Error("Output should not be empty")
	}

	if v, ok := output[`key1`]; !ok || v != `value1` {
		t.Errorf("Output should contain key %s => '%s'", `key1`, `value1`)
	}

	if v, ok := output[`key2`]; !ok || v != `true` {
		t.Errorf("Output should contain key %s => '%s'", `key2`, `true`)
	}

	if v, ok := output[`key3`]; !ok || v != `3` {
		t.Errorf("Output should contain key %s => '%s'", `key3`, `3`)
	}
}

type SubtypeTester struct {
	A int
	B int `maputil:"b"`
}

type MyStructTester struct {
	Name                  string
	Subtype1              SubtypeTester
	Active                bool           `maputil:"active"`
	Subtype2              *SubtypeTester `maputil:"subtype2"`
	TimeTest              time.Duration
	IntTest               int32
	Properties            map[string]interface{}
	StrSliceTest          []string
	InterfaceStrSliceTest []string
	StructSliceTest       []SubtypeTester
	StructSliceTest2      []SubtypeTester
	nonexported           int
}

func TestStructFromMap(t *testing.T) {
	input := map[string]interface{}{
		`Name`:           `Foo Bar`,
		`active`:         true,
		`should_not_set`: 4,
		`Subtype1`: map[string]interface{}{
			`A`: 1,
			`b`: 2,
		},
		`subtype2`: map[string]interface{}{
			`A`: 3,
			`b`: 4,
		},
		`TimeTest`: 15000000000,
		`IntTest`:  int64(5),
		`Properties`: map[string]interface{}{
			`first`:  1,
			`second`: true,
			`third`:  `three`,
		},
		`StrSliceTest`:          []string{`one`, `two`, `three`},
		`InterfaceStrSliceTest`: []interface{}{`one`, `two`, `three`},
		`StructSliceTest`:       []SubtypeTester{{10, 11}, {12, 13}, {14, 15}},
		`StructSliceTest2`: []map[string]interface{}{
			{
				`A`: 10,
				`b`: 11,
			},
			{
				`A`: 12,
				`b`: 13,
			},
			{
				`A`: 14,
				`b`: 15,
			},
		},
	}

	output := MyStructTester{}

	if err := StructFromMap(input, &output); err == nil {
		if output.Name != `Foo Bar` {
			t.Errorf("output.Name; expected: %s, got: %v", `Foo Bar`, output.Name)
		}

		if !output.Active {
			t.Errorf("output.Active; expected: true, got: false")
		}

		if output.nonexported != 0 {
			t.Errorf("output.nonexported; expected: 0, got: %v", output.nonexported)
		}

		if output.Subtype1.A != 1 {
			t.Errorf("output.Subtype1.A; expected: 1, got: %v", output.Subtype1.A)
		}

		if output.Subtype1.B != 2 {
			t.Errorf("output.Subtype1.B; expected: 2, got: %v", output.Subtype1.B)
		}

		if output.Subtype2 == nil {
			t.Errorf("output.Subtype2; is nil, should be populated with an instance")
		} else {
			if output.Subtype2.A != 3 {
				t.Errorf("output.Subtype2.A; expected: 3, got: %v", output.Subtype2.A)
			}

			if output.Subtype2.B != 4 {
				t.Errorf("output.Subtype2.B; expected: 4, got: %v", output.Subtype2.B)
			}
		}

		if output.TimeTest != time.Duration(15)*time.Second {
			t.Errorf("output.TimeTest; expected: 15s, got: %s", output.TimeTest)
		}

		if output.IntTest != int32(5) {
			t.Errorf("output.IntTest; expected: 5(int32), got: %d(%T)", output.IntTest, output.IntTest)
		}

		if output.Properties == nil {
			t.Errorf("output.Properties; is nil, should be map[string]interface{}")
		} else {
			if v, ok := output.Properties[`first`]; !ok || v != 1 {
				t.Errorf("output.Properties['first']; expected: 1, got: %v(%T)", v, v)
			}

			if v, ok := output.Properties[`second`]; !ok || v != true {
				t.Errorf("output.Properties['second']; expected: true, got: %v(%T)", v, v)
			}

			if v, ok := output.Properties[`third`]; !ok || v != `three` {
				t.Errorf("output.Properties['third']; expected: 'three', got: %v(%T)", v, v)
			}
		}

		if output.StrSliceTest == nil {
			t.Errorf("output.StrSliceTest; is nil, should be []string")
		} else {
			if l := len(output.StrSliceTest); l == 3 {
				if v := output.StrSliceTest[0]; v != `one` {
					t.Errorf("output.StrSliceTest[0]; expected: 'one', got: %s", v)
				}

				if v := output.StrSliceTest[1]; v != `two` {
					t.Errorf("output.StrSliceTest[1]; expected: 'two', got: %s", v)
				}

				if v := output.StrSliceTest[2]; v != `three` {
					t.Errorf("output.StrSliceTest[2]; expected: 'three', got: %s", v)
				}
			} else {
				t.Errorf("output.StrSliceTest; wrong length - expected: 3, got: %d", l)
			}
		}

		if output.InterfaceStrSliceTest == nil {
			t.Errorf("output.InterfaceStrSliceTest; is nil, should be []interface{}")
		} else {
			if l := len(output.InterfaceStrSliceTest); l == 3 {
				if v := output.InterfaceStrSliceTest[0]; v != `one` {
					t.Errorf("output.InterfaceStrSliceTest[0]; expected: 'one', got: %s", v)
				}

				if v := output.InterfaceStrSliceTest[1]; v != `two` {
					t.Errorf("output.InterfaceStrSliceTest[1]; expected: 'two', got: %s", v)
				}

				if v := output.InterfaceStrSliceTest[2]; v != `three` {
					t.Errorf("output.InterfaceStrSliceTest[2]; expected: 'three', got: %s", v)
				}
			} else {
				t.Errorf("output.InterfaceStrSliceTest; wrong length - expected: 3, got: %d", l)
			}
		}

		if output.StructSliceTest == nil {
			t.Errorf("output.StructSliceTest; is nil, should be []SubtypeTester")
		} else {
			if l := len(output.StructSliceTest); l == 3 {
				if v := output.StructSliceTest[0]; v.A != 10 || v.B != 11 {
					t.Errorf("output.StructSliceTest[0]; expected: {10,11}, got: %v", v)
				}

				if v := output.StructSliceTest[1]; v.A != 12 || v.B != 13 {
					t.Errorf("output.StructSliceTest[1]; expected: {12,13}, got: %v", v)
				}

				if v := output.StructSliceTest[2]; v.A != 14 || v.B != 15 {
					t.Errorf("output.StructSliceTest[2]; expected: {14,15}, got: %v", v)
				}
			} else {
				t.Errorf("output.StructSliceTest; wrong length - expected: 3, got: %d", l)
			}
		}

		if output.StructSliceTest2 == nil {
			t.Errorf("output.StructSliceTest2; is nil, should be []SubtypeTester")
		} else {
			if l := len(output.StructSliceTest2); l == 3 {
				if v := output.StructSliceTest2[0]; v.A != 10 || v.B != 11 {
					t.Errorf("output.StructSliceTest2[0]; expected: {10,11}, got: %v", v)
				}

				if v := output.StructSliceTest2[1]; v.A != 12 || v.B != 13 {
					t.Errorf("output.StructSliceTest2[1]; expected: {12,13}, got: %v", v)
				}

				if v := output.StructSliceTest2[2]; v.A != 14 || v.B != 15 {
					t.Errorf("output.StructSliceTest2[2]; expected: {14,15}, got: %v", v)
				}
			} else {
				t.Errorf("output.StructSliceTest2; wrong length - expected: 3, got: %d", l)
			}
		}
	} else {
		t.Error(err)
	}
}
