package db

import (
	"qlova.org/should"
	"qlova.org/should/test"
)

//TestablesViewer can be used to view the 'testable' table.
type TestablesViewer struct {
	View `db:"testable"`

	ID    Int64 `db:",key"`
	Value String
}

//LinkablesViewer can be used to view the 'linkable' table.
type LinkablesViewer struct {
	View `db:"linkable"`

	ID    Int64 `db:",key"`
	Value String
}

//TestSuite that can be applied to drivers to test for consistency.
//It uses the 'testable' table and as long as this table is not in use, the tests are isolated from the rest of the database.
//If the test passes, then there will be no trace of the test-data in the database.
//This makes it safe to use on databases that are in use.
type TestSuite struct {
	test.Suite

	Driver

	Testable TestablesViewer
	Linkable LinkablesViewer
}

func (ts *TestSuite) dummyRows() TestablesViewer {
	var t = ts.T()

	//Setup a few rows.
	var test = ts.Testable
	test.Value.Set("Hello")
	test.ID.Set(1)
	should.NotError(Insert(test)).Test(t)

	test = ts.Testable
	test.Value.Set("World")
	test.ID.Set(2)
	should.NotError(Insert(test)).Test(t)

	test = ts.Testable
	test.Value.Set("World")
	test.ID.Set(3)
	should.NotError(Insert(test)).Test(t)

	return test
}

func (ts *TestSuite) isolation() func() {
	var t = ts.T()

	ts.TestSync()

	return func() {
		should.NotError(
			Delete(&ts.Testable, &ts.Linkable),
		).Test(t)
	}
}

//SetupSuite gets called before any of the tests are run. Connecting the test tables.
func (ts *TestSuite) SetupSuite() {
	ts.Driver.Connect(&ts.Testable)
	ts.Driver.Connect(&ts.Linkable)
}

//TestConnect tests the driver view-connection functionality.
func (ts *TestSuite) TestConnect() {
	var Testable TestablesViewer

	ts.Driver.Connect(&Testable)

	should.Be(true)(Testable.Setup()).Test(ts.T())

	should.Be(true)(Testable.Master()).Test(ts.T())
}

//TestSync tests that the driver is able to sync the Testable table.
func (ts *TestSuite) TestSync() {
	should.NotError(
		Sync(ts.Testable),
	).Test(ts.T())

	should.NotError(
		Sync(ts.Linkable),
	).Test(ts.T())

	should.Error(
		Sync(TestablesViewer{}),
	).Test(ts.T())
}

func (ts *TestSuite) insert() {
	var t = ts.T()

	var test = ts.Testable
	test.ID.Set(1)
	test.Value.Set("TestValue")

	should.NotError(
		Insert(test),
	).Test(t)

	should.Error(
		Insert(TestablesViewer{}),
	).Test(ts.T())

}

//TestInsert tests that the driver is able to insert records into the database.
func (ts *TestSuite) TestInsert() {
	defer ts.isolation()()

	var t = ts.T()

	should.NotError(Empty(&ts.Testable)).Test(t)

	ts.insert()

	var result = ts.Testable

	should.NotError(
		If(ts.Testable.Value.Equals("TestValue")).Get(&result),
	).Test(t)

	should.Be("TestValue")(result.Value.Value()).Test(t)
}

//TestDelete tests that the driver is able to delete a table.
func (ts *TestSuite) TestDelete() {
	ts.TestSync()

	var t = ts.T()

	ts.insert()

	//Delete requires Master viewer as a data-safety precaution.
	clone := ts.Testable
	should.Error(
		Delete(&clone),
	).Test(t)

	should.NotError(
		Delete(&ts.Testable),
	).Test(t)

	//Deleted tables shouldn't have any rows.
	var result = ts.Testable
	should.Error(
		If(ts.Testable.Value.Equals("TestValue")).Get(&result),
	).Test(t)

	var test = ts.Testable
	test.Value.Set("TestValue")

	//You shoudn't be able to insert into a deleted table.
	should.Error(
		Insert(test),
	).Test(t)
}

//TestEmpty tests that the driver is able to empty a table.
func (ts *TestSuite) TestEmpty() {
	defer ts.isolation()()

	var t = ts.T()

	ts.insert()

	//Empty requires Master viewer as a data-safety precaution.
	clone := ts.Testable
	should.Error(
		Empty(&clone),
	).Test(t)

	should.NotError(
		Empty(&ts.Testable),
	).Test(t)

	//Emptied tables shouldn't have any matching rows.
	var result = ts.Testable
	should.Error(
		If(ts.Testable.Value.Equals("TestValue")).Get(&result),
	).Test(t)

	should.Be("")(result.Value.Value()).Test(t)
}

//TestPrimaryKey tests that the driver will reject duplicate values for a given primary key.
func (ts *TestSuite) TestPrimaryKey() {
	defer ts.isolation()()

	var t = ts.T()

	ts.insert()

	//Inserting a value with the same key should fail.
	var test = ts.Testable
	test.ID.Set(1)
	test.Value.Set("TestValue")

	should.Error(
		Insert(test),
	).Test(t)
}

//TestLink tests that the driver can link two tables together for filtering.
func (ts *TestSuite) TestLink() {
	defer ts.isolation()()

	var t = ts.T()

	var test = ts.Testable
	test.ID.Set(1)
	test.Value.Set("TestValue")

	should.NotError(
		Insert(test),
	).Test(t)

	var link = ts.Linkable
	link.ID.Set(1)
	link.Value.Set("LinkedValue")

	should.NotError(
		Insert(link),
	).Test(t)

	var result = ts.Linkable
	should.NotError(
		Link(test.ID.On(link.ID)).If(test.ID.Equals(1)).Get(&result),
	).Test(t)

	should.Be("LinkedValue")(result.Value.Value()).Test(t)
}
