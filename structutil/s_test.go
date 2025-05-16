package structutil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

type tStructBase struct {
	Name    string
	Age     int
	Enabled bool
	hidden  bool
	Words   []string `testing:"WORDS,omitempty"`
	Phrase  []string `structutil:"replace"`
}

type tStructChild struct {
	*tStructBase
	Address string
	Details map[string]any
}

func TestS(t *testing.T) {
	var base = &tStructChild{
		tStructBase: &tStructBase{
			Name:   `test`,
			Age:    42,
			Words:  []string{`hello`, `there`},
			Phrase: []string{`may`, `the`, `force`, `be`, `with`, `you`},
		},
		Address: `123 Fake St.`,
		Details: map[string]any{
			`global`: map[string]any{
				`name`:   `base`,
				`values`: []int{2, 4, 6},
			},
		},
	}

	var s = S(base)

	var fields = s.Fields()

	require.Len(t, fields, 7)

	for i, field := range fields {
		var name = field.Name
		var typn = field.Type.String()

		switch i {
		case 0:
			require.Equal(t, `Name`, name)
			require.Equal(t, `string`, typn)
			require.Equal(t, `test`, field.V().String())
		case 1:
			require.Equal(t, `Age`, name)
			require.Equal(t, `int`, typn)
			require.Equal(t, int64(42), field.V().Int())
		case 2:
			require.Equal(t, `Enabled`, name)
			require.Equal(t, `bool`, typn)
			require.Equal(t, false, field.V().Bool())
		case 3:
			require.Equal(t, `Words`, name)
			require.Equal(t, `[]string`, typn)
			require.Equal(t, []string{`hello`, `there`}, field.V().Strings())

			var tag, attrs, ok = field.GetTag(`testing`)

			require.True(t, ok)
			require.Equal(t, `WORDS`, tag)
			require.Equal(t, []string{
				`omitempty`,
			}, attrs)

			tag, attrs, ok = field.GetTag(`other`)

			require.False(t, ok)
			require.Equal(t, ``, tag)
			require.Empty(t, attrs)
		case 4:
			require.Equal(t, `Phrase`, name)
			require.Equal(t, `[]string`, typn)
			require.Equal(t, []string{`may`, `the`, `force`, `be`, `with`, `you`}, field.V().Strings())
		case 5:
			require.Equal(t, `Address`, name)
			require.Equal(t, `string`, typn)
			require.Equal(t, `123 Fake St.`, field.V().String())
		case 6:
			require.Equal(t, `Details`, name)
			require.Equal(t, `map[string]interface {}`, typn)
			require.Len(t, base.Details, 1)
		}
	}

	require.NoError(t, s.Merge(&tStructChild{
		tStructBase: &tStructBase{
			Words:  []string{`general`, `kenobi`},
			Phrase: []string{`and`, `also`, `with`, `you`},
		},
		Address: `987 Lulz Lane`,
		Details: map[string]any{
			`global`: map[string]any{
				`name`:   `replaced`,
				`values`: []int{1, 3, 5},
			},
		},
	}))

	require.Equal(t, []string{
		`hello`, `there`, `general`, `kenobi`,
	}, base.Words)

	require.Equal(t, []string{
		`and`, `also`, `with`, `you`,
	}, base.Phrase)

	require.Equal(t, map[string]any{
		`global`: map[string]any{
			`name`:   `replaced`,
			`values`: []any{1, 3, 5},
		},
	}, base.Details)
}
