FROM golang:1.25.1-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/server

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/docs ./docs

EXPOSE 8080
CMD ["./app"]
