package main

import (
	"doodocs/internal/handler"
	"doodocs/internal/mime"
	"doodocs/internal/processor"
	"doodocs/internal/service"
	"log"
	"net/http"
)

func main() {
	// Инициализация MimeTyper
	mimeTyper := mime.NewType()

	// Инициализация Processors
	zipProcessor := processor.NewZipProcessor()
	tarProcessor := processor.NewTarProcessor()

	processors := map[string]service.ArchiveProcessor{
		"application/zip": zipProcessor,
		"application/tar": tarProcessor,
	}

	// Инициализация ArchiveCreator (ZipProcessor переиспользуем)
	archiveCreator := processor.NewZipProcessor()

	// Инициализация ArchiveService
	archiveService := service.NewArchiveServiceImpl(mimeTyper, processors, archiveCreator)

	// Инициализация EmailService
	emailService, err := service.NewEmailServiceImpl(mimeTyper)
	if err != nil {
		log.Fatalf("failed to initialize email service: %v", err)
	}

	// Создаем обработчики
	archiveHandler := handler.NewArchiveHandler(archiveService)
	emailHandler := handler.NewEmailHandler(emailService)

	// Настройка маршрутов
	router := handler.SetupRouter(archiveHandler, emailHandler)

	// Запуск HTTP-сервера
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server is running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
