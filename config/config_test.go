package config

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {

	fmt.Println(filepath.Join(".", "etc"))
	fmt.Println(filepath.Join("..", "etc"))
}
