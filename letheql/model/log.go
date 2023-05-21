package model

import (
	"strings"

	"github.com/kuoss/lethe/letheql/parser"
)

const ValueTypeLog parser.ValueType = "log"

type Log struct {
	Name  string
	Lines []LogLine
}

func (l Log) Type() parser.ValueType { return ValueTypeLog }
func (l Log) String() string {
	entries := make([]string, len(l.Lines))
	for i, line := range l.Lines {
		entries[i] = line.String()
	}
	return strings.Join(entries, "\n")
}

type LogLine interface {
	String() string
}
