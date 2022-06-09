# mgr8

[![codecov](https://codecov.io/gh/kenji-yamane/mgr8/branch/master/graph/badge.svg?token=WUJ54P2TQQ)](https://codecov.io/gh/kenji-yamane/mgr8)

An agnostic tool that abstracts migration operations

### Suported databases
- **postgres**
- **mysql**

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

#### Set needed variables

- **database**: Database URL. <br/>
Set through `DB_HOST` environment variable or by using command flag `--database`.

- **driver**: Driver name. <br/>
Defaults to **postgres**. <br/>
Set through command flag `--driver`.

- **dir**: Migrations directory. <br/>
Set through command flag `--dir`.
### Execute migrations

Requires: **database**, **dir**, **driver** <br/>
Execute migrations by running
```bash
./bin/mgr8 apply <up|down> <number_of_migrations>
```
- number_of_migrations: Number of migrations to run (Optional). If not specified, runs only one.

### Generate migrations

Mgr8 keeps a copy of the latest schema at `.mgr8/reference.sql`. Copy the first schema version to this location, or simply run `./bin/mgr8 generate init <schemafile>` which will do the same.

Then run `./bin/mgr8 generate diff <schemafile>` to generate migrations with respect to the reference file. This will also update the reference.

When committing to a repository, check if the reference and the latest schema match. The command `./bin/mgr8 generate check <schemafile>` can be used, as it returns 0 if the files match.

To generate an empty migration (e.g. for DML), use `./bin/mgr8 generate empty`.

### Run with docker

Pull latest image with `docker pull migratemgr8/mgr8:latest` or build it yourself with `make build-docker-image`.

Run commands:
```bash
docker run -v {{ migrations path }}:/migrations --network host -e MGR8_USERNAME={{ logs username }} -e DB_HOST={{ database connection string }} migratemgr8/mgr8 <command>
```
Make sure to replace the variables surrounded by double curly braces.

## Develop

### Requirements
- [Asdf](https://asdf-vm.com/guide/getting-started.html)
- [Asdf golang plugin](https://github.com/kennyp/asdf-golang)
- [Docker compose](https://docs.docker.com/compose/install/)

Run `make install-tools` to install tooling dependencies.

### Run a database container

Run a testing database with
```bash
docker compose up [-d] <database_name> 
```
Available databases: postgres, mysql

Passing the `-d` flag is optional and will run the container in detached mode, it won't block the terminal but you won't see database logs nor be able to close the container by using ctrl+c.

Point to database by setting env **DB_HOST**.
<br/>
For postgres use DB_HOST=`postgres://root:root@localhost:5432/database_name?sslmode=disable`

### Testing

Use `make mock` to generate necessary mocks.

Use `make test`, `make display-coverage` and `make coverage-report`.

To add a new mock, add new lines to the `mock` command in the Makefile.

### Snippets

Executing migrations with postgres driver
```bash
# CLI version
./bin/mgr8 apply up --database=postgres://root:root@localhost:5432/core?sslmode=disable --dir=./migrations
# Docker version
docker run -v $PWD/migrations:/migrations -e DB_HOST=postgres://root:root@localhost:5432/core?sslmode=disable --network host -e MGR8_USERNAME=username migratemgr8/mgr8 apply up
```

Executing migrations with mysql driver
```bash
# CLI version
./bin/mgr8 apply up --database=root:root@tcp\(localhost:3306\)/core --dir=./migrations --driver=mysql
# Docker version
docker run -v $PWD/migrations:/migrations -e DB_HOST=postgres://root:root@localhost:5432/core?sslmode=disable --network host -e MGR8_USERNAME=username migratemgr8/mgr8 apply up --driver=mysql
```

