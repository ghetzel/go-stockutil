package httputil

import (
	"net/url"
	"testing"

	"github.com/ghetzel/testify/require"
)

type testOutputOne struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Count         int64
	Ok1           bool             `json:"OK1" maputil:"ok1"`
	Ok2           bool             `json:"OK2" maputil:"Ok2"`
	Ok3           bool             `json:"OK3"`
	LOL           []string         `json:"lol"`
	NonIndexedLOL []string         `json:"nilol"`
	Onesie        []string         `json:"onesie"`
	Twosie        []map[string]any `json:"twosie"`
	Empty         string           `json:"empty"`
}

func TestParseFormValues(t *testing.T) {
	var assert = require.New(t)

	var t1 = testOutputOne{
		URL: `http://test`,
	}

	assert.NoError(ParseFormValues(url.Values{
		`name`:             []string{`Tester`},
		`Count`:            []string{`42`},
		`ok1`:              []string{`true`},
		`Ok2`:              []string{`on`},
		`OK3`:              []string{`off`},
		`lol.0`:            []string{`zero`},
		`lol.1`:            []string{`one`},
		`lol.2`:            []string{`two`},
		`nilol[]`:          []string{`first`, `second`, `third`},
		`onesie[]`:         []string{`uno`},
		`twosie[0][hello]`: []string{`there`},
		`empty`:            nil,
	}, &t1))

	assert.Equal(``, t1.Empty)
	assert.Equal(`Tester`, t1.Name)
	assert.Equal(`http://test`, t1.URL)
	assert.Equal(int64(42), t1.Count)
	assert.Equal([]string{`zero`, `one`, `two`}, t1.LOL)
	assert.Equal([]string{`first`, `second`, `third`}, t1.NonIndexedLOL)
	assert.Equal([]string{`uno`}, t1.Onesie)
	assert.Equal([]map[string]any{
		{
			`hello`: `there`,
		},
	}, t1.Twosie)
	assert.True(t1.Ok1)
	assert.True(t1.Ok2)
	assert.False(t1.Ok3)
}
