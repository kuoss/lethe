package main

import (
	"testing"

	testutil "github.com/kuoss/lethe/testutil"
)

func Test_version(t *testing.T) {
	got := execute("version")
	want := `"lethetool v0.0.1"`
	testutil.CheckEqualJSON(t, got, want)
}
