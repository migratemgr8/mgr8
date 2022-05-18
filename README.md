# mgr8

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

#### How to set needed variables

- **database_url**: set through `DB_HOST` environment variable or by using command flag `--database`.
- **driver**: set through command flag `--driver` (optional). Defaults to **postgres**.
- **migrations_dir**: set through command flag `--dir`.

### Execute migrations

Execute migrations by running
```bash
./bin/mgr8 apply <up|down>
```
Needs: **database_url**, **migrations_dir**, **driver**
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

Point to database by setting env **DB_HOST**.
<br/>
For postgres use DB_HOST=`postgres://root:root@localhost:5432/database_name?sslmode=disable`

### Testing

Use `make test`, `make display-coverage` and `make coverage-report`.

### Snippets

Executing migrations with postgres driver
```bash
./bin/mgr8 apply up --database=postgres://root:root@localhost:5432/core?sslmode=disable --dir=./migrations
```

Executing migrations with mysql driver
```bash
./bin/mgr8 apply up --database=root:root@tcp\(localhost:3306\)/core --dir=./migrations --driver=mysql
```

