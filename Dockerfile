# Build Stage
FROM golang:1.21 AS build

WORKDIR /app
COPY . .
RUN go build -o main main.go

# 実行権限を設定
RUN chmod +x ./main
# EXPOSE 8000

CMD ["./main"]