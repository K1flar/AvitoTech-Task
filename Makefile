DB_CONTAINER_NAME=avitotask-db-1

PG_NAME=postgres
PG_DB=banner_service

DSN=postgres://postgres:postgres@localhost:5432/$(PG_DB)?sslmode=disable

PATH_TO_MIGRATIONS=./db/migrations
MIGRATE=docker run --rm \
	-v $(PATH_TO_MIGRATIONS):/migrations \
	--network host \
	migrate/migrate \
	-path=/migrations/ \
	
all: build run 

build:
	docker-compose build server

run: 
	docker-compose up server 

stop:
	docker-compose down

db-start:
	docker-compose up db -d

db-stop:
	docker stop $(DB_CONTAINER_NAME)

migrate-up:
	$(MIGRATE) -database $(DSN) up

migrate-down:
	$(MIGRATE) -database $(DSN) down -all

testdata:
	docker exec -it $(DB_CONTAINER_NAME) psql -U $(PG_NAME) -d $(PG_DB) -f /testdata/testdata.sql  

integration-test:
	@echo "Starting test database..."
	@docker run --name test_db --rm -d -p 5432:5432 \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=test \
		postgres
	@echo "Wait postgres..."
	@until docker exec test_db pg_isready; do sleep 1; done
	@echo "Postgres is ready"
	@echo "Start migrations..."
	@$(MIGRATE) -database postgres://postgres:postgres@localhost:5432/test?sslmode=disable up
	@echo "Start tests..."
	@go test ./internal/integration/... -v -count=1
	@docker stop test_db