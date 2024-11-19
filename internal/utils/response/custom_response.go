package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	_, err = w.Write(jsonData)
	return err
}

func WithError(w http.ResponseWriter, code int, errInput error) {
	response := ErrorResponse{
		Code:    code,
		Message: errInput.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error(err.Error())
	}
}

func ZIP(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="archive.zip"`)
	w.WriteHeader(http.StatusOK)

	_, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	return nil
}
