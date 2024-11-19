package models

type ArchiveInfo struct {
	Filename    string     `json:"filename"`
	ArchiveSize float64    `json:"archive_size"`
	TotalSize   float64    `json:"total_size"`
	TotalFiles  float64    `json:"total_files"` // float64 by task (not int)
	Files       []FileInfo `json:"files"`
}

type FileInfo struct {
	FilePath string  `json:"file_path"`
	Size     float64 `json:"size"`
	MimeType string  `json:"mimetype"`
}

func NewArchiveInfo(fileName string, archiveSize, totalSize, totalFiles float64, files []FileInfo) *ArchiveInfo {

	return &ArchiveInfo{
		Filename:    fileName,
		ArchiveSize: archiveSize,
		TotalSize:   totalSize,
		TotalFiles:  totalFiles,
		Files:       files,
	}
}

func NewFileInfo(filePath string, size float64, mimeType string) *FileInfo {
	return &FileInfo{
		FilePath: filePath,
		Size:     size,
		MimeType: mimeType,
	}
}
