package driver

import "fmt"

type PathNotFoundError struct {
	Path string
	Err  error
}

func (err PathNotFoundError) Error() string {
	return fmt.Sprintf("Path not found: %s, err: %s", err.Path, err.Err)
}
