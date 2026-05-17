# Contributing to Posixfy Cloud

Thank you for your interest in contributing to Posixfy Cloud!

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Run tests: `make test`
6. Run linting: `make lint`
7. Commit and push
8. Open a Pull Request

## Development Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- Make
- A running instance of [Posixfy Bridge](https://github.com/posixfy/posixfy-bridge)

### Quick Start
```bash
make build
make dev-backend   # Terminal 1
make dev-frontend  # Terminal 2
make test
make lint
```

## Code Style

### Go
- `go fmt`
- Write tests in `_test.go` files using table-driven tests

### TypeScript/Vue
- Vue 3 Composition API
- PascalCase component names
- TypeScript for all code

## Pull Request Guidelines

- Focused PRs — one logical change per PR
- Include tests for new behavior
- Update documentation if behavior changes
- Ensure CI passes
- Clear commit messages

## Reporting Issues

Use GitHub Issues. Include reproduction steps and environment details. Report security vulnerabilities privately.

## License

Contributions are licensed under the Apache License, Version 2.0.
