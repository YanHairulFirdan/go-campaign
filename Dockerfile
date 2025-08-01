# ---------- Stage 1: Builder ----------
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install git, curl, make (make optional untuk Air)
RUN apk add --no-cache git curl make

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build main binary
RUN go build -o main .

# Optionally, download migrate CLI binary
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz \
    | tar xvz -C /usr/local/bin

# ---------- Stage 2: Development ----------
FROM golang:1.24.4-alpine AS dev

WORKDIR /app

# Install Air and other tools
RUN apk add --no-cache git curl make && \
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz \
    | tar xvz -C /usr/local/bin

# Copy entire project (for live reload)
COPY . .

# CMD will be handled by Air
CMD ["air"]

# ---------- Stage 3: Production ----------
FROM alpine:latest AS prod

WORKDIR /app

# Copy only binary from builder
COPY --from=builder /app/main .

# Run app
CMD ["./main"]
