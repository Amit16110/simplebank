# postgres:
# 	docker run --name postgresdb -p 5432:5432 -e
createdb:
	docker exec -it postgresdb createdb -U postgres simplebank

dropdb:
	docker exec -it postgresdb dropdb simplebank

migrateUp:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simplebank?sslmode=disable" -verbose up
migrateDown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simplebank?sslmode=disable" -verbose down

cmdHistory:
	history | grep "docker run"

sqlc:
	sqlc generate

.PHONY: createdb dropdb cmdHistory migrateUp migrateDown