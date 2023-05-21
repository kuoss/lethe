package util

import "bytes"

var writer = new(bytes.Buffer)

func GetWriter() *bytes.Buffer {
	return writer
}

func Clean() {
	writer.Truncate(0)
}

func GetString() string {
	defer Clean()
	return writer.String()
}
