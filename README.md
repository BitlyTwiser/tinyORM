# tinyORM
A tiny ORM for all of your basic data layer needs

## Premise:
- TinyORM was crafted with simplicity in mind. Having used many ORM's, the desire was to curate a functinal, yet simple ORM that could take care of generic transactions only utilizing the standard library and well known drivers. 

## Usage:

Create:
Update:
Find:
Delete:
Delete can be used to delete specifc database values or all values of a table.

To delete all values in a table, you can pass an empty slice of a specific model to delete.
Example:
```

type Vehicle struct {}
type Vehicles []Vehicle

v := new(Vehicles)

db.Delete(v) // Will delete ALL vehicles.
```

For the more common operation, deleteing single values, you can pass in a model with an ID of the object you wish to delete, or you can pass a model with attributes you wish to match on.
Example:
```
type User struct {
  ID uuid.UUID
  Age int
  Name string
}

user := &Users{Name: "Carl"}

db.Delete(user) // Will delete ALL users with the name of carl.


nextUser := &Users{ID: <uuid of user record>}

db.Delete(nextUser) // Will delete specific user with this ID.

tertiaryUser := &User{Name: "Bob", Age: 111}

db.Delete(tertiaryUser) // Will delete the user bob with age 111
``` 
A mixing of attributes to fit your deletion needs can be utilized or just specific ID's of the object. You can conjoin the delete functionality with Find/Where to discover a user with a query, then delete said user.
Note: If the ID is present, the attributes are ignored! Th record with the ID that matches will be deleted

Where:

Raw:
- Raw is really just that, a rather raw implementation giving most full control over to the user. No nil value safeguards are in place nor vetting of queries/attributes.
- When calling raw, you will have the pointer receiver methods available to you: ```All``` and ```Exec```. The rather common nomenclature for ORM's.
- ```All``` expects a model (or slice of models) and will insert the data into said model. Note: Model MUST be a pointer.
- ```Exec``` will simply execute a given query and that is all. 
- A snipper of raw functionality can be seen here:
```
				if q, err := db.Raw(test.stmt, test.sliceArgs...); err == nil {
					if err := q.Exec(); err != nil {
						t.Fatalf("error executing raw query. %s", err.Error())
					}
				}
```
The ```stmt``` is the query, followed by arguments to supplament the query with data as needed.
i.e. ```stmt := "select * from foo"```
Examples of functionality are within the tinyorm_test.go

## Custom Types:
- Natively, database/sql does not offer support for slices or maps.
- To accommodate for these datatypes, the ```custom``` package was added.
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
Example:
```
		v := &Vehicle{
			ID:            uuid.New(),
			Manufacturers: custom.Slice{},
			Data:          make(custom.Map),
			Color:         "Red",
			Recall:        false,
		}

		v.Manufacturers = custom.Slice{"Ford", "Tesla", "Mercedes"}

		return v
```

custom.Map also has methods for dealing with the underlying map structure. 
```
Add(key string, value any)
Delete(key strig)
```
Example:
```
		v := &Vehicle{
			Data:   make(custom.Map),
			Color:  "Blue",
			Recall: true,
		}
		v.Data.Add("Hello Testing", 123123)

		return v
```
methods exist on the custom.Map type to insert and delete records.

Both custom.Slice and custom.Map have a ```Values()``` method to return the contents of the data structures.


## Null values:
- The ```database/sql``` pacakge does not handle nil values in the Scan functionality. The ```Custom``` package does supply the user with the ability to utilize slices and maps, the primary code wraps all queryes for model attributes into a COALESCE statement.
- If one desire, you can also utilize the SQL package sql.NullStrings, sqlNullBool, etc...
- You can also utilize a pointer to the asset that may be nil on the model.
```
type TestModel struct {
  Age int
  Name *string // Note pointer usage
}
```
- This performs identicall to the sql.NullString implementation. Quoting Russ Cox:
https://groups.google.com/g/golang-nuts/c/vOTFu2SMNeA/m/GB5v3JPSsicJ
```
There is no effective difference.
``` 
- Operationally, tinyorm handles the nil values by default using the coalesce, this is baked into the application, so the user will not have to accoint for nil values unless you are using the ```Raw``` functionality, no guards are in place there to protect the user.

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

## Multi-Tenant connections:
- tinyORM has the ability to connect and keep-alive multiple connections to different databases.
- Utilizing the multi-connect utility, you can connect to multiple databases and switch between them easily.

Example:
```
	mtc, err := tinyorm.MultiConnect(databaseConnections...)
	if err != nil {
		t.Fatal(err)
	}

	if err := mtc.SwitchDB("development").Create(&TestNoID{Stuff: "More Test PSQL"}); err != nil {
		t.Fatalf("error creating test on psqlDB. error: %v", err.Error())
	}
```
- The above example is pulled from the tinyorm_multitenant_test. 
- Utilizing the ```MultiConnect``` function, you can use the methods built into the dialects.MultiTenantDialectHandler{} struct.

```
type MultiTenantDialectHandler struct {
	Handlers map[string]DialectHandler
}

// Append will add database handlers to the Handlers slice
func (mtd *MultiTenantDialectHandler) Set(key string, handler DialectHandler) {
	mtd.Handlers[key] = handler
}

// Empty will determine if there are not database handlers present
func (mtd MultiTenantDialectHandler) Empty() bool {
	return len(mtd.Handlers) == 0
}

// Switch allows the caller to alter to different databases to perform executions again
func (mtd MultiTenantDialectHandler) SwitchDB(database string) DialectHandler {
	if db, found := mtd.Handlers[database]; found {
		return db
	}

	return nil
}
```

## Setting timeouts:
tinyorm allows for the setting of the following database/sql values for open/idle connections:
```
SetConnMaxIdleTime
SetConnMaxLifetime
SetMaxIdleConns
SetMaxOpenConns
```
These can be set wtihin the database.yml file per connnection. If left blank, defaults will be used.
Example:
```
development:
  dialect: postgres
  database: tinyorm
  user: tiny 
  password: password123!
  connect: true
  host: 127.0.0.1
  port: 5432
  name: "development"
  maxIdleTime: 60
  maxLifetime: 100
  maxIdleConn: 0
  maxOpenConn: 10
```

## Package notes:
- This ORM uses google uuid to generate UUID's for the application, the UUID's may be expected whilst using structs as models


## Notes for MYSQL:
- AS MYSQL does not have UUID as a type, one must ensure they create their columns, if using UUID's as BINARY(36).
i.e.
```
// For users table
create table users (id BINARY(36), name text, email text, username text, password text, age int);

// For vehicles table 
create table vehicles (id BINARY(36), manufacturers json, data json, color text, recall bool);
```
- this will ensure that the UUID can be marshalled correctly.

## Notes for SQLITE3:
- This ORM does support the [Auth Feature](https://github.com/mattn/go-sqlite3#user-authentication).
- One must set the ```Auth``` flag to ```true``` within the database.yml file and compile the application with the auth flag.
```go build --tags sqlite_userauth```
- Note: The default ```_auth_crypt``` used to secure the SQLITE password is SHA512
- Auth is not enabled by default and the flag does have to be used in order for Auth feature to function.
