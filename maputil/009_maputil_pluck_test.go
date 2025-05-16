package maputil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestMapPluck(t *testing.T) {
	var assert = require.New(t)

	assert.Empty(Pluck(nil, nil))
	assert.Empty(Pluck(nil, []string{`name`}))
	assert.Empty(Pluck(`test`, []string{`name`}))
	assert.Empty(Pluck([]string{`test`, `values`}, []string{`name`}))

	assert.Equal([]any{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]map[string]string{
		map[string]string{
			`name`: `Alice`,
		},
		map[string]string{
			`name`: `Bob`,
		},
		map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))

	assert.Equal([]any{
		`Alice`,
		`Mallory`,
	}, Pluck([]map[string]string{
		map[string]string{
			`name`: `Alice`,
		},
		map[string]string{
			`NAME`: `Bob`,
		},
		map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))

	assert.Equal([]any{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]map[string]map[string]any{
		map[string]map[string]any{
			`info`: map[string]any{
				`name`: `Alice`,
			},
		},
		map[string]map[string]any{
			`info`: map[string]any{
				`name`: `Bob`,
			},
		},
		map[string]map[string]any{
			`info`: map[string]any{
				`name`: `Mallory`,
			},
		},
	}, []string{`info`, `name`}))

	assert.Equal([]any{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]map[any]map[any]any{
		map[any]map[any]any{
			`info`: map[any]any{
				`name`: `Alice`,
			},
		},
		map[any]map[any]any{
			`info`: map[any]any{
				`name`: `Bob`,
			},
		},
		map[any]map[any]any{
			`info`: map[any]any{
				`name`: `Mallory`,
			},
		},
	}, []string{`info`, `name`}))

	assert.Equal([]any{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]any{
		map[string]string{
			`name`: `Alice`,
		},
		map[string]string{
			`name`: `Bob`,
		},
		map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))

	assert.Equal([]any{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]any{
		&map[string]string{
			`name`: `Alice`,
		},
		&map[string]string{
			`name`: `Bob`,
		},
		&map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))

	assert.Equal([]any{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck(&[]any{
		&map[string]string{
			`name`: `Alice`,
		},
		&map[string]string{
			`name`: `Bob`,
		},
		&map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))

	assert.Equal([]any{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck(&[]any{
		[]any{
			&map[string]string{
				`name`: `Alice`,
			},
			&map[string]string{
				`name`: `Bob`,
			},
		},
		[]any{
			&map[string]string{
				`name`: `Mallory`,
			},
		},
	}, []string{`*`, `name`}))
}
