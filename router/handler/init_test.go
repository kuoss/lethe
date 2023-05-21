package handler

import "github.com/kuoss/lethe/storage/fileservice"

var (
	handler1 = New(&fileservice.FileService{})
)
