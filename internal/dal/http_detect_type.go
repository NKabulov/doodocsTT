package dal

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type HttpDetectMimeType struct{}

func (h *HttpDetectMimeType) GetMimeType(fh *multipart.FileHeader) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	headerOfFile, err := io.ReadAll(io.LimitReader(file, maxHeaderSize))
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	return http.DetectContentType(headerOfFile), nil
}
