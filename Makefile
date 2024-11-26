.PHONY: all build run down

PREFIX=docker compose --env-file .env -f deployments/docker-compose.yml

all: build run

build:
	@${PREFIX} build

run:
	@${PREFIX} up -d

down:
	@${PREFIX} down

logs:
	@${PREFIX} logs ${AT}