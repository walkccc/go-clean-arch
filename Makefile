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

.PHONY: db_docs db_schema postgres createdb dropdb migrateup migratedown
