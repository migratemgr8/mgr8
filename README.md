# mgr8

Agnostic tool that abstracts migration operations

Database URL example: `$DB_TYPE://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable`

Postgres example: `postgres://root:root@localhost:5432/database_name?sslmode=disable`

### Asdf

We recommend using asdf to manage versions. Asdf setup is described [here](https://asdf-vm.com/guide/getting-started.html) and can be summarized as follows:

```bash
# in any folder:
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.10.0
. $HOME/.asdf/asdf.sh # ideally, this will be added to ~/.bashrc
asdf plugin add golang

# in the project folder:
asdf install
```

### Build

Build with `make build`.

## How to run

### Open a database

The following command runs a database inside a container. 

```bash
docker compose up [-d] <database_name> 
```

Available databases: postgres, mysql

Passing the `-d` flag is optional will run the container in detached mode, what means it won't block the terminal, but you won't see database logs nor be able to close the container by using ctrl+c.

### Snippets

Executing migrations with postgres driver
```bash
./bin/mgr8 apply --database=postgres://root:root@localhost:5432/core?sslmode=disable ./migrations
```

Executing migrations with mysql driver
```bash
./bin/mgr8 apply --database=root:root@tcp\(localhost:3306\)/core ./migrations mysql
```
