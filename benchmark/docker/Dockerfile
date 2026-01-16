FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]