# mgr8

An agnostic tool that abstracts migration operations

## How to use

### Requirements

- [Asdf](https://asdf-vm.com/guide/getting-started.html)
- [Asdf golang plugin](https://github.com/kennyp/asdf-golang)

### Setup

Make sure you use the project's golang version by running
```bash
asdf install
```

Build project by running
```bash
make build
```

### Run commands

### Execute migrations

Execute migrations by running
```bash
./bin/mgr8 <migrations_folder> <driver>
```
Currently supported drivers: **postgres** and **mysql**.
<br/>
Defaults to **postgres**.
<br/>
Make sure you either have `DB_HOST` environment variable with you database connection string or you pass it in by using the flag `--database`.
<br/>
Example connection string: `postgres://root:root@localhost:5432/database_name?sslmode=disable`

## Develop

### Requirements
- [Asdf](https://asdf-vm.com/guide/getting-started.html)
- [Asdf golang plugin](https://github.com/kennyp/asdf-golang)
- [Docker compose](https://docs.docker.com/compose/install/)

### Run a database container

Run a testing database with
```bash
docker compose up [-d] <database_name> 
```
Available databases: postgres, mysql

Passing the `-d` flag is optional and will run the container in detached mode, it won't block the terminal but you won't see database logs nor be able to close the container by using ctrl+c.

### Testing

Use `make test`, `make display-coverage` and `make coverage-report`.

### Snippets

Executing migrations with postgres driver
```bash
./bin/mgr8 apply --database=postgres://root:root@localhost:5432/core?sslmode=disable ./migrations
```

Executing migrations with mysql driver
```bash
./bin/mgr8 apply --database=root:root@tcp\(localhost:3306\)/core ./migrations mysql
```

