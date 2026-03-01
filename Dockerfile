FROM golang:1.25.1 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/myapp .
COPY config/ ./config/
EXPOSE 8080
CMD ["./myapp"]