package processor

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"doodocs/internal/models"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type ZipProcessor struct{}

func NewZipProcessor() *ZipProcessor {
	return &ZipProcessor{}
}

func (z *ZipProcessor) Process(fileHeader *multipart.FileHeader) (*models.ArchiveInfo, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	zipReader, err := zip.NewReader(file, fileHeader.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to read zip archive: %w", err)
	}

	var files []models.FileInfo
	var totalSize float64

	for _, f := range zipReader.File {
		if strings.HasPrefix(f.Name, "__MACOSX/") || strings.HasPrefix(filepath.Base(f.Name), "._") {
			continue
		}
		if f.FileInfo().IsDir() ||
			strings.HasSuffix(f.Name, ".DS_Store") ||
			strings.HasPrefix(filepath.Base(f.Name), "._") {
			continue
		}
		fileInArchive, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file in archive: %w", err)
		}
		defer fileInArchive.Close()

		headerOfFile, err := io.ReadAll(io.LimitReader(fileInArchive, 512))
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read file header: %w", err)
		}

		mimeType := http.DetectContentType(headerOfFile)

		fileInfo := models.NewFileInfo(
			f.Name,
			float64(f.UncompressedSize64),
			mimeType,
		)
		files = append(files, *fileInfo)

		totalSize += float64(f.UncompressedSize64)
	}

	archiveInfo := models.NewArchiveInfo(
		fileHeader.Filename,
		float64(fileHeader.Size),
		totalSize,
		float64(len(files)),
		files,
	)

	return archiveInfo, nil
}

func (z *ZipProcessor) CreateArchive(files []*multipart.FileHeader) ([]byte, error) {
	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to add file %s to archive: %w", fileHeader.Filename, err)
		}

		_, err = io.Copy(zipFile, file)
		if err != nil {
			return nil, fmt.Errorf("failed to write file %s to archive: %w", fileHeader.Filename, err)
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to finalize ZIP archive: %w", err)
	}

	return buffer.Bytes(), nil
}

type TarProcessor struct{}

func NewTarProcessor() *TarProcessor {
	return &TarProcessor{}
}

func (t *TarProcessor) Process(fileHeader *multipart.FileHeader) (*models.ArchiveInfo, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	var files []models.FileInfo
	var totalSize float64

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar archive: %w", err)
		}
		var mimeType string
		if hdr.Typeflag != tar.TypeDir {
			headerOfFile, err := io.ReadAll(io.LimitReader(tarReader, 512))
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("failed to read file header: %w", err)
			}
			mimeType = http.DetectContentType(headerOfFile)
		} else {
			mimeType = "directory"
		}

		fileInfo := models.NewFileInfo(
			hdr.Name,
			float64(hdr.Size),
			mimeType,
		)
		files = append(files, *fileInfo)

		totalSize += float64(hdr.Size)
	}

	archiveInfo := models.NewArchiveInfo(
		fileHeader.Filename,
		float64(fileHeader.Size),
		totalSize,
		float64(len(files)),
		files,
	)

	return archiveInfo, nil
}
