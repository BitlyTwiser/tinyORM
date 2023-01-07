# tinyORM
A tiny ORM for all of your basic data layer needs

## Demands:
- Utilizing a simple database.yml file, one can enter multiple database connetions for the ORM to establish connections to.
i.e. Development, Production, ReadOnly endpoint, FluentD etc..

Database.yml:
```
---
development:
  dialect: postgres
  database: development 
  user: devUser 
  connect: true
  password: devPassword 
  host: 127.0.0.1
  port: 5432
  pool: 5

production-read-only:
  dialect: postgres
  database: ro 
  user: ro-user 
  password: ro-password 
  connect: false
  host: 127.0.0.1
  port: 5432
  pool: 5
```

- Connection without a flag will create a connection to EACH specific connection.
- if Connect is false, a connection will not be established. If true or missing, a connection will be attempted to the given database. 
- You can switch to the connection anytime using the <insert stuff> command. 
- Note: conflicting connection names will not work, only the first connection will be created.
- i.e.
```
development:
  dialect: postgres
  etc..

development: 
  dialct: mysql
  etc..
```
- only the postgres connection will be established, the repated connection will be ignored.


# Notes for MYSQL:
- AS MYSQL does not have UUID as a type, one must ensure they create their columns, if using UUID's as BINARY(36).
i.e.
```
create table users (id BINARY(36), name text, email text, username text, password text, age int);
```
- this will ensure that the UUID can be marshalled correctly.

# Notes for SQLITE3:
- This ORM does support the [Auth Feature](https://github.com/mattn/go-sqlite3#user-authentication).
- One must set the ```Auth``` flag to ```true``` within the database.yml file and compile the application with the auth flag.
```go build --tags sqlite_userauth```
- Note: The default ```_auth_crypt``` used to secure the SQLITE password is SHA512
