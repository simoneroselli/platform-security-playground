# ==========================================
# STAGE 1: Build the Go binary
# ==========================================
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o webserver main.go

# ==========================================
# STAGE 2: Secure Runtime with Non-Root User
# ==========================================
# We use the specific 'nonroot' variant of the static image
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /app/webserver /webserver

# EXPLICITLY SWITCH TO NON-ROOT USER
# Distroless provides a user named 'nonroot' with UID 65532
USER nonroot:nonroot

EXPOSE 8080

CMD ["/webserver"]