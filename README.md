## Backend Project with CRUD 

### Tech used
- Go (Gin) for server
- PostgreSQL for Database (using docker)
- SQLC (ORM for using the db)
- Viper (env Variables)
- gomock (to mock the database and test api)



### Setup Database

1) run (and pull) postgres image from docker
```
docker run --name bankapp -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -p 5432:5432 -d postgres
```

--- just for testing 
2) setup the schema in `schema.sql`

3) copy in docker and execute the command
```
docker cp schema.sql bankapp:/schema.sql
```

4) apply the schema
```
psql -U postgres -d postgres -f /schema.sql
```

then check the database schema

```
psql -U postgres -d postgres

```
```
\dt
```
5) database url here is: `postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable`

### Setup database migration

1) Run
```
migrate create -ext sql -dir db/migration -seq init_schema
```

2) Add postgres, createdb, dropdb command to the make file

3) now use migrate to apply migrations and add it to the make file
```
migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up
```


### Initialize the project

1) initialize the project
```
go mod init github.com/devrvk/simplebank
```

### gomock testing

1) generate the mock 
```
mockgen -package mockdb -destination db/mock/mock.go github.com/devrvk/simplebank/db/sqlc Store 
```

```
mockgen -package <package name in output> -destination <folder where the mock db is generated> <DB package location > < Interface to mock name>
```