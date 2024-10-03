package main

import (
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/assert"
)

func TestMustConfig(t *testing.T) {
	version := "test"
	cfg := mustConfig(version)
	assert.NotEmpty(t, cfg)
	assert.Equal(t, version, cfg.Version)
}

func TestMustFileService(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.main.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	cfg := mustConfig("test")
	fileService := mustFileService(cfg)
	assert.NotNil(t, fileService, "Expected file service to be non-nil")
}

func TestStartRotator(t *testing.T) {
	_, cleanup := tester.MustSetupDir(t, map[string]string{
		"@/testdata/etc/lethe.main.yaml": "etc/lethe.yaml",
	})
	defer cleanup()

	version := "test-version"
	cfg := mustConfig(version)
	fileService := mustFileService(cfg)
	assert.NotPanics(t, func() {
		startRotator(cfg, fileService)
	}, "Expected startRotator to not panic")
}
