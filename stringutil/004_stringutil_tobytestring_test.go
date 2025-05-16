package stringutil

import (
	"testing"
)

func TestToByteString(t *testing.T) {
	var testvalues = map[any]string{
		0:       `0B`,
		1:       `1B`,
		1023:    `1023B`,
		1024:    `1KB`,
		1536:    `1.5KB`,
		2048:    `2KB`,
		1048575: `1023.9990234375KB`,
		1048576: `1MB`,
	}

	for in, out := range testvalues {
		if v, err := ToByteString(in); err != nil || v != out {
			t.Errorf("Value %v (%T) ToByteString failed: expected '%s', got '%s' (err: %v)", in, in, out, v, err)
		}
	}

	testvalues = map[any]string{
		0:       `0.00B`,
		1:       `1.00B`,
		1023:    `1023.00B`,
		1024:    `1.00KB`,
		1536:    `1.50KB`,
		2048:    `2.00KB`,
		1048575: `1024.00KB`,
		1048576: `1.00MB`,
	}

	for in, out := range testvalues {
		if v, err := ToByteString(in, `%.2f`); err != nil || v != out {
			t.Errorf("Value %v (%T) ToByteString failed: expected '%s', got '%s' (err: %v)", in, in, out, v, err)
		}
	}
}
