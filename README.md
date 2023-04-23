# Go Clean Architecture with gRPC

## Install Visual Studio Code extensions

```bash
code --install-extension esbenp.prettier-vscode
code --install-extension foxundermoon.shell-format
code --install-extension golang.go
code --install-extension matt-meyers.vscode-dbml
code --install-extension mtxr.sqltools
code --install-extension zxh404.vscode-proto3
```

## Schema Generation by `dbdocs` and `dbml2sql`

Create the [db.dbml](./docs/db.dbml) in to **decouple the design from a specific
database**. Then run `dbdocs` and `dbml2sql` to generate the docs and
[schema.sql](./docs/schemal.sql) boilerplate, we'll choose Postgres.

```bash
# Install nvs.
export NVS_HOME="$HOME/.local/share/nvs"
git clone https://github.com/jasongin/nvs "$NVS_HOME"
. "$NVS_HOME/nvs.sh" install

# Install node lts.
nvs add lts
nvs use lts
```

```bash
npm install -g dbdocs
dbdocs login
dbdocs build docs/db.dbml
dbdocs password —set secret —project go_clean_arch
```

```bash
npm install -g @dbml/cli
dbml2sql --postgres -o docs/schemal.sql docs/db.dbml
```

## Install [Docker](https://www.docker.com) and [PostgresSQL image](https://hub.docker.com/_/postgres).

```bash
# Install Docker.
brew install docker

# Run Docker app so that we can access the `docker` command.

# Pull the PostgresSQL image.
docker pull postgres:15.2-alpine

# Check the downloaded image.
docker images
```

## Run a Docker container using the official PostgresSQL image.

Creates and runs a Docker container with the name `postgres`, using the official
`postgres:15-alpine` Docker image. The container is started as a background
process (`-d` flag) and is mapped to port `5432` of the host machine
(`-p 127.0.0.1:5432:5432/tcp` flag), which is the default port for PostgreSQL.

The container is also configured with the environment variables `POSTGRES_USER`
and `POSTGRES_PASSWORD`, which set the default username and password for the
PostgreSQL database. In this case, the username is set to `root` and the
password is set to `password`.

```bash
docker run --name postgres \
  -p 127.0.0.1:5432:5432/tcp \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=password \
  -d postgres:15.2-alpine
```

```bash
# Enter the Postgres shell.
docker exec -it postgres psql -U root

# Try the following query in the shell.
SELECT now();
```

## Install [TablePlus](https://tableplus.com)

```bash
# Install TablePlus.
brew install tableplus
```

Connect to Postgres with the setting

![](https://i.imgur.com/jgHY7h3.png)

## Database Migration

```bash
# Install `migrate` command.
brew install golang-migrate

# Check the installed `migrate` command.
migrate --version

# Create the db migration directory.
mkdir -p db/migration

# Create the first migration script.
migrate create -ext sql -dir internal/db/migration -seq init_schema
```

Now, create a [Makefile](./Makefile) to save time and run the following:

```bash
# Run a PostgreSQL container.
make postgres

# Create a DB called "microservice" in this clean architecture.
make createdb

# Migrate up to create tables in the DB.
make migrateup
```

## Codegen via sqlc

```bash
# Install sqlc.
brew install sqlc

# Check the installed sqlc.
sqlc version
```

Initialize [`sqlc.yaml`](./sqlc.yaml) and copy the initial config from
[Getting started with PostgreSQL](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html#getting-started-with-postgresql)
with some modifications.

```bash
sqlc init
```

Add the queries in [user.sql](./internal/db/query/book.sql), then `make sqlc` to
codegen.

```bash
# Codegen.
make sqlc

# Eliminate red lines inside `./internal/repository/user.sql.go`.
go mod init github.com/walkccc/go-clean-arch
```

## Install [gRPC for Go](https://grpc.io/docs/languages/go/quickstart/)

```bash
# Install protobuf
brew install protobuf
```

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Add `"protoc": { "options": ["--proto_path=proto"] },` to
`~/Library/Application Support/Code/User/settings.json` to eliminate red lines.

```bash
# Generate gRPC code.
make proto

# Download required packages.
go mod tidy
```

## Install [Evans](https://github.com/ktr0731/evans) to test gRPC

```bash
# Install evans.
brew tap ktr0731/evans
brew install evans

# Start evans REPL.
make evans
```

## Install [GoMock](https://github.com/golang/mock).

```bash
# Install mockgen.
go install github.com/golang/mock/mockgen@v1.6.0

# Export the Go path.
# Add "export PATH=$PATH:~/go/bin" in your .zshrc or .bashrc

# Check the installed mockgen.
which mockgen
```

## Serve HTTP request with [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway)

1. Add [required protos](./api/proto/google/api/) from
   [googleapis/google/api](https://github.com/googleapis/googleapis/tree/master/google/api)
1. Add [required tools](./tools/tools.go)
1. Update annotation in [microservice.proto](./api/proto/microservice.proto)

```bash
make proto
```
