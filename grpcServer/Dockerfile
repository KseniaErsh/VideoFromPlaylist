# Сборка grpcServer

FROM golang:1.18

# Создём директорию приложения
RUN mkdir /app

# Копируем код
COPY . /app/

# Устанавливаем рабочую директорию
WORKDIR /app/

# Собираем приложение
RUN go build -o grpcServer-exe ./grpcServer

# Указываем рабочий порт
EXPOSE 50051

CMD ["/app/grpcServer-exe"]
