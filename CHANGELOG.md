# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added
- Structured JSON logging (`slog`) with a `LOG_LEVEL` environment variable (`debug`/`info`/`warn`/`error`, default `info`). Failures calling the bridge are now logged instead of being silently swallowed, and upstream responses with status `>= 400` are logged with the response body.
- `X-Request-Id` correlation id propagated browser → cloud → bridge for end-to-end request tracing.
- "Show hidden" toggle in the file browser; system/junk files (`.DS_Store`, `._*`, `Thumbs.db`, `desktop.ini`) are hidden by default.
- Failed API calls are now logged to the browser devtools console.

### Changed
- **BREAKING:** Renamed the bridge connection environment variable `FILEBRIDGE_URL` → `POSIXFY_BRIDGE_URL`. Existing deployments must update this variable in their `docker-compose.yml` / Swarm `stack.yml` / secrets; the old name is no longer read and the connection silently falls back to the default `http://127.0.0.1:3000`. The Go config field `FilebridgeURL` was likewise renamed to `BridgeURL`.
- Rebranded the web UI from "FileBridge" to "Posixfy" (browser title, app header, login page); the frontend npm package is now `posixfy-web`.
- The container image now runs as the `posixfy` system user (previously `filebridge`).

### Fixed
- File "Modified" column showed dates near 1970; timestamps are now interpreted as epoch milliseconds (matching the bridge).
- A deleted file could reappear after refresh due to a cached listing; `/api/fs/list` now sends `Cache-Control: no-store`.

## [0.2.0] - 2026-06-19

### Added
- Docker Swarm stack manifest (`stack.yml`) for production deployment

### Fixed
- Encode mount and path query parameters in FS proxy upstream URLs

## [0.1.0] - 2026-05-17

### Added
- Initial release of Posixfy Cloud
- JWT authentication with role-based access (admin / user)
- User management API (CRUD)
- Admin panel with RBAC
- Proxy to Posixfy Bridge with identity header injection
- Rate limiting on login endpoint
- SQLite persistent storage
- Vue 3 SPA with Element Plus UI
- File browser, upload/mkdir dialogs, admin panel
- SSE real-time file change notifications
- Graceful shutdown
- Health check endpoint
