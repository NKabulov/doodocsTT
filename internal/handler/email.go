package handler

import (
	"doodocs/internal/utils/response"
	"log/slog"
	"net/http"
	"strings"
)

func (h *EmailHandler) SendEmailFile(w http.ResponseWriter, r *http.Request) {
	// limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	// Parse form
	err := r.ParseMultipartForm(maxMemorySize)
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			response.WithError(w, http.StatusRequestEntityTooLarge, err)
			slog.Debug("Handler Email in SendEmailFile")
			return
		}
		response.WithError(w, http.StatusBadRequest, err)
		slog.Debug("Handler Email in SendEmailFile")
		return
	}
	defer r.MultipartForm.RemoveAll()

	files := r.MultipartForm.File["file"]
	if len(files) != 1 {
		response.WithError(w, http.StatusBadRequest, ErrOnlyOneFileAllowed)
		slog.Debug("Handler Email in SendEmailFile")
		return
	}

	header := files[0]
	if header.Size <= 0 {
		response.WithError(w, http.StatusNotAcceptable, ErrEmptyFile)
		slog.Debug("Handler Email in SendEmailFile")
		return
	}

	if header.Filename == "" {
		response.WithError(w, http.StatusNotAcceptable, ErrEmptyFileName)
		slog.Debug("Handler Email in SendEmailFile")
		return
	}

	emailStr := r.FormValue("emails")
	if emailStr == "" {
		response.WithError(w, http.StatusBadRequest, ErrNoEmails)
		slog.Debug("Handler Email in SendEmailFile")
		return
	}

	if err := h.service.SendEmailWithAttachment(header, emailStr); err != nil {
		response.WithError(w, http.StatusInternalServerError, err)
		slog.Debug("Handler Email in SendEmailFile")
		return
	}

	w.WriteHeader(http.StatusOK)

}
