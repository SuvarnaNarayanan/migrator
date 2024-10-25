
# Migrator - dead simple SQL migration manager

I often see myself reaching out to ORM's just because they 'happen' to make migrations easier. Flyway and the like required me to install stuff that I didn't want to spend time on.

## Features:

⭐ Migration table is decoupled from the actual database and is maintained in a separate sqlite3 db. Anytime anyone complains about why their migration isn't working - just copy and send the sqlite3 db over.

⭐ Zero installation, just download the binary apt for your hardware - here.

⭐ Multi-platform support - windows, linux. 

⭐ Language Agnostic - doesn't matter what your application code is. 

## Usage 

```
There's five commands 

1. init (i) - create initial config files
2. generate (g) - create a SQL file 
3. makemigrations (mm) - prime the SQL files
4. migrate (m) - actually execute the SQL files
5. fake - set the SQL files as migrated without actual execution of SQL

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
    driver: sqlite3 | mysql | postgres
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

...sqlite3
datasource: <path to db> 
...

...mysql
datasource: user:password@/dbname
...

for more info - check https://pkg.go.dev/github.com/go-sql-driver/mysql
```

