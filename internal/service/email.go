package service

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/smtp"
	"strings"
)

var ValidEmailFileMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/pdf": true,
}

func (s *EmailServiceImpl) SendEmailWithAttachment(header *multipart.FileHeader, emailsStr string) error {
	mimeType, err := s.mimeType.GetMimeType(header)
	if err != nil {
		slog.Debug("Service Email in GetArchiveInfo")
		return fmt.Errorf("failed to detect archive type: %w", err)
	}
	if !ValidEmailFileMimeTypes[mimeType] {
		slog.Debug("Service Email in GetArchiveInfo")
		return ErrUnsupportedFormat
	}

	emails := strings.Split(emailsStr, ",")
	validEmails := s.validateEmails(emails)
	if len(validEmails) == 0 {
		return fmt.Errorf("no valid email addresses provided")
	}

	file, err := header.Open()
	if err != nil {
		slog.Debug("Service Email in SendEmailWithAttachment: failed to open file")
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		slog.Debug("Service Email in SendEmailWithAttachment: failed to read file data")
		return fmt.Errorf("failed to read file data: %w", err)
	}

	for _, email := range validEmails {
		if err := s.sendEmail(email, header.Filename, mimeType, fileData); err != nil {
			slog.Error("Failed to send email", "email", email, "error", err)
			return fmt.Errorf("failed to send email to %s: %w", email, err)
		}
	}

	slog.Info("All emails sent successfully", "emails", validEmails)
	return nil
}

func (s *EmailServiceImpl) validateEmails(emails []string) []string {
	var validEmails []string
	for _, email := range emails {
		email = strings.TrimSpace(email)
		if strings.Contains(email, "@") && strings.Contains(email, ".") {
			validEmails = append(validEmails, email)
		} else {
			slog.Warn("Invalid email address skipped", "email", email)
		}
	}
	return validEmails
}

func (s *EmailServiceImpl) sendEmail(to, filename, mimeType string, fileData []byte) error {
	from := s.smtpUser
	password := s.smtpPassword
	host := s.smtpHost
	port := s.smtpPort

	// Кодируем файл в Base64
	encodedFile := base64.StdEncoding.EncodeToString(fileData)

	// Создаем заголовки письма
	subject := "Subject: File Attachment\n"
	mime := fmt.Sprintf("MIME-Version: 1.0\nContent-Type: multipart/mixed; boundary=boundary42\n\n")
	headers := subject + mime

	// Тело сообщения
	body := "--boundary42\n" +
		"Content-Type: text/plain; charset=utf-8\n\n" +
		"Here is the file you requested.\n\n"

	// Вложение
	attachment := fmt.Sprintf("--boundary42\nContent-Type: %s\nContent-Disposition: attachment; filename=%q\nContent-Transfer-Encoding: base64\n\n%s\n",
		mimeType, filename, encodedFile)

	// Конец письма
	closing := "--boundary42--"

	// Объединяем все части
	message := headers + body + attachment + closing

	// Настраиваем SMTP-сервер
	auth := smtp.PlainAuth("", from, password, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	// Отправляем письмо
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	slog.Info("Email sent successfully", "to", to, "filename", filename)
	return nil
}
