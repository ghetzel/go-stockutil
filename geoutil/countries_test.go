package geoutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountriesGet(t *testing.T) {
	assert.Equal(t, `United States of America`, Countries.Get(UnitedStates).Name)
	assert.Equal(t, `United States of America`, Countries.Get(`us`).Name)
	assert.Equal(t, `United States of America`, Countries.Get(`US`).Name)
	assert.Equal(t, `United States of America`, Countries.Get(`uS`).Name)
	assert.Equal(t, `United States of America`, Countries.Get(`Us`).Name)

	assert.False(t, Countries.Get(`zz`).IsValid())
}
