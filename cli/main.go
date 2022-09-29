package main

import (
	"github.com/kuoss/lethe/cli/cmd"
	cliutil "github.com/kuoss/lethe/cli/util"
	"github.com/kuoss/lethe/config"
)

func main() {
	config.LoadConfig()
	config.SetWriter(cliutil.GetWriter())
	cmd.Execute()
}
