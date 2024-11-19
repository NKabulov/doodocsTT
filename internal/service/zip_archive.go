package service

import (
	"fmt"
	"log/slog"
	"mime/multipart"
)

var validZipFileMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func (s *ArchiveServiceImpl) CreateZipArchive(files []*multipart.FileHeader) ([]byte, error) {
	if err := s.isValidFiles(files); err != nil {
		return nil, err
	}
	zipData, err := s.creator.CreateArchive(files)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZIP archive: %w", err)
	}

	return zipData, nil
}

func (s *ArchiveServiceImpl) isValidFiles(files []*multipart.FileHeader) error {
	for _, header := range files {
		if header.Filename == "" {
			return ErrEmptyFileName
		}
		if header.Size <= 0 {
			return ErrEmptyFile
		}

		mimeType, err := s.mimeType.GetMimeType(header)
		if err != nil {
			slog.Debug("Service ZIPArchive in CreateZipArchive")
			return fmt.Errorf("failed to detect archive type: %w", err)
		}

		if !validZipFileMimeTypes[mimeType] {
			slog.Debug("Service ZIPArchive in CreateZipArchive")
			return ErrUnsupportedFormat
		}
	}

	return nil

}
