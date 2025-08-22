# Development Dockerfile with CompileDaemon for auto-reload
FROM golang:1.23-alpine AS development

# Install git, ca-certificates และ CompileDaemon
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN go install github.com/githubnemo/CompileDaemon@latest

# Create working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies (without vendor)
RUN go mod download && go clean -modcache

# Copy source code
COPY . .

# Expose port
EXPOSE 9999

# Use CompileDaemon to watch for changes and auto-reload
# Development build ไม่ต้องใช้ optimization flags เพราะต้องการ compile เร็ว
CMD ["CompileDaemon", \
   "-log-prefix=false", \
   "-build=go build -mod=readonly -o /tmp/main cmd/api/main.go", \
   "-command=/tmp/main", \
   "-directory=.", \
   "-exclude-dir=.git", \
   "-exclude-dir=vendor", \
   "-exclude-dir=build", \
   "-exclude-dir=tests", \
   "-color=true", \
   "-graceful-kill=true"]
