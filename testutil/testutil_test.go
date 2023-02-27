package testutil

import (
	"testing"
)

func Test_DirectoryExists(t *testing.T) {
	var dir string

	// exists
	dir = "."
	if !DirectoryExists(dir) {
		t.Fatalf("Directory [%s] not exists.", dir)
	}

	// not exists
	dir = "./not-exists-dir"
	if DirectoryExists(dir) {
		t.Fatalf("Directory [%s] exists.", dir)
	}

	// file
	dir = "./Makefile"
	if DirectoryExists(dir) {
		t.Fatalf("Directory [%s] exists.", dir)
	}
}

func Test_CheckEqualJSON(t *testing.T) {
	CheckEqualJSON(t, "hello", `"hello"`)
}
