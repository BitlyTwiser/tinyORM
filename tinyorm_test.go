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

var Users []User

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

func TestDatabaseConnection(t *testing.T) {
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

	fUser := new(User)
	err = db.Find(fUser, 0)
	if err != nil {
		t.Fatalf("error finding user: %s", err.Error())
	}
	fmt.Println(u)
}
