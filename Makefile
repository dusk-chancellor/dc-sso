CONFIG_PATH="./configs/local.yml"
JWT_SECRET="randomencryptedhash"

local-run:
	CONFIG_PATH=${CONFIG_PATH} JWT_SECRET=${JWT_SECRET} go run cmd/sso/main.go

local-build:
	CONFIG_PATH=${CONFIG_PATH} JWT_SECRET=${JWT_SECRET} go build -o bin/sso ./cmd/sso/main.go

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
