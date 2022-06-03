# Сборка webServer

FROM golang:1.18

# Создём директорию приложения
RUN mkdir /app

# Копируем код
COPY . /app/

# Устанавливаем рабочую директорию
WORKDIR /app/

# Собираем приложение
RUN go build -o webServer-exe ./webServer

# Указываем рабочий порт
EXPOSE 8080

CMD ["/app/webServer-exe"]
