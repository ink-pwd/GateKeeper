# Создаем образ и билдер, что бы готовый образ изолировать
# от ненужных файлов по типу тестов, sql таблицы readme и тд
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем все
COPY  . .

# Компилируем проект
# CGO отключен для создания полностью статического бинарника
# (не будет искать внешние библиотеки)
# GOOS=linux - указываем целевую ОС
# Указываем путь к main.go для компиляции
RUN CGO_ENABLED=0 GOOS=linux go build -o gatekeeper ./cmd/api-server

# Финальный образ чисто с бинарником
FROM alpine:latest

# Устанавливаем сертификаты, нужны для https
RUN apk --no-cache add ca-certificates

# Копируем созданный бинарник
COPY --from=builder /app/gatekeeper .
# Копируем env
COPY --from=builder /app/.env .

CMD ["/gatekeeper"]
