package fileservice

import (
	"bufio"
	"fmt"
)

func (s *FileService) Scanner(subpath string) (*bufio.Scanner, error) {
	rc, err := s.driver.Reader(subpath)
	if err != nil {
		return nil, fmt.Errorf("reader err: %w", err)
	}
	return bufio.NewScanner(rc), nil
}
