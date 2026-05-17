# Stage 1 - build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/seed ./cmd/seed

# Stage 2 - final image
FROM alpine:3.19
WORKDIR /app
RUN apk --no-cache add ca-certificates postgresql-client netcat-openbsd

# Install golang-migrate
RUN wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

COPY --from=builder /app/server .
COPY --from=builder /app/seed .
COPY configs/ configs/
COPY internal/platform/database/migrations migrations/
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8083
ENTRYPOINT ["./entrypoint.sh"]
