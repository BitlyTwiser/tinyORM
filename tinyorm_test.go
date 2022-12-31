package tinyorm_test

import (
	"fmt"
	"testing"

	tinyorm "github.com/BitlyTwiser/tinyORM"
)

// Testing structs acting as database models
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"Email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
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
	Data          map[string]any `json:"data"`
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

	u := &User{
		ID:       0,
		Name:     "carl",
		Email:    "stuffthings@gmail.com",
		Username: "Hi",
		Password: "asdasd",
		Age:      111,
	}

	err = db.Create(u)
	if err != nil {
		t.Fatalf("error :%v", err.Error())
	}

	u2 := &User{
		ID:       1,
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
	u := &User{
		ID:       0,
		Name:     "carl",
		Email:    "stuffthings@gmail.com",
		Username: "Hi",
		Password: "asdasd",
		Age:      111,
	}

	err = db.Delete(u)

	if err != nil {
		t.Fatalf("error deleting user: %s", err.Error())
	}
}

func TestUpdateUser(t *testing.T) {

}

func TestFindUser(t *testing.T) {
	db, err := tinyorm.Connect("development")

	if err != nil {
		t.Fatalf("error was had. %v", err.Error())
	}

	fUser := new(User)
	// // With ID
	err = db.Find(fUser, 1)
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
		fmt.Println(user)
	}

	fmt.Println(fUsers)
}
