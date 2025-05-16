package rxutil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestMatch(t *testing.T) {
	var assert = require.New(t)

	var match = Match(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`, `1234.5678.9`)
	assert.NotNil(match)

	assert.Equal(`1234.5678.9`, match.Group(0))
	assert.Equal(`1234`, match.Group(1))
	assert.Equal(`5678`, match.Group(2))
	assert.Equal(`9`, match.Group(3))
	assert.Empty(match.Group(4))

	assert.Equal(`1234`, match.Group(`first`))
	assert.Equal(`9`, match.Group(`second`))
	assert.Empty(match.Group(`potato`))
}

func TestMatchAndMap(t *testing.T) {
	var assert = require.New(t)

	var match = Match(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`, `1234.5678.9`)
	assert.NotNil(match)

	assert.Equal(`1234.5678.9`, match.Group(0))
	assert.Equal(map[string]string{
		`first`:  `1234`,
		`second`: `9`,
	}, match.NamedCaptures())
}

func TestMatchCaptures(t *testing.T) {
	var assert = require.New(t)

	var match = Match(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`, `1234.5678.9`)
	assert.NotNil(match)

	assert.Equal([]string{`1234.5678.9`, `1234`, `5678`, `9`}, match.Captures())
}

func TestReplaceGroup(t *testing.T) {
	var assert = require.New(t)

	var match = Match(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`, `1234.5678.9`)
	assert.NotNil(match)

	assert.Equal(`repl`, match.ReplaceGroup(0, `repl`))
	assert.Equal(`first.5678.9`, match.ReplaceGroup(1, `first`))
	assert.Equal(`1234.second.9`, match.ReplaceGroup(2, `second`))
	assert.Equal(`1234.5678.third`, match.ReplaceGroup(3, `third`))
	assert.Equal(`1234.5678.9`, match.ReplaceGroup(4, `fourth`))
}

func TestReplaceGroupNamed(t *testing.T) {
	var assert = require.New(t)

	var match = Match(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`, `1234.5678.9`)
	assert.NotNil(match)

	assert.Equal(`first.5678.9`, match.ReplaceGroup(`first`, `first`))
	assert.Equal(`1234.5678.second`, match.ReplaceGroup(`second`, `second`))
	assert.Equal(`1234.5678.9`, match.ReplaceGroup(`third`, `third`))
}
