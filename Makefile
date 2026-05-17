.PHONY: all build build-backend build-frontend test test-backend test-frontend lint clean dev docker docker-up docker-down

all: build

build: build-backend

build-backend: build-frontend
	rm -rf backend/static/dist
	cp -r frontend/dist backend/static/dist
	cd backend && CGO_ENABLED=1 go build -o posixfy-cloud .

build-frontend:
	cd frontend && npm ci && npm run build

test: test-backend test-frontend

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npx vue-tsc --noEmit

lint:
	cd backend && go fmt ./...
	cd frontend && npx vue-tsc --noEmit

docker:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

dev-backend: build-frontend
	rm -rf backend/static/dist
	cp -r frontend/dist backend/static/dist
	cd backend && API_KEY=dev-key JWT_SECRET=dev-jwt-secret-32chars-minimum go run .

dev-frontend:
	cd frontend && npm run dev

clean:
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf backend/static/dist
	rm -f backend/posixfy-cloud
