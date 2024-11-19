
# Doodocs Backend Challenge

## Описание проекта

### Render.com:
1. https://doodocstt.onrender.com/archive/information
2. https://doodocstt.onrender.com/archive/files
3. https://doodocstt.onrender.com/mail/file
4. 
### Video:
https://www.loom.com/share/14277b818af448558ecfb4012b706e23?sid=d9ff0ce4-e6d9-4553-b24e-8074881b64be
если первая ссылка не сработает: https://drive.google.com/file/d/1TPwhGhPa03jVGH9AAl7e_HKCr4oGXYv4/view?usp=sharing

Этот проект был разработан для выполнения технического задания на стажировку в Doodocs.kz. Проект представляет собой REST API, реализующее следующие функции:

1. **Отображение информации об архиве**:
    - Принимает ZIP-архив через `multipart/form-data`, возвращает подробную информацию о его содержимом.
2. **Формирование ZIP-архива**:
    - Принимает список файлов, объединяет их в архив и отправляет обратно клиенту.
3. **Отправка файла по email**:
    - Принимает файл и список email-адресов, отправляет файл на указанные адреса.

## Особенности проекта

- **Трёхслойная архитектура**:
    - В проекте реализован аналог трёхслойной архитектуры: разделены Presentation (обработчики HTTP-запросов), Business Logic (основная логика обработки) и Data Access (работа с файлами и внешними API).
- **Ограничение ресурсов**:
    - Для защиты от атак и чрезмерного использования памяти сервером введены ограничения:
        - `maxMemorySize = 10 * 1024 * 1024` (10 MB): Максимальный размер для обработки файлов в оперативной памяти.
        - `maxBodySize = 25 * 1024 * 1024` (25 MB): Максимальный общий размер тела запроса.
- **Использование переменных окружения для безопасности**:
    - Для отправки email используются следующие переменные окружения:
        - `SMTP_HOST`: Адрес SMTP-сервера.
        - `SMTP_PORT`: Порт SMTP-сервера.
        - `SMTP_USER`: Имя пользователя для аутентификации.
        - `SMTP_PASSWORD`: Пароль для аутентификации.
    - **Как настроить**:
      Вы можете написать в терминале:
      export SMTP_HOST=smtp.gmail.com
      export SMTP_PORT=587
      export SMTP_USER=your_email@gmail.com
      export SMTP_PASSWORD=your_password
    - **Как проверить**:
      Напишите в терминале env и вы сможете найти эти окружения

## Установка и запуск

1. **Клонирование репозитория**:

    ```bash
    git clone <URL_репозитория>
    cd <название_папки_проекта>
    ```

2. **Запуск приложения**:

    ```bash
    go run cmd/main.go
    ```

Приложение запустится на порту `8080`.

## API Эндпоинты

### 1. Получение информации об архиве

**URL**: `/archive/information`  
**Метод**: `POST`  
**Описание**: Принимает архив, возвращает информацию о его структуре.

**Пример запроса**:

```http
POST /api/archive/information HTTP/1.1
Content-Type: multipart/form-data; boundary=-{boundary}

-{boundary}
Content-Disposition: form-data; name="file"; filename="archive.zip"
Content-Type: application/zip

{Binary data}
-{boundary}--
```

**Пример ответа**:

```json
{
  "filename": "archive.zip",
  "archive_size": 10240,
  "total_size": 20480,
  "total_files": 2,
  "files": [
    {
      "file_path": "document.docx",
      "size": 10240,
      "mimetype": "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
    }
  ]
}
```

---

### 2. Создание архива

**URL**: `/archive/files`  
**Метод**: `POST`  
**Описание**: Принимает список файлов и создает ZIP-архив.

**Пример запроса**:

```http
POST /api/archive/files HTTP/1.1
Content-Type: multipart/form-data; boundary=-{boundary}

-{boundary}
Content-Disposition: form-data; name="files[]"; filename="file1.png"
Content-Type: image/png

{Binary data}
-{boundary}--
```

**Пример ответа**: ZIP-файл.

---

### 3. Отправка файла по email

**URL**: `/mail/file`  
**Метод**: `POST`  
**Описание**: Отправляет файл на указанные email-адреса.

**Пример запроса**:

```http
POST /api/mail/file HTTP/1.1
Content-Type: multipart/form-data; boundary=-{boundary}

-{boundary}
Content-Disposition: form-data; name="file"; filename="document.pdf"
Content-Type: application/pdf

{Binary data}
-{boundary}
Content-Disposition: form-data; name="emails"

example1@mail.com,example2@mail.com
-{boundary}--
```

**Пример ответа**:

```http
HTTP/1.1 200 OK
```

## Основные обработчики

1. **Archive Information Handler**: Обрабатывает запросы для анализа ZIP-архива и Tar-архтва.
2. **Create Archive Handler**: Обрабатывает создание нового ZIP-архива из загружаемых файлов.
3. **Send Email File Handler**: Обрабатывает отправку загружаемых файлов на email.

