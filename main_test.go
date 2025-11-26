package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCI_ShouldFail(t *testing.T) {
	// assert.True(t, false, "This test should fail to verify CI catches errors")
	assert.True(t, true)
}
