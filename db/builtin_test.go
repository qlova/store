package db

import (
	"testing"

	"qlova.org/should"
	"qlova.org/should/test"
)

func Test_Builtin(t *testing.T) {
	var driver = Open()

	test.New(&TestSuite{
		Driver: driver,
	})(t)

	should.NotError(driver.Close()).Test(t)
}
