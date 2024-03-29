.PHONY: bash build deep-clean psql help logs migrate ps restart stop test up db_dump

SERVICE ?= bot

SLACK_TOKEN ?= `grep 'SLACK_TOKEN=' .env | cut -d '=' -f2`
KOIN_BOT_ID ?= `grep 'KOIN_BOT_ID=' .env | cut -d '=' -f2`

default: help

all: build deps up setup_db migrate #: Boot up the project with db dependencies

bash: #: Bash prompt on running container
	docker compose exec $(SERVICE) bash

build: #: Build containers
	touch .env
	docker compose build

deep-clean: stop
	docker system prune --all --volumes --force

logs: #: Tail the service container's logs
	docker compose logs -tf $(SERVICE)

migrate: build up #: Run migrations
	docker compose run --rm --entrypoint /bin/migration bot

ps: #: Show running processes
	docker compose ps

psql: #: psql prompt on db container, password is postgres
	docker compose exec db psql -d akc -h localhost -U postgres -W

restart: #: Restart the service container
	docker compose restart $(SERVICE)

setup: build up
	docker compose run --rm --entrypoint /bin/local-dev-setup bot

stop: #: Stop running containers
	docker compose stop

test: #: Run tests
	docker compose run --rm -e MIX_ENV=test $(SERVICE) mix test

up: build #: Start containers
	docker compose up -d

down: #: Bring down the service
	docker compose down

db_dump: #: Dump the current database
	docker compose exec db pg_dump -U postgres alex_koin_dev > akc_backup

help: #: Show help topics
	@grep "#:" Makefile* | grep -v "@grep" | sort | sed "s/\([A-Za-z_ -]*\):.*#\(.*\)/$$(tput setaf 3)\1$$(tput sgr0)\2/g"
