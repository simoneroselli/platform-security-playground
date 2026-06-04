# ==========================================
# STAGE 1: Build the Go binary
# ==========================================

FROM: golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to the workspace
COPY main.go .

# Cross-compile for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webserver main.go

