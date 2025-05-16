package stringutil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestUUID(t *testing.T) {
	var assert = require.New(t)
	var input = []byte{
		0x01, 0x02, 0x03, 0x01,
		0x02, 0x03, 0x01, 0x02,
		0x03, 0x01, 0x02, 0x03,
		0x01, 0x02, 0x03, 0x01,
	}

	uuid, err := UuidFromBytes(input)

	assert.NoError(err)

	assert.Equal(`01020301-0203-0102-0301-020301020301`, uuid.String())
	assert.Equal(input, uuid.Bytes())
	assert.Equal(`01020301020301020301020301020301`, uuid.Hex())
	assert.Equal(`AQIDAQIDAQIDAQIDAQIDAQ==`, uuid.Base64())
	assert.Equal(`8DfbUyTr2zZABVZdbmdo6`, uuid.Base58())
}
