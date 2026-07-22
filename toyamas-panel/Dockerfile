# ==============================================================================
# Toyamas Panel Dockerfile
# Target Spec: Extremely Lightweight (< 30 MB RAM footprint)
# Build Type: Multi-stage Pure Go (CGO-Free)
# ==============================================================================

# --- Stage 1: Builder ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git & CA certificates
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Compile lightweight static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o toyamas-panel ./cmd/server

# --- Stage 2: Minimal Runtime ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata procps systemd

WORKDIR /app

# Copy binary & web assets
COPY --from=builder /app/toyamas-panel /app/toyamas-panel
COPY --from=builder /app/web /app/web

# Create data directory for SQLite
RUN mkdir -p /app/data

# Environment Variables
ENV PORT=8080
ENV DB_PATH=/app/data/toyamas.db
ENV ADMIN_USER=admin
ENV ADMIN_PASS=toyamas123

EXPOSE 8080

ENTRYPOINT ["/app/toyamas-panel"]
