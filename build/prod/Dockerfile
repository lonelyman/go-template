# Production Dockerfile - Multi-stage build for optimized binary
# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create a non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application with production optimizations
# คำอธิบาย build flags:
# CGO_ENABLED=0     : ปิดการใช้ CGO (C bindings) ทำให้ได้ static binary
# GOOS=linux        : Target OS เป็น Linux (สำหรับ container)
# GOARCH=amd64      : Target architecture เป็น x86_64
# -ldflags='-w -s'  : ลบ debug info (-w) และ symbol table (-s) ลดขนาด binary
# -extldflags "-static" : Static linking ไม่ต้องพึ่ง shared libraries
# -a                : Force rebuild ทุก package
# -installsuffix cgo : ป้องกัน conflict กับ CGO builds
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
   -ldflags='-w -s -extldflags "-static"' \
   -a -installsuffix cgo \
   -o main cmd/api/main.go

# Stage 2: Final production stage using scratch image
# scratch = เป็น empty image ขนาดเล็กที่สุด (0 MB)
FROM scratch

# Import SSL certificates from builder (ต้องมีสำหรับ HTTPS calls)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import passwd file (สำหรับ non-root user)
COPY --from=builder /etc/passwd /etc/passwd

# Copy the static binary from builder stage
COPY --from=builder /app/main /app/main

# Copy assets (configs removed - using environment variables instead)
COPY --from=builder /app/assets /app/assets

# Use non-root user for security
USER appuser

# Expose port
EXPOSE 9999

# Set the entrypoint
ENTRYPOINT ["/app/main"]
