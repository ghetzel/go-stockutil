package mathutil

import (
	"testing"

	"github.com/ghetzel/testify/assert"
)

func TestClamp(t *testing.T) {
	var assert = assert.New(t)

	assert.True(Clamp(0, 0, 0) == 0)
	assert.True(Clamp(0, -1, 1) == 0)
	assert.True(Clamp(-1, -1, 1) == -1)
	assert.True(Clamp(-0.9999, -1, 1) == -0.9999)
	assert.True(Clamp(-1.0001, -1, 1) == -1)
	assert.True(Clamp(-2, -1, 1) == -1)
	assert.True(Clamp(1, -1, 1) == 1)
	assert.True(Clamp(0.9999, -1, 1) == 0.9999)
	assert.True(Clamp(1.0001, -1, 1) == 1)
	assert.True(Clamp(2, -1, 1) == 1)
}

func TestRound(t *testing.T) {
	var assert = assert.New(t)

	assert.Equal(float64(1), Round(1.499999999999999))
	assert.Equal(float64(2), Round(1.5))
	assert.Equal(float64(2), Round(1.999999999999999))
	assert.Equal(float64(2), Round(1.999999999999999))

	assert.Equal(float64(2), RoundPlaces(2.49, 0))
	assert.Equal(float64(2.5), RoundPlaces(2.49, 1))
	assert.Equal(float64(2.49), RoundPlaces(2.49, 2))
	assert.Equal(float64(2.49), RoundPlaces(2.490000000049, 10))
	assert.Equal(float64(2.4900000001), RoundPlaces(2.49000000005, 10))
}

func TestLeadingSignificantZeros(t *testing.T) {
	var assert = assert.New(t)

	assert.Equal(0, LeadingSignificantZeros(1.1, 10))
	assert.Equal(1, LeadingSignificantZeros(1.0100, 10))
	assert.Equal(2, LeadingSignificantZeros(1.0010, 10))
	assert.Equal(3, LeadingSignificantZeros(1.0001, 10))
	assert.Equal(4, LeadingSignificantZeros(1.00001, 10))
	assert.Equal(5, LeadingSignificantZeros(1.000001, 10))
	assert.Equal(0, LeadingSignificantZeros(1.00000000000001, 12))
	assert.Equal(13, LeadingSignificantZeros(1.00000000000001, 13))
}
