package main

import (
	"github.com/kuoss/lethe/cmd/lethetool/cmd"
)

var (
	Version = "development"
)

func main() {
	cmd.Execute(Version)
}
