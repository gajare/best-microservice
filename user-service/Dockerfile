# Build stage
FROM golang:1.24.2 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./main.go

# Run stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/server /app/server

EXPOSE 50051
CMD ["/app/server"]