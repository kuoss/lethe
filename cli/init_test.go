package main

import (
	"bytes"
	"strings"

	"github.com/kuoss/lethe/cli/cmd"
	cliutil "github.com/kuoss/lethe/cli/util"
	"github.com/kuoss/lethe/config"
	testutil "github.com/kuoss/lethe/testutil"
)

func init() {
	config.LoadConfig()
	config.GetConfig().Set("retention.time", "3h")
	config.GetConfig().Set("retention.size", "10m")
	config.SetNow(testutil.GetNow())
	config.SetLogRoot("/tmp/log")
	config.SetWriter(cliutil.GetWriter())
}

func execute(args ...string) string {
	buf := new(bytes.Buffer)
	cmd := cmd.GetRootCmd()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	cmd.Execute()
	return strings.TrimSpace(buf.String() + cliutil.GetString())
}
