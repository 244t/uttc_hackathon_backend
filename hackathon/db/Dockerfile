# Build Stage
FROM golang:1.21 AS build

WORKDIR /app
COPY main.go .
RUN go build -o server main.go

# 実行権限を設定
RUN chmod +x ./server
EXPOSE 8080

CMD ["./server"]