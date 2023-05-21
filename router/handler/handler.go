package handler

import (
	"github.com/kuoss/lethe/storage/fileservice"
)

type Handler struct {
	fileService *fileservice.FileService
}

func New(fileService *fileservice.FileService) *Handler {
	return &Handler{fileService}
}
