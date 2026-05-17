# Stage 1: Build frontend
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go backend with embedded frontend
FROM golang:1.26-bookworm AS backend-builder

WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

# Copy built frontend into the embed directory
RUN rm -rf static/dist
COPY --from=frontend-builder /app/frontend/dist ./static/dist/

# Build with CGO enabled (required for go-sqlite3)
RUN CGO_ENABLED=1 go build -o posixfy-cloud .

# Stage 3: Runtime
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

COPY --from=backend-builder /app/backend/posixfy-cloud /usr/local/bin/posixfy-cloud

RUN useradd -r -s /bin/false -m filebridge
USER filebridge
WORKDIR /home/filebridge

EXPOSE 8080

ENV LISTEN_ADDR=0.0.0.0:8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["posixfy-cloud"]
