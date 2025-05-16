package maputil

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/testify/require"
)

type testMstruct struct {
	ID     string `maputil:"id" json:"_id"`
	Name   string `json:"NAME"`
	Factor float64
}

func TestM(t *testing.T) {
	var assert = require.New(t)
	var input = M(map[string]any{
		`first`: true,
		`second`: map[string]any{
			`s1`:     `test`,
			`values`: []int{1, 2, 3, 4},
			`truthy`: `True`,
			`strnum`: `42`,
			`then`:   `2006-01-02`,
		},
		`now`:    time.Now(),
		`third`:  3.1415,
		`fourth`: 42,
		`mapslice`: []any{
			map[string]any{
				`id`: 1,
			},
			map[string]any{
				`id`: 2,
			},
			map[string]any{
				`id`: 3,
			},
		},
	})

	assert.Equal(``, M(nil).String(`lol`))
	assert.False(M(nil).Bool(`lol`))
	assert.Equal(int64(0), M(nil).Int(`lol`))
	assert.Equal(float64(0), M(nil).Float(`lol`))
	assert.Len(M(nil).Slice(`lol`), 0)
	assert.Nil(M(nil).Auto(`second.strnum`))
	assert.Zero(M(nil).Time(`now`))

	assert.Equal(`test`, input.String(`second.s1`))
	assert.True(input.Bool(`first`))
	assert.True(input.Bool(`second.truthy`))
	assert.True(input.Bool(`second.s1`))
	assert.Equal(3.1415, input.Float(`third`))
	assert.Equal(int64(3), input.Int(`third`))
	assert.Equal(int64(42), input.Int(`fourth`))
	assert.Equal(int64(3), input.Int(`second.values.2`))
	assert.Equal(int64(0), input.Int(`second.values.99`))
	assert.Equal(float64(42), input.Float(`fourth`))
	assert.Len(input.Slice(`second.values`), 4)
	assert.Equal(int64(42), input.Auto(`second.strnum`))
	assert.Equal(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), input.Time(`second.then`))

	var mapvals []int64

	for _, ms := range input.SliceOfMaps(`mapslice`) {
		mapvals = append(mapvals, ms.Int(`id`))
	}

	assert.Equal([]int64{1, 2, 3}, mapvals)

	assert.Equal(2, input.JSONPath(`$..values[1]`, 42))
	assert.Equal(42, input.JSONPath(`$..values[99]`, 42))

	assert.Equal(6, input.Len())
	var k = make([]string, 6)
	var i = 0

	assert.NoError(input.Each(func(key string, value typeutil.Variant) error {
		k[i] = key
		i++
		return nil
	}))

	assert.ElementsMatch(k, []string{`first`, `second`, `third`, `fourth`, `now`, `mapslice`})
}

func TestMSet(t *testing.T) {
	var assert = require.New(t)
	var input = M(nil)

	assert.Equal(``, input.String(`lol`))

	assert.Equal(`2funny4me`, input.Set(`lol`, `2funny4me`).String())

	assert.Equal(`2funny4me`, input.String(`lol`))
}

func TestMStruct(t *testing.T) {
	var assert = require.New(t)
	var input = M(&testMstruct{
		ID:     `123`,
		Name:   `tester`,
		Factor: 3.14,
	})

	assert.Equal(`123`, input.String(`id`))
	assert.EqualValues(123, input.Int(`id`))

	assert.Equal(`tester`, input.String(`Name`))
	assert.Equal(3.14, input.Float(`Factor`))

	assert.Equal(map[string]any{
		`id`:     `123`,
		`Name`:   `tester`,
		`Factor`: 3.14,
	}, input.MapNative())

	assert.Equal(map[string]any{
		`_id`:    `123`,
		`NAME`:   `tester`,
		`Factor`: 3.14,
	}, input.MapNative(`json`))

	var keys []string
	var values []any

	for item := range input.Iter() {
		keys = append(keys, item.K)
		values = append(values, item.Value)
	}

	assert.ElementsMatch([]string{`id`, `Name`, `Factor`}, keys)
	assert.ElementsMatch([]any{`123`, `tester`, float64(3.14)}, values)
}

func TestMUrlValues(t *testing.T) {
	var assert = require.New(t)
	var input = M(url.Values{
		`a`: []string{`1`},
		`b`: []string{},
		`c`: []string{`2`, `3`},
	})

	assert.Equal(`1`, input.String(`a`))
	assert.EqualValues(1, input.Int(`a`))

	assert.Equal(``, input.String(`b`))
	assert.Equal(float64(0), input.Float(`b`))
	assert.Nil(input.Auto(`b`))

	assert.Equal([]string{`2`, `3`}, input.Strings(`c`))
}

func TestMHttpHeader(t *testing.T) {
	var assert = require.New(t)
	var input = M(http.Header{
		`a`: []string{`1`},
		`b`: []string{},
		`c`: []string{`2`, `3`},
	})

	assert.Equal(`1`, input.String(`a`))
	assert.EqualValues(1, input.Int(`a`))

	assert.Equal(``, input.String(`b`))
	assert.Equal(float64(0), input.Float(`b`))
	assert.Nil(input.Auto(`b`))

	assert.Equal([]string{`2`, `3`}, input.Strings(`c`))
}

func TestMStructNested(t *testing.T) {
	var assert = require.New(t)

	type msecond struct {
		S1     string
		Values []int
		Truthy any
		Strnum string
		Then   string
	}

	type mtop struct {
		First  bool
		Second msecond
		Now    time.Time
		Third  float64
		Fourth int
	}

	var input = M(mtop{
		First: true,
		Second: msecond{
			S1:     `test`,
			Values: []int{1, 2, 3, 4},
			Truthy: `True`,
			Strnum: `42`,
			Then:   `2006-01-02`,
		},
		Now:    time.Now(),
		Third:  3.1415,
		Fourth: 42,
	})

	assert.Equal(`test`, input.String(`Second.S1`))
	assert.True(input.Bool(`First`))
	assert.True(input.Bool(`Second.Truthy`))
	assert.True(input.Bool(`Second.S1`))
	assert.Equal(3.1415, input.Float(`Third`))
	assert.Equal(int64(3), input.Int(`Third`))
	assert.Equal(int64(42), input.Int(`Fourth`))
	assert.Equal(int64(3), input.Int(`Second.Values.2`))
	assert.Equal(int64(0), input.Int(`Second.Values.99`))
	assert.Equal(float64(42), input.Float(`Fourth`))
	assert.Len(input.Slice(`Second.Values`), 4)
	assert.Equal(int64(42), input.Auto(`Second.Strnum`))
	assert.Equal(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), input.Time(`Second.Then`))
}

func TestMMarshalJSON(t *testing.T) {
	var assert = require.New(t)

	var m = M(map[string]any{
		`hello`: 1,
		`there`: true,
		`general`: map[string]any{
			`kenobi`: true,
		},
		`xyz`: []string{`a`, `b`, `c`},
		`zzz`: []map[string]any{
			map[string]any{
				`name`:  `a`,
				`value`: 0,
			},
			map[string]any{
				`name`:  `b`,
				`value`: 1,
			},
			map[string]any{
				`name`:  `c`,
				`value`: 2,
			},
		},
	})

	out, err := json.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`{"general":{"kenobi":true},"hello":1,"there":true,"xyz":["a","b","c"],"zzz":[{"name":"a","value":0},{"name":"b","value":1},{"name":"c","value":2}]}`), out)
}

func TestMUnmarshalJSON(t *testing.T) {
	var assert = require.New(t)

	var js = []byte(`{"general":{"kenobi":true},"hello":1,"there":true,"xyz":["a","b","c"],"zzz":[{"name":"a","value":0},{"name":"b","value":1},{"name":"c","value":2}]}`)

	var m Map

	assert.NoError(json.Unmarshal(js, &m))
	assert.Equal(map[string]any{
		`hello`: float64(1),
		`there`: true,
		`general`: map[string]any{
			`kenobi`: true,
		},
		`xyz`: []any{`a`, `b`, `c`},
		`zzz`: []any{
			map[string]any{
				`name`:  `a`,
				`value`: float64(0),
			},
			map[string]any{
				`name`:  `b`,
				`value`: float64(1),
			},
			map[string]any{
				`name`:  `c`,
				`value`: float64(2),
			},
		},
	}, m.data)

	assert.EqualValues(1, m.Int(`hello`))
	assert.True(m.Bool(`there`))
	assert.Equal(true, m.Bool(`general.kenobi`))
	assert.EqualValues(2, m.Float(`zzz.2.value`))
}

func TestMMarshalXML(t *testing.T) {
	var assert = require.New(t)

	var m = M(map[string]any{
		`hello`: 1,
		`there`: true,
		`general`: map[string]any{
			`kenobi`: true,
		},
		`xyz`: []string{`a`, `b`, `c`},
		`zzz`: []map[string]any{
			map[string]any{
				`name`:  `a`,
				`value`: 0,
			},
			map[string]any{
				`name`:  `b`,
				`value`: 1,
			},
			map[string]any{
				`name`:  `c`,
				`value`: 2,
			},
		},
	})

	// default marshal
	out, err := xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<data><general><kenobi>true</kenobi></general><hello>1</hello><there>true</there><xyz><element>a</element><element>b</element><element>c</element></xyz><zzz><element><name>a</name><value>0</value></element><element><name>b</name><value>1</value></element><element><name>c</name><value>2</value></element></zzz></data>`), out)

	// custom root tagname
	m.SetRootTagName(`nub_nub`)
	out, err = xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<nub_nub><general><kenobi>true</kenobi></general><hello>1</hello><there>true</there><xyz><element>a</element><element>b</element><element>c</element></xyz><zzz><element><name>a</name><value>0</value></element><element><name>b</name><value>1</value></element><element><name>c</name><value>2</value></element></zzz></nub_nub>`), out)

	// generic XML structure
	m.SetMarshalXmlGeneric(true)
	out, err = xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<nub_nub><item type="object" key="general"><item key="kenobi" type="bool">true</item></item><item key="hello" type="int">1</item><item key="there" type="bool">true</item><item type="array" key="xyz"><item key="element" type="str">a</item><item key="element" type="str">b</item><item key="element" type="str">c</item></item><item type="array" key="zzz"><item type="object" key="element"><item key="name" type="str">a</item><item key="value" type="int">0</item></item><item type="object" key="element"><item key="name" type="str">b</item><item key="value" type="int">1</item></item><item type="object" key="element"><item key="name" type="str">c</item><item key="value" type="int">2</item></item></item></nub_nub>`), out)

	// compact structure, custom keyfunc
	m.SetMarshalXmlGeneric(false)
	m.SetMarshalXmlKeyFunc(func(in string) string {
		return stringutil.Camelize(in)
	})

	out, err = xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<NubNub><General><Kenobi>true</Kenobi></General><Hello>1</Hello><There>true</There><Xyz><Element>a</Element><Element>b</Element><Element>c</Element></Xyz><Zzz><Element><Name>a</Name><Value>0</Value></Element><Element><Name>b</Name><Value>1</Value></Element><Element><Name>c</Name><Value>2</Value></Element></Zzz></NubNub>`), out)

	// generic structure, custom keyfunc
	m.SetMarshalXmlGeneric(true)
	out, err = xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<NubNub><item type="object" key="General"><item key="Kenobi" type="bool">true</item></item><item key="Hello" type="int">1</item><item key="There" type="bool">true</item><item type="array" key="Xyz"><item key="Element" type="str">a</item><item key="Element" type="str">b</item><item key="Element" type="str">c</item></item><item type="array" key="Zzz"><item type="object" key="Element"><item key="Name" type="str">a</item><item key="Value" type="int">0</item></item><item type="object" key="Element"><item key="Name" type="str">b</item><item key="Value" type="int">1</item></item><item type="object" key="Element"><item key="Name" type="str">c</item><item key="Value" type="int">2</item></item></item></NubNub>`), out)
}

func TestMIter(t *testing.T) {
	var assert = require.New(t)
	var input = M(map[string]any{
		`a`: 1,
		`b`: 2,
		`c`: 3,
		`d`: 4,
	})

	var keys []string
	var vals []int

	for item := range input.Iter() {
		keys = append(keys, item.K)
		vals = append(vals, int(item.V.Int()))
	}

	assert.ElementsMatch([]string{`c`, `d`, `a`, `b`}, keys)
	assert.ElementsMatch([]int{4, 2, 3, 1}, vals)

	keys = nil
	vals = nil

	for item := range input.Iter(IterOptions{
		SortKeys: true,
	}) {
		keys = append(keys, item.K)
		vals = append(vals, int(item.V.Int()))
	}

	assert.Equal([]string{`a`, `b`, `c`, `d`}, keys)
	assert.Equal([]int{1, 2, 3, 4}, vals)
}

func TestMJson(t *testing.T) {
	var assert = require.New(t)
	var m = M(`{"hello": "there", "general": "kenobi"}`)

	assert.Equal(map[string]any{
		`hello`:   `there`,
		`general`: `kenobi`,
	}, m.MapNative())
	assert.Equal(`there`, m.String(`hello`))
	assert.Equal(`kenobi`, m.String(`general`))

	m.Delete(`general`)

	assert.Equal(map[string]any{
		`hello`: `there`,
	}, m.MapNative())

	assert.Equal([]byte(`{"hello":"there"}`), m.JSON())
	assert.Equal([]byte("{\n..\"hello\": \"there\"\n}"), m.JSON(`..`))
}

func TestMMerge(t *testing.T) {
	var m = NewMap()

	assert.Equal(t, map[string]any{}, m.MapNative())

	m.Merge(map[string]any{
		`a`: 1,
		`b`: 2,
	})

	assert.Equal(t, map[string]any{
		`a`: 1,
		`b`: 2,
	}, m.MapNative())

	m.Merge(map[string]any{
		`b`: 2.1,
		`c`: 3,
	})

	assert.Equal(t, map[string]any{
		`a`: 1,
		`b`: 2.1,
		`c`: 3,
	}, m.MapNative())

	m.Merge(map[string]any{
		`b`: nil,
	})

	assert.Equal(t, map[string]any{
		`a`: 1,
		`b`: nil,
		`c`: 3,
	}, m.MapNative())

	m.Compact()

	assert.Equal(t, map[string]any{
		`a`: 1,
		`c`: 3,
	}, m.MapNative())
}
