package sliceutil

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/ghetzel/testify/require"
)

func TestContains(t *testing.T) {
	var assert = require.New(t)
	var input any

	input = []int{1, 3, 5}
	assert.True(Contains(input, 1))
	assert.True(Contains(input, 3))
	assert.True(Contains(input, 5))
	assert.False(Contains(input, -1))
	assert.False(Contains(input, 2))
	assert.False(Contains(input, -3))
	assert.False(Contains(input, 4))
	assert.False(Contains(input, -5))
	assert.False(Contains([]int{}, 1))
	assert.False(Contains([]int{}, 2))
	assert.False(Contains([]int{}, 0))

	input = []string{"one", "three", "five"}
	assert.True(Contains(input, "one"))
	assert.True(Contains(input, "three"))
	assert.True(Contains(input, "five"))
	assert.False(Contains(input, "One"))
	assert.False(Contains(input, "two"))
	assert.False(Contains(input, "Three"))
	assert.False(Contains(input, "four"))
	assert.False(Contains(input, "Five"))
	assert.False(Contains([]string{}, "one"))
	assert.False(Contains([]string{}, "two"))
	assert.False(Contains([]string{}, ""))
}

func TestAt(t *testing.T) {
	var assert = require.New(t)
	var input any
	var out any
	var ok bool

	input = []int{1, 3, 5}
	out, ok = At(input, 0)
	assert.True(ok)
	assert.Equal(1, out)

	out, ok = At(input, 1)
	assert.True(ok)
	assert.Equal(3, out)

	out, ok = At(input, 2)
	assert.True(ok)
	assert.Equal(5, out)

	out, ok = At(input, 99999)
	assert.False(ok)
	assert.Nil(out)
}

func TestLen(t *testing.T) {
	var assert = require.New(t)
	var input any

	assert.Zero(Len(nil))
	assert.Zero(Len(input))
	input = []int{1, 3, 5}
	assert.Equal(3, Len(input))
	assert.Equal(3, Len(`123`))
}

func TestGet(t *testing.T) {
	var assert = require.New(t)
	var input any

	input = []int{1, 3, 5}
	assert.Equal(1, Get(input, 0))
	assert.Equal(3, Get(input, 1))
	assert.Equal(5, Get(input, 2))
	assert.Nil(Get(input, 99999))
	assert.Nil(Get(nil, 0))
}

func TestFirst(t *testing.T) {
	var assert = require.New(t)
	var input any

	assert.Nil(First(nil))
	assert.Nil(First(input))

	input = []int{}
	assert.Nil(First(input))

	input = []int{1, 3, 5}
	assert.Equal(1, First(input))

}

func TestRest(t *testing.T) {
	var assert = require.New(t)
	var input any

	assert.Nil(Rest(nil))
	assert.Nil(Rest(input))

	input = []int{1}
	assert.Nil(Rest(input))

	input = []int{1, 3, 5}
	assert.Equal([]any{3, 5}, Rest(input))
}

func TestLast(t *testing.T) {
	var assert = require.New(t)
	var input any

	assert.Nil(Last(nil))
	assert.Nil(Last(input))

	input = []int{}
	assert.Nil(Last(input))

	input = []int{1, 3, 5}
	assert.Equal(5, Last(input))
}

func TestContainsString(t *testing.T) {
	var assert = require.New(t)

	var input = []string{"one", "three", "five"}

	assert.True(ContainsString(input, "one"))
	assert.True(ContainsString(input, "three"))
	assert.True(ContainsString(input, "five"))
	assert.False(ContainsString(input, "One"))
	assert.False(ContainsString(input, "two"))
	assert.False(ContainsString(input, "Three"))
	assert.False(ContainsString(input, "four"))
	assert.False(ContainsString(input, "Five"))
	assert.False(ContainsString([]string{}, "one"))
	assert.False(ContainsString([]string{}, "two"))
	assert.False(ContainsString([]string{}, ""))
}

func TestContainsAnyString(t *testing.T) {
	var assert = require.New(t)

	var input = []string{"one", "three", "five"}
	var any = []string{"one", "two", "four"}

	assert.True(ContainsAnyString(input, any...))
	assert.False(ContainsAnyString(input))
	assert.False(ContainsAnyString([]string{}, "one"))
	assert.False(ContainsAnyString([]string{}, "two"))
	assert.False(ContainsAnyString([]string{}, ""))
	assert.False(ContainsAnyString(input, []string{"six", "seven"}...))
}

func TestContainsAllStrings(t *testing.T) {
	var assert = require.New(t)

	var input = []string{"one", "three", "five"}

	assert.True(ContainsAllStrings(input, "one"))
	assert.True(ContainsAllStrings(input, "three"))
	assert.True(ContainsAllStrings(input, "five"))
	assert.True(ContainsAllStrings(input, "one", "three"))
	assert.True(ContainsAllStrings(input, "one", "five"))
	assert.True(ContainsAllStrings(input, "one", "three", "five"))
	assert.False(ContainsAllStrings(input, "one", "four"))
	assert.True(ContainsAllStrings(input))
}

func TestCompact(t *testing.T) {
	var assert = require.New(t)

	assert.Nil(Compact(nil))

	assert.Equal([]any{
		0, 1, 2, 3,
	}, Compact([]any{
		0, 1, 2, 3,
	}))

	assert.Equal([]any{
		1, 3, 5,
	}, Compact([]any{
		nil, 1, nil, 3, nil, 5,
	}))

	assert.Equal([]any{
		`one`, `three`, ` `, `five`,
	}, Compact([]any{
		`one`, ``, `three`, ``, ` `, `five`,
	}))

	assert.Equal([]any{
		[]int{1, 2, 3},
	}, Compact([]any{
		nil, []string{}, []int{1, 2, 3}, map[string]bool{},
	}))
}

func TestCompactString(t *testing.T) {
	var assert = require.New(t)

	assert.Nil(CompactString(nil))

	assert.Equal([]string{
		`one`, `three`, `five`,
	}, CompactString([]string{
		`one`, `three`, `five`,
	}))

	assert.Equal([]string{
		`one`, `three`, ` `, `five`,
	}, CompactString([]string{
		`one`, ``, `three`, ``, ` `, `five`,
	}))
}

func TestStringify(t *testing.T) {
	var assert = require.New(t)

	assert.Nil(Stringify(nil))

	assert.Equal([]string{
		`0`, `1`, `2`,
	}, Stringify([]any{
		0, 1, 2,
	}))

	assert.Equal([]string{
		`0.5`, `0.55`, `0.555`, `0.555001`,
	}, Stringify([]any{
		0.5, 0.55, 0.55500, 0.555001,
	}))

	assert.Equal([]string{
		`true`, `true`, `false`,
	}, Stringify([]any{
		true, true, false,
	}))
}

func TestOr(t *testing.T) {
	var assert = require.New(t)

	assert.Nil(Or())
	assert.Nil(Or(nil))
	assert.Equal(1, Or(0, 1, 0, 2, 0, 3, 4, 5, 6))
	assert.Equal(true, Or(false, false, true))
	assert.Equal(`one`, Or(`one`))
	assert.Equal(4.0, Or(nil, ``, false, 0, 4.0))
	assert.Nil(Or(false, false, false))
	assert.Nil(Or(0, 0, 0))

	assert.Equal(`three`, Or(``, ``, `three`))

	type testStruct struct {
		name string
	}

	assert.Equal(testStruct{`three`}, Or(testStruct{}, testStruct{}, testStruct{`three`}))
}

func TestOrString(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(``, OrString())
	assert.Equal(``, OrString(``))

	assert.Equal(`one`, OrString(`one`))
	assert.Equal(`two`, OrString(``, `two`, ``, `three`))
}

func TestEach(t *testing.T) {
	var assert = require.New(t)

	type testStruct struct {
		Name  string
		Hello bool
		unex  string
	}

	assert.Nil(Each(nil, nil))

	assert.Nil(Each([]string{`one`, `two`, `three`}, func(i int, v any) error {
		return Stop
	}))

	var count = 0
	assert.Nil(Each([]string{`one`, `two`, `three`}, func(i int, v any) error {
		if v.(string) == `two` {
			return Stop
		} else {
			count += 1
			return nil
		}
	}))

	var values = []any{}

	assert.Nil(Each(map[string]string{
		`one`:   `first`,
		`two`:   `second`,
		`three`: `third`,
	}, func(i int, v any) error {
		values = append(values, v)
		return nil
	}))

	values = []any{}

	assert.Nil(Each(&testStruct{
		Name:  `test`,
		Hello: true,
		unex:  `should not see me`,
	}, func(i int, v any) error {
		values = append(values, v)
		return nil
	}))

	assert.ElementsMatch([]any{`test`, true}, values)

	var valchan = make(chan string)

	go func() {
		defer close(valchan)
		for i := 0; i < 3; i++ {
			valchan <- fmt.Sprintf("value%d", i)
		}
	}()

	// test Each-ing over a channel
	var valuesS = make([]string, 0)

	assert.Nil(Each(valchan, func(i int, v any) error {
		valuesS = append(valuesS, fmt.Sprintf("%v", v))
		return nil
	}))

	assert.ElementsMatch([]string{`value0`, `value1`, `value2`}, valuesS)
}

func TestUnique(t *testing.T) {
	var assert = require.New(t)

	assert.Equal([]any{`one`, `two`, `three`}, Unique([]string{`one`, `one`, `two`, `three`}))
	assert.Equal([]any{1, 2, 3}, Unique([]int{1, 2, 2, 3}))
	assert.NotEqual([]any{1, 2, 3}, Unique([]int64{1, 2, 2, 3}))
}

func TestMap(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(
		[]any{10, 20, 30},
		Map([]int{1, 2, 3}, func(_ int, v any) any {
			return v.(int) * 10
		}),
	)

	assert.Equal(
		[]any{true, false, true},
		Map([]bool{false, true, false}, func(_ int, v any) any {
			return !v.(bool)
		}),
	)
}

func TestMapString(t *testing.T) {
	var assert = require.New(t)

	assert.Equal(
		[]string{`1-1thousand`, `2-1thousand`, `3-1thousand`},
		MapString([]int{1, 2, 3}, func(_ int, v string) string {
			return v + `-1thousand`
		}),
	)

	assert.Equal(
		[]string{`first`, `third`, `fifth`},
		CompactString(MapString([]string{`first`, `second`, `third`, `fourth`, `fifth`}, func(_ int, v string) string {
			switch v {
			case `second`, `fourth`:
				return ``
			default:
				return v
			}
		})),
	)
}

func TestChunks(t *testing.T) {
	var assert = require.New(t)
	var input = []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23}

	assert.Equal([][]any{
		[]any{1},
		[]any{3},
		[]any{5},
		[]any{7},
		[]any{9},
		[]any{11},
		[]any{13},
		[]any{15},
		[]any{17},
		[]any{19},
		[]any{21},
		[]any{23},
	}, Chunks(input, 1))

	assert.Equal([][]any{
		[]any{1, 3},
		[]any{5, 7},
		[]any{9, 11},
		[]any{13, 15},
		[]any{17, 19},
		[]any{21, 23},
	}, Chunks(input, 2))

	assert.Equal([][]any{
		[]any{1, 3, 5},
		[]any{7, 9, 11},
		[]any{13, 15, 17},
		[]any{19, 21, 23},
	}, Chunks(input, 3))

	assert.Equal([][]any{
		[]any{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23},
	}, Chunks(input, 1000))
}

func TestFlatten(t *testing.T) {
	var assert = require.New(t)

	assert.Equal([]any{`one`, `two`, `three`}, Flatten([]string{`one`, `two`, `three`}))
	assert.Equal([]any{`one`, `two`, `three`}, Flatten([]any{[]string{`one`, `two`}, `three`}))
	assert.Equal([]any{`one`, `two`, `three`}, Flatten([]any{[]string{`one`}, []string{`two`}, []string{`three`}}))
}

func TestIntersect(t *testing.T) {
	var assert = require.New(t)

	assert.Empty(IntersectStrings(nil, nil))
	assert.Empty(IntersectStrings([]string{`a`, `b`, `c`}, nil))
	assert.Empty(IntersectStrings(nil, []string{`a`, `c`, `e`}))

	assert.Equal(
		[]string{`a`, `c`},
		IntersectStrings([]string{`a`, `b`, `c`}, []string{`a`, `c`, `e`}),
	)
}

func TestDifference(t *testing.T) {
	var assert = require.New(t)

	assert.Empty(Difference(nil, nil))
	assert.Empty(Difference(nil, []string{`a`, `c`, `e`}))
	assert.ElementsMatch([]string{`a`, `b`, `c`}, Difference([]string{`a`, `b`, `c`}, nil))

	assert.ElementsMatch(
		[]any{`a`},
		Difference([]string{`a`, `b`, `c`}, []string{`b`, `c`}),
	)

	assert.ElementsMatch([]string{`a`, `b`, `c`},
		Difference([]string{`a`, `b`, `c`}, []string{`x`, `y`, `z`}),
	)
}

func TestSlice(t *testing.T) {
	var assert = require.New(t)
	var in = []any{1, 2, 3, 4, 5}

	assert.EqualValues([]any{}, Slice(in, 99, -1))
	assert.EqualValues([]any{1}, Slice(in, 0, 1))
	assert.EqualValues([]any{1, 2}, Slice(in, 0, 2))
	assert.EqualValues([]any{1, 2, 3, 4}, Slice(in, 0, 4))
	assert.EqualValues([]any{3, 4}, Slice(in, 2, 4))
	assert.EqualValues([]any{3, 4}, Slice(in, 2, 4))
	assert.EqualValues([]any{4, 5}, Slice(in, -2, -1))
	assert.EqualValues([]any{1, 2, 3, 4, 5}, Slice(in, -5, -1))
	assert.EqualValues([]any{2, 3, 4}, Slice(in, -4, -2))

	assert.EqualValues([]any{}, Slice(in, -6, -6))
	assert.EqualValues([]any{1}, Slice(in, -5, -5))
	assert.EqualValues([]any{2}, Slice(in, -4, -4))
	assert.EqualValues([]any{3}, Slice(in, -3, -3))
	assert.EqualValues([]any{4}, Slice(in, -2, -2))
	assert.EqualValues([]any{5}, Slice(in, -1, -1))
	assert.EqualValues([]any{1}, Slice(in, 0, 1))
	assert.EqualValues([]any{2}, Slice(in, 1, 2))
	assert.EqualValues([]any{3}, Slice(in, 2, 3))
	assert.EqualValues([]any{4}, Slice(in, 3, 4))
	assert.EqualValues([]any{5}, Slice(in, 4, 5))
	assert.EqualValues([]any{}, Slice(in, 5, 6))

	assert.EqualValues([]any{1, 2, 3, 4, 5}, Slice(in, -100, -1))
	assert.EqualValues([]any{1, 2, 3, 4}, Slice(in, -100, -2))
	assert.EqualValues([]any{1, 2, 3}, Slice(in, -100, -3))
	assert.EqualValues([]any{1, 2}, Slice(in, -100, -4))
	assert.EqualValues([]any{1}, Slice(in, -100, -5))
	assert.EqualValues([]any{}, Slice(in, -100, -6))

	assert.EqualValues([]any{2, 3, 4, 5}, Slice(in, 1, -1))
	assert.EqualValues([]any{2, 3, 4}, Slice(in, 1, -2))
	assert.EqualValues([]any{2, 3}, Slice(in, 1, -3))
	assert.EqualValues([]any{2}, Slice(in, 1, -4))
	assert.EqualValues([]any{}, Slice(in, 1, -5))
	assert.EqualValues([]any{}, Slice(in, 1, -6))
}

func TestTrimSpace(t *testing.T) {
	var assert = require.New(t)

	assert.Nil(TrimSpace(nil))
	assert.Equal([]string{`aaa`, `bbb`, `ccc`}, TrimSpace([]string{`aaa`, `   bbb `, ` ccc    `}))
}

func TestFirstNonZero(t *testing.T) {
	assert.Nil(t, FirstNonZero())
	assert.Equal(t, 42, FirstNonZero(42))
	assert.Equal(t, 42, FirstNonZero(``, 0, false, 42, false, true, 96))
	assert.Equal(t, 8, FirstNonZero([]int{0, 0, 0}, 8, []int{69}))
	assert.Equal(t, 84, FirstNonZero([]int{0, 0, 0}, 0, []int{84}))
}

func TestSplitCompact(t *testing.T) {
	assert.Equal(t, []string{}, SplitCompact(``, `,`))
	assert.Equal(t, []string{` `}, SplitCompact(` `, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitCompact(`a,b,c`, `,`))
	assert.Equal(t, []string{` a `, `  b  `, `   c   `}, SplitCompact(` a ,  b  ,   c   `, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitCompact(`a,,b,c`, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitCompact(`a,,b,c,`, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitCompact(`,,,a,,,b,,,,,c,,,`, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitCompact(`a,,,b,,,,,c`, `,`))
}

func TestSplitTrimSpaceCompact(t *testing.T) {
	assert.Equal(t, []string{}, SplitTrimSpaceCompact(``, `,`))
	assert.Equal(t, []string{}, SplitTrimSpaceCompact(` `, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitTrimSpaceCompact(`a,b,c`, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitTrimSpaceCompact(` a ,  b  ,   c   `, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitTrimSpaceCompact(`a,,b,c`, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitTrimSpaceCompact(`a,,b,c,`, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitTrimSpaceCompact(`,,,a,,,b,,,,,c,,,`, `,`))
	assert.Equal(t, []string{`a`, `b`, `c`}, SplitTrimSpaceCompact(`a,,,b,,,,,c`, `,`))
}
