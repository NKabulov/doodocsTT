package handler

import (
	"doodocs/internal/utils/response"
	"log/slog"
	"net/http"
	"strings"
)

func (h *ArchiveHandler) GetArchiveInfo(w http.ResponseWriter, r *http.Request) {
	// limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	// Parse form
	err := r.ParseMultipartForm(maxMemorySize)
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			response.WithError(w, http.StatusRequestEntityTooLarge, err)
			slog.Debug("Handler Archive in GetArchiveInfo")
			return
		}
		response.WithError(w, http.StatusBadRequest, err)
		slog.Debug("Handler Archive in GetArchiveInfo")
		return
	}
	defer r.MultipartForm.RemoveAll()

	// Validate uploaded file
	files := r.MultipartForm.File["file"]
	if len(files) != 1 {
		response.WithError(w, http.StatusBadRequest, ErrOnlyOneFileAllowed)
		slog.Debug("Handler Archive in GetArchiveInfo")
		return
	}

	header := files[0]
	if header.Size <= 0 {
		response.WithError(w, http.StatusNotAcceptable, ErrEmptyFile)
		slog.Debug("Handler Archive in GetArchiveInfo")
		return
	}

	if header.Filename == "" {
		response.WithError(w, http.StatusNotAcceptable, ErrEmptyFileName)
		slog.Debug("Handler Archive in GetArchiveInfo")
		return
	}

	data, err := h.archiveService.GetArchiveInfo(header)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err)
		slog.Debug("Handler Archive in GetArchiveInfo")
		return
	}

	if err := response.JSON(w, data); err != nil {
		response.WithError(w, http.StatusInternalServerError, err)
		slog.Debug("Handler Archive in GetArchiveInfo")
		return
	}

}
