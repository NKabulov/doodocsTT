package service

import (
	"doodocs/internal/models"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
)

type ArchiveServiceImpl struct {
	mimeType   MimeTyper
	processors map[string]ArchiveProcessor
	creator    ArchiveCreator
}

type MimeTyper interface {
	GetMimeType(file *multipart.FileHeader) (string, error)
}

type ArchiveProcessor interface {
	Process(header *multipart.FileHeader) (*models.ArchiveInfo, error)
}

type ArchiveCreator interface {
	CreateArchive(files []*multipart.FileHeader) ([]byte, error)
}

func NewArchiveServiceImpl(mimeType MimeTyper, processors map[string]ArchiveProcessor, creator ArchiveCreator) *ArchiveServiceImpl {
	return &ArchiveServiceImpl{
		mimeType:   mimeType,
		processors: processors,
		creator:    creator,
	}
}

type EmailServiceImpl struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	mimeType     MimeTyper
}

func NewEmailServiceImpl(mimeType MimeTyper) (*EmailServiceImpl, error) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPassword == "" {
		return nil, fmt.Errorf("missing required SMTP configuration: check SMTP_HOST, SMTP_PORT, SMTP_USER, and SMTP_PASSWORD")
	}

	return &EmailServiceImpl{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUser:     smtpUser,
		smtpPassword: smtpPassword,
		mimeType:     mimeType,
	}, nil
}

var (
	ErrOnlyOneFileAllowed     = errors.New("only one file is allowed")
	ErrEmptyFile              = errors.New("uploaded file is empty")
	ErrEmptyFileName          = errors.New("uploaded file name is empty")
	ErrNoFilesUploaded        = errors.New("no files uploaded")
	ErrUnsupportedFormat      = errors.New("unsupported file format")
	ErrNotCorrectEmailAddress = errors.New("not correct email address")
)
