CONFIG_PATH = ./configs/local.yml
JWT_SECRET = randomencryptedhash

go-run:
	go run cmd/sso/main.go

go-build:
	go build -o bin/sso ./cmd/sso/main.go

sso-up:
	docker compose up -d sso

sso-stop:
	docker compose stop sso

db-up:
	docker compose up -d postgres

db-stop:
	docker compose stop postgres

redis-up:
	docker compose up -d redis

redis-stop:
	docker compose stop redis

up:
	docker compose up -d

down:
	docker compose down
