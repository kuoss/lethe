package filesystem

import (
	"fmt"
	"testing"
)

func TestWalk(t *testing.T) {
	d := New(DriverParameters{RootDirectory: "THIS_WILL_OVERRIDE_BY_DRIVER_CODE"})
	infos, err := d.WalkDir("./tmp/log")
	if err != nil {
		t.Fatalf("err from walk")
	}
	for _, info := range infos {
		fmt.Println(info)
	}
}
