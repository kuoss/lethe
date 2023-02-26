package testutil

import (
	"log"
	"os"
	"testing"

	"github.com/kuoss/lethe/config"
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

	logDirectory := config.GetLogRoot()
	if _, err := os.Stat(logDirectory); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("logDirectory [%s] not exists: %s", logDirectory, err)
		}
		log.Fatalf("Cannot stat logDirectory [%s]: %s", logDirectory, err)
	}
}
