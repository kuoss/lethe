package main

import (
	"bytes"
	"strings"

	"github.com/kuoss/lethe/cli/cmd"
	cliutil "github.com/kuoss/lethe/cli/util"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/logs"
	testutil "github.com/kuoss/lethe/testutil"
)

var rotator *logs.Rotator

func init() {
	testutil.Init()
	testutil.SetTestLogFiles()
	rotator = logs.NewRotator()
	config.SetWriter(cliutil.GetWriter())
	// time.Sleep(500 * time.Millisecond)
}

func execute(args ...string) string {
	buf := new(bytes.Buffer)
	cmd := cmd.GetRootCmd()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	cmd.Execute()
	result := strings.TrimSpace(buf.String() + cliutil.GetString())
	// time.Sleep(500 * time.Millisecond)
	return result
}
