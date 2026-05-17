# Posixfy Cloud

A web-based file management UI and authentication layer for [Posixfy Bridge](https://github.com/posixfy/posixfy-bridge). Users authenticate via the web interface, and all file operations are executed with their mapped Unix identity through the Bridge service.

> **Posixfy Cloud** 是一个面向 [Posixfy Bridge](https://github.com/posixfy/posixfy-bridge) 的 Web 文件管理 UI 和认证层。用户通过 Web 界面登录后，所有文件操作均通过 Bridge 服务以其映射的 Unix 身份执行。

## Architecture

```
Browser ───▶ Posixfy Cloud (Go + Vue, :8080) ───▶ Posixfy Bridge (Rust, :3000) ───▶ Filesystem
              JWT Auth │ User Mgmt                 setfsuid/setfsgid                  UID/GID
                       │ SQLite                    API Key auth                       enforced
                       │ Proxy + headers                                              by kernel
```

| Component | Tech | Role |
|-----------|------|------|
| **Backend** | Go, Gin, SQLite | JWT auth, user management, API proxy |
| **Frontend** | Vue 3, TypeScript, Element Plus | SPA for file browsing and admin |
| **Bridge** | [posixfy-bridge](https://github.com/posixfy/posixfy-bridge) | Low-level file operations (required dependency) |

## Features

- JWT authentication with role-based access (admin / user)
- File operations: browse, upload, download, delete, create directory
- Real-time file change notifications via SSE
- Admin panel for user management
- Works with local mounts and NFS
- SQLite for persistent storage

## Quick Start

### Prerequisites

Posixfy Cloud requires a running [Posixfy Bridge](https://github.com/posixfy/posixfy-bridge) instance. Start the Bridge first:

```bash
# In a separate terminal, start posixfy-bridge
docker run -d --name posixfy-bridge -p 3000:3000 \
  -e API_KEY=dev-key \
  -e MOUNT_POINTS=data:/data \
  ghcr.io/posixfy/posixfy-bridge:latest
```

### Using Docker

```bash
docker pull ghcr.io/posixfy/posixfy-cloud:latest

docker run -d \
  --name posixfy-cloud \
  -p 8080:8080 \
  -e API_KEY=dev-key \
  -e JWT_SECRET=your-jwt-secret-32chars-minimum \
  -e FILEBRIDGE_URL=http://host.docker.internal:3000 \
  -e ADMIN_INIT_PASSWORD=admin123 \
  -v cloud-data:/data \
  ghcr.io/posixfy/posixfy-cloud:latest
```

### Using Docker Compose

```yaml
services:
  bridge:
    image: ghcr.io/posixfy/posixfy-bridge:latest
    environment:
      - API_KEY=change-me
      - MOUNT_POINTS=data:/data
    volumes:
      - /path/to/data:/data

  cloud:
    image: ghcr.io/posixfy/posixfy-cloud:latest
    ports:
      - "8080:8080"
    environment:
      - API_KEY=change-me
      - JWT_SECRET=change-me-to-a-random-jwt-secret
      - FILEBRIDGE_URL=http://bridge:3000
      - ADMIN_INIT_PASSWORD=admin123
    depends_on:
      - bridge
```

### Manual Setup

```bash
make build

API_KEY=dev-key \
JWT_SECRET=dev-jwt-secret-32chars-minimum \
FILEBRIDGE_URL=http://127.0.0.1:3000 \
LISTEN_ADDR=0.0.0.0:8080 \
ADMIN_INIT_PASSWORD=admin123 \
./backend/posixfy-cloud
```

## Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `API_KEY` | Yes | — | Must match the Bridge service's API key |
| `JWT_SECRET` | Yes | — | Secret for signing JWT tokens (min 32 chars) |
| `FILEBRIDGE_URL` | Yes | — | URL of the running Posixfy Bridge instance |
| `LISTEN_ADDR` | No | `0.0.0.0:8080` | Bind address |
| `DB_PATH` | No | `./posixfy-cloud.db` | SQLite database path |
| `ADMIN_INIT_PASSWORD` | Yes | — | Password for the initial admin user |
| `CORS_ORIGINS` | No | `*` | Allowed CORS origins (comma-separated) |

## API

### Authentication

```
POST /api/auth/login    # Body: { "username", "password" } → { "token" }
GET  /api/auth/me       # Header: Authorization: Bearer <token>
```

### File Operations

All file operations are proxied to Posixfy Bridge with identity headers injected automatically.

```
GET    /api/fs/mounts
GET    /api/fs/list?mount=&path=&page=&limit=
GET    /api/fs/file?mount=&path=
POST   /api/fs/upload?mount=&path=
DELETE /api/fs/delete?mount=&path=
POST   /api/fs/mkdir?mount=&path=
GET    /api/fs/watch?mount=&path=    # SSE stream
```

### Admin

```
GET    /api/admin/users
POST   /api/admin/users
PUT    /api/admin/users/:id
DELETE /api/admin/users/:id
```

### Health

```
GET /health
```

## Project Structure

```
posixfy-cloud/
├── backend/
│   ├── main.go                 # Entry point, server setup
│   ├── go.mod
│   ├── router/router.go        # Route definitions
│   ├── models/user.go          # User model
│   ├── service/
│   │   ├── fs_client.go        # Bridge HTTP client
│   │   └── user_service.go     # User management
│   └── static/embed.go         # Embedded frontend (dist/)
├── frontend/
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── src/                    # Vue 3 SPA
├── Dockerfile
├── Makefile
└── VERSION
```

## Development

```bash
make dev-backend   # Terminal 1: Go backend
make dev-frontend  # Terminal 2: Vue dev server (hot reload)
make test          # Run all tests
make lint          # Go fmt + Vue type check
```

## Security

- Use strong, randomly generated `JWT_SECRET` and `API_KEY` in production
- Remove `ADMIN_INIT_PASSWORD` after initial admin setup
- Configure `CORS_ORIGINS` to specific origins (do not use `*` in production)
- Run behind HTTPS (use nginx, Caddy, or similar reverse proxy)
- The Bridge service should not be exposed directly to the internet

## License

Licensed under the Apache License, Version 2.0.
