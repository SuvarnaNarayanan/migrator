
# Migrator - dead simple SQL migration manager

I often see myself reaching out to ORM's just because they 'happen' to make migrations easier. Flyway and the like required me to install stuff that I didn't want to spend time on.

![intro](intro.gif)

## Features:

⭐ Migration table is decoupled from the actual database and is maintained in a separate sqlite db. Anytime anyone complains about why their migration isn't working - just copy and send the sqlite db over.

⭐ Zero installation, just download the binary apt for your hardware - [here](https://github.com/SuvarnaNarayanan/migrator/releases).

⭐ Multi-platform support - windows, linux. 

⭐ Language Agnostic - doesn't matter what your application code is. 

## Usage 

```
There's six commands 

1. init (i) - create initial config files
2. generate (g) - create a SQL file 
3. makemigrations (mm) - prime the SQL files
4. migrate (m) - actually execute the SQL files
5. fake - set the SQL files as migrated without actual execution of SQL
6. help

```

### $ migrator init

Creates a `migrations` directory if not found - this holds all the SQL files to be migrated. 

Creates a sample `migrations.config.yaml` file that needs to be altered before any other commands could be run.

```
migrations.config.yaml

migration:
    dbname: migrations
    tablename: migrations
    dir: migrations
targetdb:
    driver: sqlite | mysql | postgres
    datasource: ""
    username: "" // optional - not of any significance
    password: "" // optional - not of any significance

```

Make sure that the value for `datasource` is a valid connection string.

for example: 

```
... postgres
datasource: host=<host> port=<port> user=<user> password=<password> dbname=<dbname> sslmode=disable  
...

for more info - check https://pkg.go.dev/github.com/lib/pq

...sqlite
datasource: <path to db> 
...

...mysql
datasource: user:password@/dbname
...

for more info - check https://pkg.go.dev/github.com/go-sql-driver/mysql
```

### $ migrator generate 

Generate a valid SQL file to use with `migrator`. SQL files to be used with migrator needs to be named correctly. 

The format is as follows:

```
ID_DBOPERATION_DESCRIPTOR

ID: incrementing number
DB-OPERATION: db operation - CREATE, UPDATE, DELETE
DESCRIPTOR: short description of the operation

example: 1_CREATE_PERSON.sql

```

generate command looks for the last entry in migrations table and generates an appropriate id, incremented 1 from it. 

This command also inserts the record into the migrations table. 

### $ migrator makemigrations

Read all SQL files in the migrations directory and enters a record for each in the migration table with status as `PENDING`. 

The status can be either `COMPLETED, PENDING, FAILED`.

### $ migrator migrate

Executes the SQL files in `PENDING/FAILED` state in migration table.

### $ migrator fake

Fakes the pending migrations by setting the status as `COMPLETED` without actually executing any SQL.

## FAQ's

### How do I integrate this with an existing system?

You can capture the current state of the DB as `1_CREATE_INIT.sql` file and then mark it as migrated using `migrator fake` 