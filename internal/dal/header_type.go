package dal

import (
	"errors"
	"mime/multipart"
)

const (
	contentTypeHeader = "Content-Type"
)

var (
	ErrEmptyContentTypeHeader = errors.New("content type header is empty")
)

type HeaderMimeType struct{}

func (h *HeaderMimeType) GetMimeType(fh *multipart.FileHeader) (string, error) {
	mimeType := fh.Header.Get(contentTypeHeader)
	if mimeType == "" {
		return "", ErrEmptyContentTypeHeader
	}
	return mimeType, nil
}
