# Docker Setup Guide

โปรเจคนี้มี Docker configuration 2 แบบ สำหรับ Development และ Production

## Development Mode

### Features

-  ใช้ CompileDaemon สำหรับ auto-reload เมื่อมีการเปลี่ยนแปลงโค้ด
-  Mount source code เข้าไปใน container
-  Debug mode เปิดใช้งาน
-  Interactive mode สำหรับการ debug

### การใช้งาน Development

```bash
# Build และ run development container
make docker-dev-up

# หรือใช้ docker-compose โดยตรง
docker-compose -f docker-compose.dev.yml up --build

# ดู logs
make docker-dev-logs

# หยุด development container
make docker-dev-down
```

### Files สำหรับ Development

-  `build/package/Dockerfile.dev` - Development Dockerfile พร้อม CompileDaemon
-  `docker-compose.dev.yml` - Development compose file

## Production Mode

### Features

-  Multi-stage build สำหรับ optimized image
-  ใช้ scratch image สำหรับ final stage
-  Binary ขนาดเล็กและปลอดภัย
-  Non-root user
-  Static linking

### การใช้งาน Production

```bash
# Build และ run production container
make docker-prod-up

# หรือใช้ docker-compose โดยตรง
docker-compose up --build -d

# ดู logs
make docker-prod-logs

# หยุด production container
make docker-prod-down
```

### Files สำหรับ Production

-  `build/package/Dockerfile` - Production Dockerfile (multi-stage build)
-  `docker-compose.yml` - Production compose file

## Port Configuration

| Mode        | Service | Container Port | Host Port |
| ----------- | ------- | -------------- | --------- |
| Development | App     | 9999           | 9999      |
| Development | Redis   | 6379           | 6380      |
| Production  | App     | 9999           | 9090      |
| Production  | Redis   | 6379           | 6379      |

## Environment Variables

การตั้งค่า environment variables ทั้งสอง mode:

```env
DB_HOST=host.docker.internal
DB_PORT=7430
DB_USER=root
DB_PASSWORD=12345678
DB_NAME=postgres
```

-  **Development**: `FIBER_MODE=debug`
-  **Production**: `FIBER_MODE=release`

## Available Make Commands

```bash
# Development
make docker-dev-build     # Build development image
make docker-dev-up        # Start development containers
make docker-dev-down      # Stop development containers
make docker-dev-logs      # View development logs
make docker-rebuild-dev   # Rebuild development containers

# Production
make docker-prod-build    # Build production image
make docker-prod-up       # Start production containers
make docker-prod-down     # Stop production containers
make docker-prod-logs     # View production logs
make docker-rebuild-prod  # Rebuild production containers

# Utility
make docker-clean         # Clean up unused Docker resources
```

## การเชื่อมต่อ Database

โปรเจคใช้ PostgreSQL ที่รันอยู่ในเครื่องแล้ว (port 7430) แทนการใช้ PostgreSQL container

-  **Host**: `host.docker.internal`
-  **Port**: `7430`
-  **User**: `root`
-  **Password**: `12345678`
-  **Database**: `postgres`

## Auto-reload ใน Development

Development mode ใช้ [CompileDaemon](https://github.com/githubnemo/CompileDaemon) สำหรับ auto-reload:

-  ตรวจจับการเปลี่ยนแปลงไฟล์ Go
-  Compile และ restart application อัตโนมัติ
-  แสดง colored output
-  Graceful shutdown

## Directory Structure

```
build/
├── package/
│   ├── Dockerfile      # Production Dockerfile
│   └── Dockerfile.dev  # Development Dockerfile
docker-compose.yml      # Production compose
docker-compose.dev.yml  # Development compose
```

## Tips

1. **Development**: ใช้ `docker-compose.dev.yml` สำหรับการพัฒนา
2. **Production**: ใช้ `docker-compose.yml` สำหรับ production
3. **Logs**: ใช้ `make docker-*-logs` เพื่อดู real-time logs
4. **Clean up**: ใช้ `make docker-clean` เพื่อล้าง unused images และ containers
