package convutil

import (
	"testing"

	"github.com/ghetzel/testify/assert"
)

func TestConvert(t *testing.T) {
	var assert = assert.New(t)

	// temperatures
	assert.EqualValues(0, MustConvert(32, Fahrenheit, Celsius))
	assert.EqualValues(0, MustConvert(-459.67, Fahrenheit, Kelvin))
	assert.EqualValues(32, MustConvert(0, Celsius, Fahrenheit))
	assert.EqualValues(273.15, MustConvert(0, Celsius, Kelvin))
	assert.EqualValues(-273.15, MustConvert(0, Kelvin, Celsius))
	assert.EqualValues(-459.67, MustConvert(0, Kelvin, Fahrenheit))

	// lengths
	assert.EqualValues(1609.344, MustConvert(1, Miles, Meters))
	assert.EqualValues(1, MustConvert(1609.344, Meters, Miles))
	assert.EqualValues(4200, MustConvert(2.609759007397, Miles, Meters))
	assert.EqualValues(2.609759000000, MustConvert(4200, Meters, Miles))
	assert.EqualValues(1, MustConvert(9.4607304725808e+15, Meters, Lightyears))
	assert.EqualValues(9.4607304725808e+15, MustConvert(1, Lightyears, Meters))
}

func TestByteStringer(t *testing.T) {
	var assert = assert.New(t)

	assert.Equal(`1B`, Byte.String())
	assert.Equal(`1KB`, Kilobyte.String())
	assert.Equal(`1MB`, Megabyte.String())
	assert.Equal(`1GB`, Gigabyte.String())
	assert.Equal(`1TB`, Terabyte.String())
	assert.Equal(`1PB`, Petabyte.String())
	assert.Equal(`1EB`, Exabyte.String())
	assert.Equal(`1ZB`, Zettabyte.String())
	assert.Equal(`1YB`, Yottabyte.String())
	assert.Equal(`1RB`, Ronnabyte.String())
	assert.Equal(`1QB`, Quettabyte.String())

	assert.Equal(`512B`, (512 * Byte).String())
	assert.Equal(`512KB`, (512 * Kilobyte).String())
	assert.Equal(`512MB`, (512 * Megabyte).String())
	assert.Equal(`512GB`, (512 * Gigabyte).String())
	assert.Equal(`512TB`, (512 * Terabyte).String())
	assert.Equal(`512PB`, (512 * Petabyte).String())
	assert.Equal(`512EB`, (512 * Exabyte).String())
	assert.Equal(`512ZB`, (512 * Zettabyte).String())
	assert.Equal(`512YB`, (512 * Yottabyte).String())
	assert.Equal(`512RB`, (512 * Ronnabyte).String())
	assert.Equal(`512QB`, (512 * Quettabyte).String())

	assert.Equal(`0.5KB`, Bytes(512).To(Kilobyte))
	assert.Equal(`0.0005MB`, Bytes(512).To(Megabyte))
	assert.Equal(`0.0000005GB`, Bytes(512).To(Gigabyte))
	assert.Equal(`0.0000000005TB`, Bytes(512).To(Terabyte))
	assert.Equal(`0.0000000000005PB`, Bytes(512).To(Petabyte))
	assert.Equal(`0.0000000000000004EB`, Bytes(512).To(Exabyte))
	// assert.Equal(`0.0000000000000000005ZB`, Bytes(512).To(Zettabyte))
	// assert.Equal(`0.0000000000000000000005YB`, Bytes(512).To(Yottabyte))
	// assert.Equal(`0.0000000000000000000000005BB`, Bytes(512).To(Brontobyte))
}
