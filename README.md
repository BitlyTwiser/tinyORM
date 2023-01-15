# tinyORM
A tiny ORM for all of your basic data layer needs

## Premise:
- TinyORM was crafted with (hopefully) simplicity in mind. Having used many ORM's, the desire was to curate a functinal, yet simple minded ORM that could take care of generic transactions only utilizing the standard library and well known drivers. 

## Usage:

Create:
Update:
Find:
Delete:
Where:

## Custom Types:
- Natively, database/sql does not offer support for slices or maps.
- To accomodate for these datatypes, the ```custom``` pacakge was added.
- One can create custom types to utilize within their models akin to the ```Vehicle``` struct found in the tests.

Example:
```
type Vehicle struct {
	ID            uuid.UUID    `json:"id"`
	Manufacturers custom.Slice `json:"manufacturers"`
	Data          custom.Map   `json:"data"`
	Color         string       `json:"color"`
	Recall        bool         `json:"recall"`
}

// Creating a vehicle using the custom types:
  v := &Vehicle{
    ID:            uuid.New(),
    Manufacturers: custom.Slice{},
    Data:          make(custom.Map),
    Color:         "Red",
    Recall:        false,
  }
```

The custom types of customer.Slice{} has a built in ```Append``` method for inserting other types into the slice.

custom.Map also has methods for dealing with the underlying map structure. 
```
Add(key string, value any)
Delete(key strig)
```
methods exist on the custom.Map type to insert and delete records.

Both custom.Slice and custom.Map have a ```Values()``` method to return the contents of the data structures.

## Operational Notes:
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

# Package notes:
- This ORM uses google uuid to generate UUID's for the application, the UUID's may be expected whilst using structs as models


# Notes for MYSQL:
- AS MYSQL does not have UUID as a type, one must ensure they create their columns, if using UUID's as BINARY(36).
i.e.
```
// For users table
create table users (id BINARY(36), name text, email text, username text, password text, age int);

// For vehicles table 
create table vehicles (id BINARY(36), manufacturers json, data json, color text, recall bool);
```
- this will ensure that the UUID can be marshalled correctly.

# Notes for SQLITE3:
- This ORM does support the [Auth Feature](https://github.com/mattn/go-sqlite3#user-authentication).
- One must set the ```Auth``` flag to ```true``` within the database.yml file and compile the application with the auth flag.
```go build --tags sqlite_userauth```
- Note: The default ```_auth_crypt``` used to secure the SQLITE password is SHA512
