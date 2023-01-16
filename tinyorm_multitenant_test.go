package tinyorm_test

import (
	"testing"

	tinyorm "github.com/BitlyTwiser/tinyORM"
)

var databaseConnections = []string{"development", "development-mysql"}

func TestMultiTenant(t *testing.T) {

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
