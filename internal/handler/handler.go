package handler

import (
	"doodocs/internal/models"
	"errors"
	"mime/multipart"
)

type ArchiveHandler struct {
	archiveService ArchiveService
}

type ArchiveService interface {
	GetArchiveInfo(fileHeader *multipart.FileHeader) (*models.ArchiveInfo, error)
	CreateZipArchive(fileHeaders []*multipart.FileHeader) ([]byte, error)
}

func NewArchiveHandler(archServ ArchiveService) *ArchiveHandler {
	return &ArchiveHandler{
		archiveService: archServ,
	}
}

type EmailService interface {
	SendEmailWithAttachment(header *multipart.FileHeader, emailsStr string) error
}

type EmailHandler struct {
	service EmailService
}

func NewEmailHandler(serv EmailService) *EmailHandler {
	return &EmailHandler{service: serv}
}

// explained in the video
const (
	maxMemorySize = 10 * 1024 * 1024 // 10 MB for RAM (
	maxBodySize   = 25 * 1024 * 1024 // 25 MB. This is a protection against malicious users who might send an excessively large file or multiple files to exhaust the server's resources.
)

var (
	ErrOnlyOneFileAllowed = errors.New("only one file is allowed")
	ErrEmptyFile          = errors.New("uploaded file is empty")
	ErrEmptyFileName      = errors.New("uploaded file name is empty")
	ErrNoFilesUploaded    = errors.New("no files uploaded")
	ErrNoEmails           = errors.New("no emails attached")
)
