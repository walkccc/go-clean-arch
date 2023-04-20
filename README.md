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
