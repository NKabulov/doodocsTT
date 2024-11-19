package service

import (
	"doodocs/internal/models"
	"fmt"
	"log/slog"
	"mime/multipart"
)

func (s *ArchiveServiceImpl) GetArchiveInfo(fileHeader *multipart.FileHeader) (*models.ArchiveInfo, error) {
	mimeType, err := s.mimeType.GetMimeType(fileHeader)
	if err != nil {
		slog.Debug("Service Archive in GetArchiveInfo")
		return nil, fmt.Errorf("failed to detect archive type: %w", err)
	}
	fmt.Printf("archive type: %s\n", mimeType)

	processor, ok := s.processors[mimeType]
	if !ok {
		return nil, ErrUnsupportedFormat
	}

	result, err := processor.Process(fileHeader)
	if err != nil {
		slog.Debug("Service Archive in GetArchiveInfo")
		return nil, err
	}

	return result, nil
}
