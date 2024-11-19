package handler

import "net/http"

func SetupRouter(archiveHandler *ArchiveHandler, emailHandler *EmailHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/archive/information", http.HandlerFunc(archiveHandler.GetArchiveInfo))
	mux.Handle("/archive/files", http.HandlerFunc(archiveHandler.CreateZipArchive))

	mux.Handle("/mail/file", http.HandlerFunc(emailHandler.SendEmailFile))

	return mux
}
