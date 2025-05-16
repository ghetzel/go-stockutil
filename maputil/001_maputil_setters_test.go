package maputil

import (
	_ "encoding/json"
	"fmt"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestDeepSetNothing(t *testing.T) {
	var assert = require.New(t)

	var output = make(map[string]any)
	output = DeepSet(output, []string{}, "yay").(map[string]any)

	assert.Empty(output)
}

func TestDeepSetReplace(t *testing.T) {
	var assert = require.New(t)

	var output = map[string]any{
		`this`: map[string]any{
			`test`: `1`,
		},
	}

	DeepSet(output, []string{"this", "test"}, `2`)
	assert.Equal(map[string]any{
		`this`: map[string]any{
			`test`: `2`,
		},
	}, output)
}

func TestDeepSetString(t *testing.T) {
	var assert = require.New(t)

	var output = make(map[string]any)
	var testValue = "test-string"

	output = DeepSet(output, []string{"str"}, testValue).(map[string]any)

	value, ok := output["str"]
	assert.True(ok)
	assert.Equal(testValue, value)
}

func TestDeepSetBool(t *testing.T) {
	var assert = require.New(t)

	var output = make(map[string]any)
	var testValue = true

	output = DeepSet(output, []string{"bool"}, testValue).(map[string]any)

	value, ok := output["bool"]
	assert.True(ok)
	assert.Equal(testValue, value)
}

func TestDeepSetArray(t *testing.T) {
	var assert = require.New(t)

	var output = make(map[string]any)
	var testValues = []string{"first", "second"}

	for i, tv := range testValues {
		output = DeepSet(output, []string{"top-array", fmt.Sprint(i)}, tv).(map[string]any)
	}

	// output = DeepSet(output, []string{"top-array"}, 3.4).(map[string]any)

	// fmt.Println(typeutil.Dump(output))
	topArray, ok := output["top-array"]
	assert.True(ok)

	assert.ElementsMatch(testValues, topArray)
}

func TestDeepSetArrayIndices(t *testing.T) {
	var assert = require.New(t)

	var input = map[string]any{
		`things`: map[string]any{
			`type1`: []string{
				`first`,
				`second`,
				`third`,
			},
			`type2`: []string{
				`first`,
				`second`,
				`third`,
			},
			`type3`: []any{
				map[string]any{
					`name`:  `first`,
					`index`: 0,
				},
				map[string]any{
					`name`:  `first`,
					`index`: 1,
				},
				map[string]any{
					`name`:  `first`,
					`index`: 2,
				},
			},
		},
	}

	var output = DeepSet(input, []string{`things`, `type1`, `0`}, `First`)
	DeepSet(output, []string{`things`, `type1`, `2`}, `Third`)
	DeepSet(output, []string{`things`, `type2`, `1`}, `Second`)
	DeepSet(output, []string{`things`, `type2`, `2`}, nil)
	DeepSet(output, []string{`things`, `type2`, `3`}, `third`)
	DeepSet(output, []string{`things`, `type3`, `0`, `index`}, map[string]any{
		`num`: 0,
	})
	DeepSet(output, []string{`things`, `type3`, `1`, `index`}, map[string]any{
		`num`: 1,
	})
	DeepSet(output, []string{`things`, `type3`, `2`, `index`}, map[string]any{
		`num`: 2,
	})

	assert.Equal(map[string]any{
		`things`: map[string]any{
			`type1`: []any{
				`First`,
				`second`,
				`Third`,
			},
			`type2`: []any{
				`first`,
				`Second`,
				nil,
				`third`,
			},
			`type3`: []any{
				map[string]any{
					`name`: `first`,
					`index`: map[string]any{
						`num`: 0,
					},
				},
				map[string]any{
					`name`: `first`,
					`index`: map[string]any{
						`num`: 1,
					},
				},
				map[string]any{
					`name`: `first`,
					`index`: map[string]any{
						`num`: 2,
					},
				},
			},
		},
	}, output)
}

func TestDeepSetNestedMapCreation(t *testing.T) {
	var assert = require.New(t)

	var output = make(map[string]any)
	output = DeepSet(output, []string{"deeply", "nested", "map"}, true).(map[string]any)
	output = DeepSet(output, []string{"deeply", "nested", "count"}, 2).(map[string]any)

	deeply, ok := output["deeply"]
	assert.True(ok)

	var deeplyMap = deeply.(map[string]any)

	nested, ok := deeplyMap["nested"]
	assert.True(ok)

	var nestedMap = nested.(map[string]any)

	_, ok = nestedMap["map"]
	assert.True(ok)

	_, ok = nestedMap["count"]
	assert.True(ok)
}

func TestDeepSetStructField(t *testing.T) {
	var assert = require.New(t)
	type testStructDeepSet struct {
		String string
		Int    int
		Float  float64
		Bool   bool
	}

	var instance testStructDeepSet

	DeepSet(&instance, []string{"String"}, `Hello`)
	DeepSet(&instance, []string{"Int"}, 123)
	DeepSet(&instance, []string{"Float"}, 3.14)
	DeepSet(&instance, []string{"Bool"}, true)

	assert.Equal(`Hello`, instance.String)
	assert.Equal(123, instance.Int)
	assert.Equal(3.14, instance.Float)
	assert.True(instance.Bool)
}

func TestDiffuseMap(t *testing.T) {
	var assert = require.New(t)

	var output = make(map[string]any)

	output["name"] = "test.thing.name"
	output["enabled"] = true
	output["cool.beans"] = "yep"
	output["tags.0"] = "base"
	output["tags.1"] = "other"
	output["tags.2"] = "more"
	output["tags.3"] = "still-more"
	output["devices.0.name"] = "lo"
	output["devices.1.name"] = "eth0"
	output["devices.1.peers.0"] = "0.0.0.0"
	output["devices.1.peers.1"] = "1.1.1.1"
	output["devices.1.peers.2"] = "2.2.2.2"
	output["devices.1.switch.0.name"] = "aa:bb:cc:dd:ee:ff"
	output["devices.1.switch.0.ip"] = "111.222.0.1"
	output["devices.1.switch.1.name"] = "cc:dd:ee:ff:bb:dd"
	output["devices.1.switch.1.ip"] = "111.222.0.2"

	output, err := DiffuseMap(output, ".")
	assert.NoError(err)

	//  name
	v, _ := output["name"]
	assert.Equal("test.thing.name", v)

	//  enabled
	v, _ = output["enabled"]
	assert.Equal(true, v)

	//  tags[]
	v, ok := output["tags"]
	assert.True(ok)

	assert.Len(v, 4)

	var vArray = v.([]any)

	assert.Equal("base", vArray[0])
	assert.Equal("other", vArray[1])
	assert.Equal("more", vArray[2])
	assert.Equal("still-more", vArray[3])
}

func TestDelete(t *testing.T) {
	var assert = require.New(t)
	var in = map[string]any{
		`a`: 1,
		`b`: 2,
		`c`: 3,
	}

	assert.NoError(Delete(in, `b`))

	assert.EqualValues(map[string]any{
		`a`: 1,
		`c`: 3,
	}, in)
}
