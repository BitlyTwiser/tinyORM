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
  password: devPassword 
  host: 127.0.0.1
  port: 5432
  pool: 5

production-read-only:
  dialect: postgres
  database: ro 
  user: ro-user 
  password: ro-password 
  host: 127.0.0.1
  port: 5432
  pool: 5
```

- Connection without a flag will create a connection to EACH specific connection.
- You can switch to the connection anytime using the db.conn command. 
