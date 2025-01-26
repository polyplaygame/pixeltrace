package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToSliceByte(t *testing.T) {
	s := "hello world!"
	s1 := []byte(s)
	b := StringToSliceByte(s)
	assert.Equal(t, b, s1)
}

func TestSliceByteToString(t *testing.T) {
	s1 := "0123"
	b1 := []byte(s1)
	b2 := StringToSliceByte(s1)
	assert.Equal(t, len(b1), len(b1))
	assert.Equal(t, string(b2), s1)
}
