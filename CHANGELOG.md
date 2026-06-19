# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

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
