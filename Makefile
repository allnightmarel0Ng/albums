.PHONY: all build run down

PREFIX=docker compose --env-file .env -f deployments/docker-compose.yml

all: build run

build:
	@${PREFIX} build

run:
	@${PREFIX} up -d

down:
	@${PREFIX} down --volumes 

logs:
	@${PREFIX} logs ${AT}

ps:
	@${PREFIX} ps -a

exec:
	@${PREFIX} exec ${AT} ${CMD}