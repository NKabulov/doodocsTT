package handler

import (
	"doodocs/internal/utils/response"
	"log/slog"
	"net/http"
	"strings"
)

func (h *ArchiveHandler) CreateZipArchive(w http.ResponseWriter, r *http.Request) {

	// limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	// Parse form
	err := r.ParseMultipartForm(maxMemorySize)
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			response.WithError(w, http.StatusRequestEntityTooLarge, err)
			slog.Debug("Handler ZIPArchive in CreateZipArchive")
			return
		}
		response.WithError(w, http.StatusBadRequest, err)
		slog.Debug("Handler ZIPArchive in CreateZipArchive")
		return
	}
	defer r.MultipartForm.RemoveAll()

	files := r.MultipartForm.File["files[]"]
	if len(files) < 1 {
		response.WithError(w, http.StatusBadRequest, ErrNoFilesUploaded)
		slog.Debug("Handler ZIPArchive in CreateZipArchive")
		return
	}

	for _, header := range files {
		if header.Size <= 0 {
			response.WithError(w, http.StatusNotAcceptable, ErrEmptyFile)
			slog.Debug("Handler ZIPArchive in CreateZipArchive")
			return
		}

		if header.Filename == "" {
			response.WithError(w, http.StatusNotAcceptable, ErrEmptyFileName)
			slog.Debug("Handler ZIPArchive in CreateZipArchive")
			return
		}
	}

	newZip, err := h.archiveService.CreateZipArchive(files)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err)
		slog.Debug("Handler ZIPArchive in CreateZipArchive")
		return
	}

	if err := response.ZIP(w, newZip); err != nil {
		response.WithError(w, http.StatusInternalServerError, err)
		slog.Debug("Handler ZIPArchive in CreateZipArchive")
		return
	}

}
