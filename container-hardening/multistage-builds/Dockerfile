# ==========================================
# STAGE 1: Build the Go binary
# ==========================================

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to the workspace
COPY main.go .

# Cross-compile for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webserver main.go

# ==========================================
# STAGE 2: The Final Runtime (The Lean Production Image)
# ==========================================
# We switch to Google's Distroless static image (No shell, no package manager, ~2MB)
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy ONLY the compiled binary from the 'builder' stage
COPY --from=builder /app/webserver /webserver

EXPOSE 8080

# Run the binary
CMD ["/webserver"]