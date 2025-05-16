package maputil

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ghetzel/testify/require"
)

type MyTestThing struct {
	Name  string
	Other int
}

func TestMapJoin(t *testing.T) {
	var input = map[string]any{
		`key1`: `value1`,
		`key2`: true,
		`key3`: 3,
	}

	var output = Join(input, `=`, `&`)

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
	var assert = require.New(t)

	var i1 = map[string]any{
		`1`: 1,
		`2`: true,
		`3`: `three`,
	}

	var i2 = map[string]bool{
		`1`: true,
		`2`: false,
		`3`: true,
	}

	var i3 = map[string]MyTestThing{
		`1`: MyTestThing{},
		`2`: MyTestThing{},
		`3`: MyTestThing{},
	}

	var i4 sync.Map

	i4.Store(`1`, MyTestThing{})
	i4.Store(`2`, 2)
	i4.Store(`3`, 3.14)

	var output = []string{`1`, `2`, `3`}

	assert.Empty(StringKeys(nil))

	assert.Equal(output, StringKeys(i1))
	assert.Equal(output, StringKeys(i2))
	assert.Equal(output, StringKeys(i3))
	assert.Equal(output, StringKeys(&i4))

	assert.Empty(StringKeys(true))
	assert.Empty(StringKeys(4))
	assert.Empty(StringKeys([]int{1, 2, 3}))
}

func TestMapSplit(t *testing.T) {
	var input = `key1=value1&key2=true&key3=3`

	var output = Split(input, `=`, `&`)

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
	Properties            map[string]any
	StrSliceTest          []string
	InterfaceStrSliceTest []string
	StructSliceTest       []SubtypeTester
	StructSliceTest2      []SubtypeTester
	StructSliceTest3      []SubtypeTester
	nonexported           int
}

func TestStructFromMapEmbedded(t *testing.T) {
	type tAddress struct {
		Number   string
		Street   string
		City     string `potato:"city"`
		State    string `potato:"state"`
		Country  string `potato:"country"`
		LoadedAt time.Time
	}

	type tPerson struct {
		Name    string
		Age     int `potato:"age"`
		Address *tAddress
	}

	type tUser struct {
		tPerson
		Email  string `potato:"email"`
		Active bool   `potato:"ACTIVE"`
	}

	var assert = require.New(t)

	var tgt tUser

	assert.NoError(TaggedStructFromMap(map[string]any{
		`Name`:   `Rusty Shackleford`,
		`age`:    420,
		`email`:  `none+of@your.biz`,
		`ACTIVE`: true,
		`Address`: map[string]any{
			`Number`:   350,
			`Street`:   `Fifth Avenue`,
			`city`:     `New York`,
			`state`:    `NY`,
			`country`:  `US`,
			`LoadedAt`: `2006-01-02`,
		},
	}, &tgt, `potato`))

	// assert.Equal(`Rusty Shackleford`, tgt.Name)
	// assert.Equal(420, tgt.Age)
	assert.Equal(`none+of@your.biz`, tgt.Email)
	assert.True(tgt.Active)
	// assert.Equal(`350`, tgt.Address.Number)
	// assert.Equal(`Fifth Avenue`, tgt.Address.Street)
	// assert.Equal(`New York`, tgt.Address.City)
	// assert.Equal(`NY`, tgt.Address.State)
	// assert.Equal(`US`, tgt.Address.Country)
}

func TestStructFromMap(t *testing.T) {
	var assert = require.New(t)

	var input = map[string]any{
		`Name`:           `Foo Bar`,
		`active`:         true,
		`should_not_set`: 4,
		`subtype2`: map[string]any{
			`A`: 3,
			`b`: 4,
		},
		`TimeTest`: 15000000000,
		`Subtype1`: map[string]any{
			`A`: 1111,
			`b`: 2222,
		},
		`IntTest`: int64(5),
		`Properties`: map[string]any{
			`first`:  1,
			`second`: true,
			`third`:  `three`,
		},
		`StrSliceTest`:          []string{`one`, `two`, `three`},
		`InterfaceStrSliceTest`: []any{`one`, `two`, `three`},
		`StructSliceTest`:       []SubtypeTester{{10, 11}, {12, 13}, {14, 15}},
		`StructSliceTest2`: []map[string]any{
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
		`StructSliceTest3`: []any{
			map[string]any{
				`A`: 10,
				`b`: 11,
			},
			map[string]any{
				`A`: 12,
				`b`: 13,
			},
			map[string]any{
				`A`: 14,
				`b`: 15,
			},
		},
	}

	var output = MyStructTester{}

	var err = StructFromMap(input, &output)
	assert.NoError(err)

	assert.Equal(`Foo Bar`, output.Name)
	assert.True(output.Active)
	assert.Zero(output.nonexported)

	assert.Equal(1111, output.Subtype1.A)
	assert.Equal(2222, output.Subtype1.B)

	assert.NotNil(output.Subtype2)
	assert.Equal(3, output.Subtype2.A)
	assert.Equal(4, output.Subtype2.B)

	assert.Equal(time.Duration(15)*time.Second, output.TimeTest)
	assert.Equal(int32(5), output.IntTest)

	assert.NotNil(output.Properties)
	assert.EqualValues(1, output.Properties[`first`])
	assert.EqualValues(true, output.Properties[`second`])
	assert.Equal(`three`, output.Properties[`third`])

	assert.NotNil(output.StrSliceTest)
	assert.Len(output.StrSliceTest, 3)
	assert.Equal(`one`, output.StrSliceTest[0])
	assert.Equal(`two`, output.StrSliceTest[1])
	assert.Equal(`three`, output.StrSliceTest[2])

	assert.NotNil(output.InterfaceStrSliceTest)
	assert.Len(output.InterfaceStrSliceTest, 3)
	assert.EqualValues(`one`, output.InterfaceStrSliceTest[0])
	assert.EqualValues(`two`, output.InterfaceStrSliceTest[1])
	assert.EqualValues(`three`, output.InterfaceStrSliceTest[2])

	assert.NotNil(output.StructSliceTest)
	assert.Len(output.StructSliceTest, 3)
	assert.EqualValues(10, output.StructSliceTest[0].A)
	assert.EqualValues(11, output.StructSliceTest[0].B)

	assert.EqualValues(12, output.StructSliceTest[1].A)
	assert.EqualValues(13, output.StructSliceTest[1].B)

	assert.EqualValues(14, output.StructSliceTest[2].A)
	assert.EqualValues(15, output.StructSliceTest[2].B)

	assert.NotNil(output.StructSliceTest2)
	assert.Len(output.StructSliceTest2, 3)
	assert.EqualValues(10, output.StructSliceTest2[0].A)
	assert.EqualValues(11, output.StructSliceTest2[0].B)

	assert.EqualValues(12, output.StructSliceTest2[1].A)
	assert.EqualValues(13, output.StructSliceTest2[1].B)

	assert.EqualValues(14, output.StructSliceTest2[2].A)
	assert.EqualValues(15, output.StructSliceTest2[2].B)

	assert.NotNil(output.StructSliceTest3)
	assert.Len(output.StructSliceTest3, 3)
	assert.EqualValues(10, output.StructSliceTest3[0].A)
	assert.EqualValues(11, output.StructSliceTest3[0].B)

	assert.EqualValues(12, output.StructSliceTest3[1].A)
	assert.EqualValues(13, output.StructSliceTest3[1].B)

	assert.EqualValues(14, output.StructSliceTest3[2].A)
	assert.EqualValues(15, output.StructSliceTest3[2].B)
}

func TestMapAppend(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(map[string]any{}, Append())

	assert.Equal(map[string]any{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, Append(map[string]any{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}))

	assert.Equal(map[string]any{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, Append(nil, map[string]any{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, nil))

	assert.Equal(map[string]any{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
		`d`: 4,
		`e`: false,
		`f`: 6.1,
	}, Append(map[string]any{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, map[string]any{
		`d`: 4,
		`e`: false,
		`f`: 6.1,
	}))

	assert.Equal(map[string]any{
		`a`: 1,
		`b`: true,
		`c`: `Five`,
	}, Append(map[string]any{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, map[string]any{
		`c`: `Five`,
	}))
}

// func TestMapValues(t *testing.T) {
// 	assert := require.New(t)

// 	assert.Equal([]any{
// 		1, 3, 5,
// 	}, MapValues(map[string]int{
// 		`first`:  1,
// 		`second`: 3,
// 		`third`:  5,
// 	}))
// }

func TestApply(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(map[string]any{
		`a`: 10,
		`b`: 20,
		`c`: 30,
	}, Apply(map[string]any{
		`a`: 1,
		`b`: 2,
		`c`: 3,
	}, func(_ []string, value any) (any, bool) {
		return value.(int) * 10, true
	}))
}
