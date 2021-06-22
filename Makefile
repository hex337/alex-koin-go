.PHONY: bash build deps help logs migrate ps restart seed setup_db stop test up db_dump

SERVICE ?= bot

SLACK_TOKEN ?= `grep 'SLACK_TOKEN=' .env | cut -d '=' -f2`
KOIN_BOT_ID ?= `grep 'KOIN_BOT_ID=' .env | cut -d '=' -f2`

default: help

all: build deps up setup_db migrate #: Boot up the project with db dependencies

bash: #: Bash prompt on running container
	docker-compose exec $(SERVICE) bash

build: #: Build containers
	touch .env
	docker-compose build

deps: #: Install the dependencies
	docker-compose run --rm -e MIX_ENV=test $(SERVICE) mix deps.get

format: #: Run mix format
	docker-compose run --rm $(SERVICE) mix format

iex: #: Interactive elixir shell on container
	docker-compose exec $(SERVICE) iex -S mix

logs: #: Tail the service container's logs
	docker-compose logs -tf $(SERVICE)

migrate: #: Run migrations
	docker-compose run --rm $(SERVICE) mix ecto.migrate

ps: #: Show running processes
	docker-compose ps

restart: #: Restart the service container
	docker-compose restart $(SERVICE)

seed: #: Seed the DB
	docker-compose exec -T $(SERVICE) mix run priv/repo/seeds.exs

setup_db: #: Create the db table(s)
	docker-compose run --rm $(SERVICE) mix ecto.create

stop: #: Stop running containers
	docker-compose stop

test: #: Run tests
	docker-compose run --rm -e MIX_ENV=test $(SERVICE) mix test

up: #: Start containers
	docker-compose up -d

down: #: Bring down the service
	docker-compose down

db_dump: #: Dump the current database
	docker-compose exec db pg_dump -U postgres alex_koin_dev > akc_backup

help: #: Show help topics
	@grep "#:" Makefile* | grep -v "@grep" | sort | sed "s/\([A-Za-z_ -]*\):.*#\(.*\)/$$(tput setaf 3)\1$$(tput sgr0)\2/g"
