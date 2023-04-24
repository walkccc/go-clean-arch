DB_URL=postgresql://root:password@localhost:5432/microservice?sslmode=disable

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schemal.sql docs/db.dbml

# Create and run a Docker container with the name `postgres`, using the official
# `postgres:15.2-alpine` Docker image.
postgres:
	docker run --name postgres \
			-p 127.0.0.1:5432:5432/tcp \
			-e POSTGRES_USER=root \
			-e POSTGRES_PASSWORD=password \
			-d postgres:15.2-alpine

# Create a DB called "microservice".
createdb:
	docker exec -it postgres createdb --username=root --owner=root microservice

# Drop a DB called "microservice".
dropdb:
	docker exec -it postgres dropdb microservice

# Migrate up to add tables in "microservice".
migrateup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

# Migrate down to drop tables in "microservice".
migratedown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

# Codegen CRUD code from "./internal/db/query/" to "./internal/db/sqlc/".
sqlc:
	sqlc --file ./internal/db/sqlc.yaml generate

# Run all tests and generate code coverage reports for all packages in the
# current module.
test:
	go test -v -cover ./...

# Generate gRPC code.
proto:
	rm pkg/*.go
	protoc --proto_path=api/proto \
			--go_out=pkg \
			--go_opt=paths=source_relative \
			--go-grpc_out=pkg \
			--go-grpc_opt=paths=source_relative \
			api/proto/*.proto

server:
	go run cmd/main.go

evans:
	evans -r repl

# Generate a mock implementation located at `internal/db/mock/store.go` of the
# `Store` interface in the
# `github.com/walkccc/go-clean-arch/internal/repository` package.
#
# This mock implementation can then be used for testing purposes.
mockgen:
	mockgen -package mock -destination internal/db/mock/store.go \
			github.com/walkccc/go-clean-arch/internal/repository Store

.PHONY: db_docs db_schema postgres createdb dropdb migrateup migratedown sqlc \
		test proto server evans mockgen
