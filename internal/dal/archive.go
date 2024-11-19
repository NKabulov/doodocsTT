package dal

//
//import (
//	"archive/tar"
//	"archive/zip"
//	"bytes"
//	"doodocs/internal/models"
//	"errors"
//	"fmt"
//	"io"
//	"mime/multipart"
//	"net/http"
//)
//
//type ArchiveRepoImpl struct{}
//
//func NewArchiveRepoImpl() *ArchiveRepoImpl {
//	return &ArchiveRepoImpl{}
//}
//
//func (r *ArchiveRepoImpl) DetectArchiveType(file multipart.File) (string, error) {
//	//TODO: optionally add DTECTCOntenttype
//	mimeType := http.DetectContentType(file)
//	if isZip := r.isZip(file); isZip {
//		return "zip", nil
//	}
//
//	if isTar := r.isTar(file); isTar {
//		return "tar", nil
//	}
//
//	return "", errors.New("unsupported archive format")
//}
//
//func (r *ArchiveRepoImpl) isZip(file multipart.File) bool {
//	header := make([]byte, 4)
//	_, err := file.Read(header)
//	if err != nil {
//		return false
//	}
//
//	if err := resetFilePointer(file); err != nil {
//		return false
//	}
//
//	return bytes.HasPrefix(header, []byte{0x50, 0x4B, 0x03, 0x04})
//}
//
//func (r *ArchiveRepoImpl) isTar(file multipart.File) bool {
//	tarReader := tar.NewReader(file)
//
//	_, err := tarReader.Next()
//
//	if err := resetFilePointer(file); err != nil {
//		return false
//	}
//
//	return err == nil
//}
//
//func (r *ArchiveRepoImpl) GetFilesInfo(file multipart.File, archiveType string, size int64) ([]models.FileInfo, float64, error) {
//	switch archiveType {
//	case "zip":
//		return r.processZip(file, size)
//	case "tar":
//		return r.processTar(file)
//	default:
//		return nil, 0, errors.New("unsupported archive format")
//	}
//}
//
//func (r *ArchiveRepoImpl) processZip(file multipart.File, size int64) ([]models.FileInfo, float64, error) {
//	reader, ok := file.(io.ReaderAt)
//	if !ok {
//		return nil, 0, errors.New("file does not support random access")
//	}
//
//	zipReader, err := zip.NewReader(reader, size)
//	if err != nil {
//		return nil, 0, fmt.Errorf("failed to read zip archive: %w", err)
//	}
//
//	var files []models.FileInfo
//	var totalSize float64
//
//	for _, f := range zipReader.File {
//		files = append(files, models.FileInfo{
//			FilePath: f.Name,
//			Size:     float64(f.UncompressedSize64),
//			MimeType: "unknown", // MIME-тип можно дополнительно определить
//		})
//		totalSize += float64(f.UncompressedSize64)
//	}
//
//	return files, totalSize, nil
//}
//
//func (r *ArchiveRepoImpl) processTar(file multipart.File) ([]models.FileInfo, float64, error) {
//	tarReader := tar.NewReader(file)
//
//	var files []models.FileInfo
//	var totalSize float64
//
//	for {
//		hdr, err := tarReader.Next()
//		if err == io.EOF {
//			break // Конец архива
//		}
//		if err != nil {
//			return nil, 0, fmt.Errorf("failed to read tar archive: %w", err)
//		}
//
//		// Пропускаем директории
//		if hdr.Typeflag == tar.TypeDir {
//			continue
//		}
//
//		files = append(files, models.FileInfo{
//			FilePath: hdr.Name,
//			Size:     float64(hdr.Size),
//			MimeType: "unknown", // MIME-тип можно дополнительно определить
//		})
//		totalSize += float64(hdr.Size)
//	}
//
//	return files, totalSize, nil
//}
//
//func resetFilePointer(file multipart.File) error {
//	if seeker, ok := file.(io.Seeker); ok {
//		_, err := seeker.Seek(0, io.SeekStart)
//		return err
//	}
//	return errors.New("file does not support Seek")
//}
