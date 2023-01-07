# postgres:
# 	docker run --name postgresdb -p 5432:5432 -e
createdb:
	docker exec -it postgresdb createdb -U postgres simplebank

dropdb:
	docker exec -it postgresdb dropdb simplebank
# Not a root user
dropdbNonRoot:
	docker exec -it postgresdb psql -U postgres -d postgres -c "DROP DATABASE simplebank"
migrateUp:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simplebank?sslmode=disable" -verbose up
migrateDown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simplebank?sslmode=disable" -verbose down

cmdHistory:
	history | grep "docker run"

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go
.PHONY: createdb dropdb cmdHistory migrateUp migrateDown test server dropdbNonRoot