package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWalk(t *testing.T) {
	d := New(DriverParameters{RootDirectory: "THIS_WILL_OVERRIDE_BY_DRIVER_CODE"})
	userHomeDir, _ := os.UserHomeDir()
	//filepath.
	infos, err := d.WalkDir(filepath.Join(userHomeDir, "tmp", "log"))
	if err != nil {
		t.Fatalf("err from walk")
	}
	for _, info := range infos {
		fmt.Println(info)
	}
}
