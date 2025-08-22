# Development Dockerfile with CompileDaemon for auto-reload
FROM golang:alpine AS development

# Install git, ca-certificates and CompileDaemon
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
