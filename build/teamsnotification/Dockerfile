FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates tzdata build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata curl wget
WORKDIR /root/

ENV TZ=Asia/Taipei

COPY --from=builder /app/server ./server
COPY --from=builder /app/configs ./configs

EXPOSE 8080 9090
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]


