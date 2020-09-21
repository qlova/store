package postgres_test

import (
	"os"
	"testing"

	"qlova.org/should"
	"qlova.org/should/test"
	"qlova.store/db"
	"qlova.store/db/driver/postgres"
)

func Test_Open(t *testing.T) {

	var password = os.Getenv("POSTGRES_PASSWORD")

	if password == "" {
		t.Skip()
	}

	var driver = postgres.Open("host=localhost sslmode=disable user=postgres dbname=postgres password=" + password + " port=5433")

	test.New(&db.TestSuite{
		Driver: driver,
	})(t)

	should.NotError(driver.Close()).Test(t)
}
