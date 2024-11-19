package mime

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"
)

type Type struct{}

func NewType() *Type {
	return &Type{}
}

const maxHeaderSize = 512

func (m *Type) GetMimeType(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		slog.Debug("Service Archive in GetArchiveInfo")
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	headerOfFile, err := io.ReadAll(io.LimitReader(file, maxHeaderSize))
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}
	file.Seek(0, io.SeekStart)

	if len(headerOfFile) < 4 {
		return "", fmt.Errorf("file header is too short to determine format")
	}

	mimeType := http.DetectContentType(headerOfFile)

	if mimeType == "application/zip" || (mimeType == "application/octet-stream" && m.isZip(headerOfFile)) {
		if m.isOpenXML(file, fileHeader.Size) {
			return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", nil
		}
		return "application/zip", nil
	}
	if mimeType == "application/octet-stream" && m.isTar(file) {
		return "application/tar", nil
	}
	if mimeType == "text/xml; charset=utf-8" && m.isXML(headerOfFile) {
		return "application/xml", nil
	}

	return mimeType, nil
}

func (m *Type) isZip(headerOfFile []byte) bool {
	return bytes.HasPrefix(headerOfFile, []byte{0x50, 0x4B, 0x03, 0x04})
}

func (m *Type) isTar(file multipart.File) bool {
	tarReader := tar.NewReader(file)

	_, err := tarReader.Next()
	file.Seek(0, io.SeekStart)

	return err == nil
}

func (m *Type) isXML(data []byte) bool {
	trimmed := bytes.TrimSpace(data)
	return bytes.HasPrefix(trimmed, []byte("<?xml")) || bytes.Contains(trimmed, []byte("<?xml"))
}

func (m *Type) isOpenXML(file multipart.File, fileSize int64) bool {
	zipReader, err := zip.NewReader(file, fileSize)
	if err != nil {
		return false
	}

	for _, f := range zipReader.File {
		if f.Name == "[Content_Types].xml" || strings.HasPrefix(f.Name, "word/") ||
			strings.HasPrefix(f.Name, "xl/") || strings.HasPrefix(f.Name, "ppt/") {
			return true
		}
	}
	return false
}
