package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustConfig(t *testing.T) {
	version := "test-version"
	cfg := mustConfig(version)
	assert.NotEmpty(t, cfg)
	assert.Equal(t, version, cfg.Version)
}

func TestMustFileService(t *testing.T) {
	version := "test-version"
	cfg := mustConfig(version)
	fileService := mustFileService(cfg)
	assert.NotNil(t, fileService, "Expected file service to be non-nil")
}

func TestStartRotator(t *testing.T) {
	version := "test-version"
	cfg := mustConfig(version)
	fileService := mustFileService(cfg)
	assert.NotPanics(t, func() {
		startRotator(cfg, fileService)
	}, "Expected startRotator to not panic")
}
