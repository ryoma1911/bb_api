# Goのイメージを指定
FROM --platform=linux/amd64 golang:1.23-alpine

WORKDIR /code

# go.mod と go.sum を先にコピー（これが超重要！）
COPY go.mod go.sum ./
RUN go mod download

# 全コードをコピー
COPY . .

RUN go build -o main ./cmd

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN apk add --no-cache tzdata
ENV TZ="Asia/Tokyo"
RUN echo $TZ > /etc/timezone

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
