all: destroy build init start

build:
	docker-compose build --no-cache

init: reinit_db config accounts

start:
	docker-compose up

restart: init start

config:
	rm -rf .bigchaindb_workshop_docker
	touch .bigchaindb_workshop_docker
	docker-compose run --rm bdb bigchaindb -yc .bigchaindb_workshop configure
	docker-compose run --rm bdb bigchaindb -c .bigchaindb_workshop init

accounts:
	docker-compose run --rm bdb python -m bigchaindb_workshop.backend.init_accounts

drop_db:
	docker-compose stop rdb
	docker-compose rm -f rdb

start_db:
	docker-compose up -d rdb

reinit_db: drop_db start_db
	sleep 10

stop:
	docker-compose down

destroy:
	docker-compose rm -f
