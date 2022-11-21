# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /server cmd/http/main.go
RUN go build -o /migrations cmd/migrations/postgres/main.go

EXPOSE 8081

CMD "/migrations"; "/server"