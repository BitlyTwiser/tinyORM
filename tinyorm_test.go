package tinyorm_test

import (
	"fmt"
	"testing"

	tinyorm "github.com/BitlyTwiser/tinyORM"
	"github.com/BitlyTwiser/tinyORM/pkg/custom"
	"github.com/google/uuid"
)

// Testing structs acting as database models
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Username string    `json:"username,omitempty"`
	Password string    `json:"password"`
	Age      int       `json:"age,omitempty"`
}

type Users []User

type Vehicle struct {
	ID            uuid.UUID    `json:"id"`
	Manufacturers custom.Slice `json:"manufacturers"`
	Data          custom.Map   `json:"data"`
	Color         string       `json:"color"`
	Recall        bool         `json:"recall"`
}

type Vehicles []Vehicle

// Testing the models with no IDs present
type TestNoID struct {
	Stuff string `json:"stuff"`
	Data  string `json:"data"`
}

type TestNoIDs []TestNoID

// Database from within the database.yml file to test against
const DATABASE = "development"

var (
	userID    = uuid.MustParse("4c0ea40b-4aeb-4b67-a407-4da25901ec8d")
	vehicleID = uuid.MustParse("4c0ea40b-4aeb-4b67-a407-4da25901ec8d")
)

// A select few tests need to be ran prior to the main series of tests. This subset is establishing 1 row for each struct that must exist for the Find queries.
// Though this does make the tests slightly more brittle, it also showcases the find by id functionality, thus was accepted
func TestCreate(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)
	if err != nil {
		t.Fatalf("error occurred connecting to database %s. %v", DATABASE, err.Error())
	}
	createTests := map[string]struct {
		action          string
		want            struct{}
		stmt            string
		limit           int
		args            any
		adjustModel     bool                  // Denotes if we run the adjustModelFun
		adjustModelFunc func(model any) error // Performs selected operations on model to alter model per test
		model           any
	}{
		"Test Create User with sepcific ID and without selected fields": {action: "create", adjustModel: false, model: createUserModel(true, userID)},
		"Test Create Vehicle with specific ID without selected fields":  {action: "create", adjustModel: false, model: createVehicleModel(true, vehicleID)},
	}

	for name, test := range createTests {
		t.Run(name, func(t *testing.T) {
			if err := db.Create(test.model); err != nil {
				t.Fatalf("error updating model. error: %v", err.Error())
			}
		})
	}

}

func TestORMFunctionality(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error occurred connecting to database %s. %v", DATABASE, err.Error())
	}

	tests := map[string]struct {
		action          string
		want            struct{}
		stmt            string
		limit           int
		args            any
		sliceArgs       []any
		adjustModel     bool                  // Denotes if we run the adjustModelFun
		adjustModelFunc func(model any) error // Performs selected operations on model to alter model per test
		model           any
	}{
		"Test Create User":               {action: "create", adjustModel: false, model: createUserModel(false)},
		"Test Create User without ID":    {action: "create", adjustModel: false, model: createUserModel(true)},
		"Test Create Vehicle":            {action: "create", adjustModel: false, model: createVehicleModel(false)},
		"Test Create Vehicle without ID": {action: "create", adjustModel: false, model: createVehicleModel(true)},
		"Test Update User":               {action: "update", adjustModel: true, model: createUserModel(true, userID), adjustModelFunc: UpdateUser},
		"Test Update Vehicle":            {action: "update", adjustModel: true, model: createVehicleModel(true, vehicleID), adjustModelFunc: UpdateVehicle},
		"Test Find User with ID":         {action: "find", adjustModel: false, model: new(User), args: userID},
		"Test Find Users":                {action: "find", adjustModel: false, model: new(Users)}, // Note the pluralization
		"Test Find User without ID":      {action: "find", adjustModel: false, model: new(User)},  // i.e. First functionality
		"Test Find Vehicle with ID":      {action: "find", adjustModel: false, model: new(Vehicle), args: vehicleID},
		"Test Find All Vehicles":         {action: "find", adjustModel: false, model: new(Vehicles)}, // Note the pluralization
		"Test Find Vehicle without ID":   {action: "find", adjustModel: false, model: new(Vehicle)},  // i.e. First functionality
		"Test Raw Query":                 {action: "raw-all", adjustModel: false, model: new(Vehicles), stmt: fmt.Sprintf("SELECT %s FROM vehicles", vehicleCoalesceQuery()), sliceArgs: []any{}},
		"Test Raw Query Singular":        {action: "raw-all", adjustModel: false, model: new(Vehicle), stmt: fmt.Sprintf("SELECT %s FROM vehicles LIMIT 1", vehicleCoalesceQuery()), sliceArgs: []any{}},
		"Test Raw Exec PSQL":             {action: "raw-exec", adjustModel: false, model: new(Vehicles), stmt: "insert into test_no_ids VALUES($1, $2)", sliceArgs: []any{"Things", "TestTest"}}, // Run this test for PSQL
		//"Test Raw Exec Mysql/Sqlite":       {action: "raw-exec", adjustModel: false, model: new(Vehicles), stmt: "insert into test_no_ids VALUES(?, ?)", sliceArgs: []any{"Things", "TestTest"}}, // Run this test for Mysql
		"Test Create model that has no id": {action: "create", adjustModel: false, model: createNodIdModel()},
		"Test Where User":                  {action: "where", adjustModel: false, model: new(User), limit: 0, stmt: "name = ?", args: "TestCreate"},
		"Test Where User using LIKE":       {action: "where", adjustModel: false, model: new(Users), limit: 2, stmt: "name LIKE ?", args: "%TestCreate%"},
		"Test Where Vehicle":               {action: "where", adjustModel: false, model: new(Vehicle), limit: 0, stmt: "color = ?", args: "Blue"},
		"Test Where Vehicle using LIKE":    {action: "where", adjustModel: false, model: new(Vehicles), limit: 2, stmt: "color LIKE ?", args: "%red%"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			switch test.action {
			case "create":
				if test.adjustModel {
					test.adjustModelFunc(test.model)
				}
				if err := db.Create(test.model); err != nil {
					t.Fatalf("error updating model. error: %v", err.Error())
				}
			case "update":
				if test.adjustModel {
					test.adjustModelFunc(test.model)
				}
				if err := db.Update(test.model); err != nil {
					t.Fatalf("error updating model. error: %v", err.Error())
				}
			case "find":
				if test.adjustModel {
					test.adjustModelFunc(test.model)
				}
				if test.args == nil {
					if err := db.Find(test.model); err != nil {
						t.Fatalf("error updating model. error: %v", err.Error())
					}
				} else {
					if err := db.Find(test.model, test.args); err != nil {
						t.Fatalf("error updating model. error: %v", err.Error())
					}
				}
			case "where":
				if test.adjustModel {
					test.adjustModelFunc(test.model)
				}
				if err := db.Where(test.model, test.stmt, test.limit, test.args); err != nil {
					t.Fatalf("error updating model. error: %v", err.Error())
				}
			case "raw-all":
				if test.adjustModel {
					test.adjustModelFunc(test.model)
				}
				if q, err := db.Raw(test.stmt, test.sliceArgs...); err == nil {
					if err := q.All(test.model); err != nil {
						t.Fatalf("error executing raw query. errorr: %s", err.Error())
					}
				} else {
					t.Fatalf("error executing raw query. errorr: %s", err.Error())
				}
			case "raw-exec":
				if test.adjustModel {
					test.adjustModelFunc(test.model)
				}
				if q, err := db.Raw(test.stmt, test.sliceArgs...); err == nil {
					if err := q.Exec(); err != nil {
						t.Fatalf("error executing raw query. %s", err.Error())
					}
				} else {
					t.Fatalf("error executing raw query. %s", err.Error())
				}
			}
		})
	}

}

// Similar case as the create tests with the Delete tests. Run these after all other tests have executed
func TestDeleteData(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)
	if err != nil {
		t.Fatalf("error occurred connecting to database %s. %v", DATABASE, err.Error())
	}

	deleteTests := map[string]struct {
		action          string
		want            struct{}
		stmt            string
		limit           int
		args            any
		adjustModel     bool                  // Denotes if we run the adjustModelFun
		adjustModelFunc func(model any) error // Performs selected operations on model to alter model per test
		model           any
	}{
		"Test Delete User by ID":                    {action: "delete", adjustModel: false, model: &User{ID: userID}},
		"Test Delete Vehicle by ID":                 {action: "delete", adjustModel: false, model: &Vehicle{ID: vehicleID}},
		"Test Delete TestWithNoIds with attributes": {action: "delete", adjustModel: false, model: &TestNoID{Stuff: "Things"}},
		"Test Delete User using attributes":         {action: "delete", adjustModel: false, model: &User{Name: "TestCreate"}},
		"Test Delete Vehicle using attributes":      {action: "delete", adjustModel: false, model: &Vehicle{Color: "Blue"}},
		"Test Should not Delete User":               {action: "delete", adjustModel: false, model: new(User)},
		"Test Should not Delete Vehicle":            {action: "delete", adjustModel: false, model: new(Vehicle)},
		"Test Delete Users":                         {action: "delete-bulk", adjustModel: false, model: new(Users)},    // Will Delete ALL Users
		"Test Delete Vehicles":                      {action: "delete-bulk", adjustModel: false, model: new(Vehicles)}, // Will Delete ALL Vehicles
		"Test Delete no ID model":                   {action: "delete", adjustModel: false, model: new(TestNoIDs)},
	}

	for name, test := range deleteTests {
		switch test.action {
		case "delete":
			t.Run(name, func(t *testing.T) {
				if err := db.Delete(test.model); err != nil {
					t.Fatalf("error updating model. error: %v", err.Error())
				}
			})
		case "delete-bulk":
			t.Run(name, func(t *testing.T) {
				if err := db.BulkDelete(test.model); err != nil {
					t.Fatalf("error updating model. error: %v", err.Error())
				}
			})

		}
	}
}

func createUserModel(withID bool, id ...uuid.UUID) *User {
	if !withID {
		return &User{
			Name:     "TestCreate",
			Email:    "TestCreate@email.com",
			Password: "password",
		}
	}

	if withID && len(id) == 0 {
		return &User{
			ID:       uuid.New(),
			Name:     "TestCreateUser2",
			Email:    "email@gmail.com",
			Username: "TestTest",
			Password: "password",
			Age:      111,
		}
	}

	return &User{
		ID:    id[0],
		Name:  "TestCreate3",
		Age:   10000,
		Email: "moreemail@gmail.com",
	}
}

func createVehicleModel(withID bool, id ...uuid.UUID) *Vehicle {
	if !withID {
		v := &Vehicle{
			Data:   make(custom.Map),
			Color:  "Blue",
			Recall: true,
		}
		v.Data.Add("Hello Testing", 123123)

		return v
	}

	if withID && len(id) == 0 {
		v := &Vehicle{
			ID:            uuid.New(),
			Manufacturers: custom.Slice{},
			Data:          make(custom.Map),
			Color:         "Red",
			Recall:        false,
		}

		v.Manufacturers = custom.Slice{"Ford", "Tesla", "Mercedes"}

		return v
	}

	return &Vehicle{
		ID:            id[0],
		Manufacturers: custom.Slice{},
		Data:          make(custom.Map),
		Color:         "Red",
		Recall:        false,
	}
}

func createNodIdModel() *TestNoID {
	return &TestNoID{Stuff: "Things", Data: "More things"}
}

func vehicleCoalesceQuery() string {
	return "COALESCE(id, '00000000-00000000-00000000-00000000'), COALESCE(manufacturers, '[]'), COALESCE(data, '{}'), COALESCE(color, ''), COALESCE(recall, false)"
}

// Updates pointer to user with different attributes
func UpdateUser(model any) error {
	if u, ok := model.(*User); ok {
		u.Name = "TestUpdate"
		u.Age = 9999
	} else {
		return fmt.Errorf("Error reflecting user pointer in UpdateUser test")
	}

	return nil
}

func UpdateVehicle(model any) error {
	m := custom.NewMap()
	m.Add("Two", "moremoremroemo")

	s := custom.Slice{"one", "two"}

	if v, ok := model.(*Vehicle); ok {
		v.Data = m
		v.Recall = true
		v.Manufacturers = s
	} else {
		return fmt.Errorf("Error reflecting value of vehicle model in UpdateVehicle")
	}

	return nil
}

var databaseConnections = []string{"development", "development-mysql"}

func TestMultiTenant(t *testing.T) {
	t.Skipf("Skipping multi-tenant test")
	mtc, err := tinyorm.MultiConnect(databaseConnections...)
	if err != nil {
		t.Fatal(err)
	}

	if err := mtc.SwitchDB("development").Create(&TestNoID{Stuff: "More Test PSQL"}); err != nil {
		t.Fatalf("error creating test on psqlDB. error: %v", err.Error())
	}

	if err := mtc.SwitchDB("development-mysql").Create(&TestNoID{Stuff: "More Test MySql"}); err != nil {
		t.Fatalf("error creating test on mysql. error: %v", err.Error())
	}
}
