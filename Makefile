CONFIG_FILE=./configs/config.yaml

# DB
DB_SERVICE_NAME=db
DB_CONTAINER_NAME=avitotask-db-1

PG_NAME=postgres
PG_DB=banner_service

DSN=$(shell sed -n 's/dsn: //p' $(CONFIG_FILE))

PATH_TO_MIGRATIONS=./db/migrations
MIGRATE=docker run --rm \
	-v $(PATH_TO_MIGRATIONS):/migrations \
	--network host \
	migrate/migrate \
	-path=/migrations/ \
	-database $(DSN) 

stop:
	docker-compose down

db-start:
	docker-compose up $(DB_SERVICE_NAME) -d

db-stop:
	docker stop $(DB_CONTAINER_NAME)

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down -all

testdata:
	docker exec -it $(DB_CONTAINER_NAME) psql -U $(PG_NAME) -d $(PG_DB) -f /testdata/testdata.sql  
