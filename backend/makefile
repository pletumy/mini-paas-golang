# DB Config
DB_CONTAINER_NAME = mini-paas-postgres
DB_PORT = 5432
DB_USER = postgres
DB_PASSWORD = 123
DB_NAME = mini_paas

###############################
db-up:
	docker run -d --name $(DB_CONTAINER_NAME)	\
	-e POSTGRES_USER=$(DB_USER)	\
	-e POSTGRES_PASSWORD=$(DB_PASSWORD)	\
	-e POSTGRES_DB=$(DB_NAME)	\
	-p $(DB_PORT):5432	\
	postgres:15

db-down:
	docker stop $(DB_CONTAINER_NAME) && docker rm $(DB_CONTAINER_NAME)

db-logs:
	docker logs -f $(DB_CONTAINER_NAME)

run:
	go run cmd/server/main.go
	
migrate:
	go run cmd/server/main.go migrate

test-integration:
	go test ./internal/tests/... -v