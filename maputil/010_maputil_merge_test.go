package maputil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestMapMerge(t *testing.T) {
	var assert = require.New(t)

	out, err := Merge(nil, nil)
	assert.NoError(err)
	assert.Empty(out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`: `First`,
	}, map[string]any{
		`age`: 2,
	})

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`: `First`,
		`age`:  2,
	}, out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`: []string{`First`, `Second`},
	}, map[string]any{
		`name`: `Third`,
	})

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`: []any{`First`, `Second`, `Third`},
	}, out)

	// ------------------------------------------------------------------------
	// FIXME: wrong output
	//	expected: map[string]interface {}{"name":[]interface {}{"First", "Second", "Third"}}
	//	received: map[string]interface {}{"name":[]interface {}{"First", "Second", []interface {}{"First", "Third"}}}
	//
	// out, err = Merge(map[string]any{
	// 	`name`: []string{`First`, `Second`},
	// }, map[string]any{
	// 	`name`: []string{`Third`},
	// })

	// assert.NoError(err)
	// assert.Equal(map[string]any{
	// 	`name`: []any{`First`, `Second`, `Third`},
	// }, out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`: `First`,
	}, map[string]any{
		`name`: `First`,
	})

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`: `First`,
	}, out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`:    `First`,
		`enabled`: true,
	}, nil)

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`:    `First`,
		`enabled`: true,
	}, out)

	// ------------------------------------------------------------------------
	out, err = Merge(nil, map[string]any{
		`name`:    `Second`,
		`enabled`: true,
	})

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`:    `Second`,
		`enabled`: true,
	}, out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`: `First`,
	}, map[string]any{
		`name`: `Second`,
		`age`:  2,
	})

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`: `Second`,
		`age`:  2,
	}, out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`: `First`,
	}, map[string]any{
		`name`: `Second`,
		`age`:  2,
	}, AppendValues)

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`: []any{`First`, `Second`},
		`age`:  2,
	}, out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`:    `First`,
		`enabled`: nil,
	}, map[string]any{
		`name`:    `Second`,
		`enabled`: true,
	})

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`:    `Second`,
		`enabled`: true,
	}, out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`:    `First`,
		`enabled`: nil,
	}, map[string]any{
		`name`:    `Second`,
		`enabled`: true,
	}, AppendValues)

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`:    []any{`First`, `Second`},
		`enabled`: true,
	}, out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`: `First`,
		`age`:  `yes`,
	}, map[string]any{
		`name`: `Second`,
		`age`:  42,
	})

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`: `Second`,
		`age`:  42,
	}, out)

	// ------------------------------------------------------------------------
	out, err = Merge(map[string]any{
		`name`: `First`,
		`age`:  `yes`,
	}, map[string]any{
		`name`: `Second`,
		`age`:  42,
	}, AppendValues)

	assert.NoError(err)
	assert.Equal(map[string]any{
		`name`: []any{`First`, `Second`},
		`age`:  []any{`yes`, 42},
	}, out)
}
