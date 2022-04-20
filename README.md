# mgr8

Agnostic tool that abstracts migration operations

Database URL example: `$DB_TYPE://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable`

Postgres example: `postgres://root:root@localhost:5432/database_name?sslmode=disable`

### Asdf

Run `asdf install`.

### Build

Build with `make build`.

## How to run

### Open a database

```bash
docker compose up <database_name>
```

Available databases: postgres

### Run migrations

```bash
./bin/mgr8 apply --database=postgres://root:root@localhost:5432/core?sslmode=disable ./migrations
```