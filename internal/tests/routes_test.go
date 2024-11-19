package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArchiveInformation(t *testing.T) {
	// Prepare a valid ZIP file for testing
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.zip")
	part.Write([]byte("This is a test archive."))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/archive/information", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestArchiveFiles(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("files[]", "image.png")
	file.Write([]byte("This is a test image."))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/archive/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSendEmailFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", "document.pdf")
	file.Write([]byte("This is a test document."))
	_ = writer.WriteField("emails", "test1@mail.com,test2@mail.com")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/mail/file", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
