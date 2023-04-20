db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schemal.sql docs/db.dbml

.PHONY: db_docs db_shcema
