package util

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
)

func NewTableWriter(buf *bytes.Buffer) *tablewriter.Table {
	table := tablewriter.NewWriter(buf)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetNoWhiteSpace(true)
	table.SetTablePadding("")
	table.SetColumnSeparator("")
	table.SetTablePadding("   ")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	return table
}
