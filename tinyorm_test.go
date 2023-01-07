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

// Database from within the database.yml file to test against
const DATABASE = "development-sqlite"

func TestCreateUser(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	u := &User{
		Name:     "yo",
		Email:    "penis@gmail.com",
		Password: "asdasdasd",
	}

	err = db.Create(u)
	if err != nil {
		t.Fatalf("error :%v", err.Error())
	}

	u2 := &User{
		ID:       uuid.New(),
		Name:     "yoyo",
		Email:    "yoyo@gmail.com",
		Username: "SupDawg",
		Password: "123123",
		Age:      2,
	}

	err = db.Create(u2)
	if err != nil {
		t.Fatalf("error :%v", err.Error())
	}

	u3 := &User{
		ID:    uuid.MustParse("4c0ea40b-4aeb-4b67-a407-4da25901ec8d"),
		Name:  "Carlton",
		Age:   10000,
		Email: "SupDawg@gmail.com",
	}

	err = db.Create(u3)

	if err != nil {
		t.Fatalf("error :%v", err.Error())
	}
}

func TestCreateVehicle(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}
	v := &Vehicle{
		ID:            uuid.MustParse("4c0ea40b-4aeb-4b67-a407-4da25901ec8d"),
		Manufacturers: custom.Slice{},
		Data:          make(custom.Map),
		Color:         "Red",
		Recall:        false,
	}

	v.Manufacturers.Append("Carl", "Sagan")

	err = db.Create(v)
	if err != nil {
		t.Fatalf("error creating vehicle %s", err.Error())
	}

	v2 := &Vehicle{
		Data:   make(custom.Map),
		Color:  "Blue",
		Recall: true,
	}

	v2.Data.Add("SupSup", 123123)

	err = db.Create(v2)

	if err != nil {
		t.Fatalf("error creating vehicle 2. %v", err.Error())
	}
}

func TestUpdateUser(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	u := &User{}
	err = db.Find(u, "4c0ea40b-4aeb-4b67-a407-4da25901ec8d")

	if err != nil {
		t.Fatalf("error finding user. Error %s", err.Error())
	}

	u.Name = "SupFool"
	u.Age = 420

	err = db.Update(u)

	if err != nil {
		t.Errorf("error updating user: %s", err.Error())
	}
}

func TestUpdateVehicle(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	v := new(Vehicle)

	err = db.Find(v, "4c0ea40b-4aeb-4b67-a407-4da25901ec8d")

	if err != nil {
		t.Fatalf("could not find vehicle. Error: %v", err.Error())
	}

	m := custom.NewMap()
	m.Add("Two", "moremoremroemo")

	s := custom.Slice{"one", "two"}

	v.Data = m
	v.Recall = true
	v.Manufacturers = s

	err = db.Update(v)

	if err != nil {
		t.Fatalf("error updating vehicle with id: %v. Error: %v", v.ID, err.Error())
	}
}

func TestFindUser(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	fUser := new(User)
	// // With ID
	err = db.Find(fUser, uuid.MustParse("4c0ea40b-4aeb-4b67-a407-4da25901ec8d"))
	if err != nil {
		t.Fatalf("error finding user: %s", err.Error())
	}

	fmt.Println(fUser)

	fUsers := new(Users)
	// No id passed, array is expected
	err = db.Find(fUsers)
	if err != nil {
		t.Fatalf("error finding users: %s", err.Error())
	}

	for _, user := range *fUsers {
		fmt.Printf("Found user id: %s", user.ID)
	}

	fmt.Println(fUsers)

	Users2 := new(User)
	err = db.Find(Users2)
	if err != nil {
		t.Fatalf("error finding users: %s", err.Error())
	}

	fmt.Println(Users2)

}

func TestFindVehicle(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	v := new(Vehicle)

	err = db.Find(v)

	if err != nil {
		t.Fatalf("error finding vehicle. Error: %v", err.Error())

	}

	t.Logf("Found first vehicle record. %v", v)

	v2 := new(Vehicles)
	err = db.Find(v2)

	if err != nil {
		t.Fatalf("error finding vehicles. Error: %v", err.Error())
	}

	t.Logf("found all vehicles. %v", v2)

	v3 := new(Vehicle)
	err = db.Find(v3, "4c0ea40b-4aeb-4b67-a407-4da25901ec8d")

	if err != nil {
		t.Fatal("could not find user")
	}

	t.Log(v3)
}

func TestWhere(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	u := new(User)

	err = db.Where(u, "name = ?", 0, "yo")

	if err != nil {
		t.Errorf("error obtaining record. error %v", err.Error())
	}

	fmt.Println(u)

	u2 := new(Users)

	//err = db.Where(u2, "name ILIKE ?", 0, "%yo%") // PSQL ILIKE support
	err = db.Where(u2, "name LIKE ?", 0, "%yo%")

	if err != nil {
		t.Errorf("error obtaining record. error %v", err.Error())
	}

	for _, user := range *u2 {
		fmt.Println(user)
	}
}
func TestDeleteUser(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}
	// u := &User{
	// 	ID:       uuid.MustParse("7feb1891-38f2-45b6-80d7-54e5d0217b78"),
	// 	Name:     "carl",
	// 	Email:    "stuffthings@gmail.com",
	// 	Username: "Hi",
	// 	Password: "asdasd",
	// 	Age:      111,
	// }

	u := &User{}

	err = db.Delete(u)

	if err != nil {
		t.Fatalf("error deleting user: %s", err.Error())
	}

}

func TestDeleteVehicle(t *testing.T) {
	db, err := tinyorm.Connect(DATABASE)

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	v := &Vehicle{}

	err = db.Delete(v)

	if err != nil {
		t.Fatalf("error deleting thing")
	}
}
