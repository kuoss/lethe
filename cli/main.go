package main

import (
	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/cli/cmd"
	cliutil "github.com/kuoss/lethe/cli/util"
	"github.com/kuoss/lethe/config"
)

func main() {
	var err error
	err = config.LoadConfig()
	if err != nil {
		logger.Fatalf("error on LoadConfig: %s", err)
	}
	config.SetWriter(cliutil.GetWriter())
	err = cmd.Execute()
	if err != nil {
		logger.Fatalf("error on Execute: %s", err)
	}
}
