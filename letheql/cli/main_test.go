package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	main()
}

func Test_version(t *testing.T) {
	actual := execute("version")
	expected := "lethetool v0.0.1"
	assert.Equal(t, expected, actual)
}
