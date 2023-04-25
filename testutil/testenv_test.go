package testutil

import (
	"os"
	"testing"

	"github.com/kuoss/lethe/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	Init()
}

func Test_Init(t *testing.T) {
	testMode := os.Getenv("TEST_MODE")
	if testMode != "1" {
		t.Fatalf("TEST_MODE=[%s], not 1", testMode)
	}
}

func Test_SetTestLogFiles(t *testing.T) {
	SetTestLogFiles()

	logDirectory := config.GetLogDataPath()
	assert.DirExists(t, logDirectory)
}
