# ---------- Builder ----------
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# ---------- Runtime ----------
FROM alpine:latest

WORKDIR /app

# copy binary
COPY --from=builder /app/server .

# 🔥 copy migrations folder
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]