package tinyorm_test

import (
	"fmt"
	"testing"

	tinyorm "github.com/BitlyTwiser/tinyORM"
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

type Dog struct {
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Color  string `json:"color"`
	Name   string `json:"name"`
}

type Vehicle struct {
	Manufacturers []string       `json:"manufacturers"`
	Data          map[string]int `json:"data"`
	Color         string         `json:"color"`
	Recall        bool           `json:"recall"`
}

func TestInitializeDatabase(t *testing.T) {
	tests := []struct {
		name string
		have string
		want string
	}{
		{
			name: "",
			have: "",
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.have != test.want {
				t.Errorf("test %s failed. Have: %v Want: %v", test.name, test.have, test.want)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	db, err := tinyorm.Connect("development")

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	// u := &User{
	// 	ID:       uuid.New(),
	// 	Name:     "carl",
	// 	Email:    "stuffthings@gmail.com",
	// 	Username: "Hi",
	// 	Password: "asdasd",
	// 	Age:      111,
	// }

	u := &User{
		Name:     "yo",
		Email:    "penis@gmail.com",
		Password: "asdasdasd",
	}

	err = db.Create(u)
	if err != nil {
		t.Fatalf("error :%v", err.Error())
	}

	// v := &Vehicle{
	// 	Manufacturers: []string{"Ford", "Tesla"},
	// 	Data:          map[string]int{"asdasd": 10},
	// 	Color:         "Red",
	// 	Recall:        false,
	// }

	// err = db.Create(v)
	// if err != nil {
	// 	t.Fatalf("error creating vehicle %s", err.Error())
	// }

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
}

func TestDeleteUser(t *testing.T) {
	db, err := tinyorm.Connect("development")

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

	// v := &Vehicle{Color: "Red"}

	// err = db.Delete(v)

	// if err != nil {
	// 	t.Fatalf("error deleting thing")
	// }
}

func TestUpdateUser(t *testing.T) {
	db, err := tinyorm.Connect("development")

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	u := &User{}
	err = db.Find(u, "418caee0-fce1-431e-a26b-e73b84750f37")

	if err != nil {
		t.Fatalf("error finding user. Error %s", err.Error())
	}

	u.Name = "SomethingElse"
	u.Age = 42069

	err = db.Update(u)

	if err != nil {
		t.Errorf("error updating user: %s", err.Error())
	}
}

func TestFindUser(t *testing.T) {
	db, err := tinyorm.Connect("development")

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	fUser := new(User)
	// // With ID
	err = db.Find(fUser, uuid.MustParse("418caee0-fce1-431e-a26b-e73b84750f37"))
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
}
