package db

import (
	"math"
	"testing"

	"qlova.org/should"
)

//TestTypes tests the db types.
func (ts *TestSuite) TestTypes() {
	defer ts.isolation()()

	var t = ts.T()

	//General Tests for each type.
	var Types = []interface {
		Test(*testing.T)
	}{
		Int8{}, Int16{}, Int32{}, Int64{},

		Float32{}, Float64{},

		Rune{},

		Bool{}, Bytes{}, String{},

		Time{},
	}
	for _, T := range Types {
		T.Test(t)
	}
}

//TestResultsDelete tests the deletion of rows.
func (ts *TestSuite) TestResultsDelete() {
	defer ts.isolation()()

	var t = ts.T()

	//Setup a few rows.
	var test = ts.dummyRows()

	should.NotError(
		If(test.Value.Equals("World")).Delete(),
	).Test(t)

	var result = ts.Testable
	should.Error(
		If(ts.Testable.Value.Equals("World")).Get(&result),
	).Test(t)

	should.Be("")(result.Value.Value()).Test(t)

	//Shouldn't delete anything else.
	result = ts.Testable
	should.NotError(
		If(ts.Testable.Value.Equals("Hello")).Get(&result),
	).Test(t)

	should.Be("Hello")(result.Value.Value()).Test(t)

	//The last delete shouldn't error or panic.
	should.NotPanic(func() {
		should.NotError(
			If(test.Value.Equals("Hello")).Delete(),
		).Test(t)
	}).Test(t)
}

//TestResultsUpdate tests the updating of rows.
func (ts *TestSuite) TestResultsUpdate() {
	defer ts.isolation()()

	var t = ts.T()

	//Setup a few rows.
	var test = ts.dummyRows()

	should.NotError(
		If(test.Value.Equals("Hello")).Update(
			test.Value.To("olleH"),
		),
	).Test(t)

	var result = ts.Testable
	should.NotError(
		If(ts.Testable.Value.Equals("olleH")).Get(&result),
	).Test(t)

	should.Be("olleH")(result.Value.Value()).Test(t)

	//Shouldn't update anything else.
	result = ts.Testable
	should.Error(
		If(ts.Testable.Value.Equals("Hello")).Get(&result),
	).Test(t)

	//Shouldn't update anything else.
	result = ts.Testable
	should.NotError(
		If(ts.Testable.Value.Equals("World")).Get(&result),
	).Test(t)

	should.Be("World")(result.Value.Value()).Test(t)
}

//TestResultsRead tests the retrieval of given columns.
func (ts *TestSuite) TestResultsRead() {
	defer ts.isolation()()

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

	should.NotError(
		If(test.Value.Equals("Hello")).Read(&test.Value),
	).Test(t)

	should.Be("Hello")(test.Value.Value()).Test(t)
	should.Be(int64(0))(test.ID.Value()).Test(t)
}

//TestResultsGet tests the retrieval of rows.
func (ts *TestSuite) TestResultsGet() {
	defer ts.isolation()()

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

	should.NotError(
		If(test.Value.Equals("Hello")).Get(&test),
	).Test(t)

	should.Be("Hello")(test.Value.Value()).Test(t)
	should.Be(int64(1))(test.ID.Value()).Test(t)
}

//TestResultsJSON tests the json marshalling of rows.
func (ts *TestSuite) TestResultsJSON() {
	defer ts.isolation()()

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

	json, err := If(test.Value.Equals("Hello")).MarshalJSON()
	should.NotError(err).Test(t)

	should.Be(`[{"ID":1,"Value":"Hello"}]`)(string(json)).Test(t)

	json, err = If(test.Value.NotEquals("")).SortBy(test.Value.Increasing()).MarshalJSON()
	should.NotError(err).Test(t)

	should.Be(`[{"ID":1,"Value":"Hello"},{"ID":2,"Value":"World"}]`)(string(json)).Test(t)
}

//TestResultsCount tests the results counting function.
func (ts *TestSuite) TestResultsCount() {
	defer ts.isolation()()

	var t = ts.T()

	//Setup a few rows.
	var test = ts.dummyRows()

	count, err := If(test.Value.Equals("Hello")).Count(test.Value)
	should.NotError(err).Test(t)

	should.Be(1)(count).Test(t)

	count, err = If(test.Value.Equals("World")).Count(test.Value)
	should.NotError(err).Test(t)

	should.Be(2)(count).Test(t)

	count, err = If(test.Value.Equals("")).Count(test.Value)
	should.NotError(err).Test(t)

	should.Be(0)(count).Test(t)
}

//TestResultsSum tests the result sum function.
func (ts *TestSuite) TestResultsSum() {
	defer ts.isolation()()

	var t = ts.T()

	//Setup a few rows.
	var test = ts.dummyRows()

	should.NotError(If(test.Value.Equals("Hello")).Sum(&test.ID)).Test(t)
	should.Be(int64(1))(test.ID.Value()).Test(t)

	should.NotError(If(test.Value.Equals("World")).Sum(&test.ID)).Test(t)
	should.Be(int64(5))(test.ID.Value()).Test(t)

	should.NotError(If(test.Value.Equals("")).Sum(&test.ID)).Test(t)
	should.Be(int64(0))(test.ID.Value()).Test(t)
}

//TestResultsAverage tests the result average function.
func (ts *TestSuite) TestResultsAverage() {
	defer ts.isolation()()

	var t = ts.T()

	//Setup a few rows.
	var test = ts.dummyRows()

	avg, err := If(test.Value.Equals("Hello")).Average(test.ID)
	should.NotError(err).Test(t)
	should.Be(float64(1))(avg).Test(t)

	avg, err = If(test.Value.Equals("World")).Average(test.ID)
	should.NotError(err).Test(t)
	should.Be(float64(2.5))(avg).Test(t)

	avg, err = If(test.Value.Equals("")).Average(test.ID)
	should.NotError(err).Test(t)
	should.Be(true)(math.IsNaN(avg)).Test(t)
}

//TestResultsSlice tests filter slicing.
func (ts *TestSuite) TestResultsSlice() {
	defer ts.isolation()()

	var t = ts.T()

	//Setup a few rows.
	var test = ts.dummyRows()

	should.NotError(
		If(test.Value.Equals("World")).SortBy(test.ID.Increasing()).Slice(0, 2).Into(&test),
	).Test(t)

	i := Range(&test)

	i.Next()
	should.Be(int64(2))(test.ID.Value()).Test(t)
	i.Next()
	should.Be(int64(3))(test.ID.Value()).Test(t)

	should.NotError(
		If(test.Value.NotEquals("")).SortBy(test.ID.Increasing()).Slice(0, 3, &test.ID).Read(),
	).Test(t)

	i = Range(&test)

	i.Next()
	should.Be(int64(1))(test.ID.Value()).Test(t)
	i.Next()
	should.Be(int64(2))(test.ID.Value()).Test(t)
	i.Next()
	should.Be(int64(3))(test.ID.Value()).Test(t)
}

//TestResultsNotFound tests that ErrNotFound is returned if a record does not exist.
func (ts *TestSuite) TestResultsNotFound() {
	defer ts.isolation()()

	var t = ts.T()

	//Setup a few rows.
	var test = ts.dummyRows()

	should.Be(ErrNotFound)(
		If(test.Value.Equals("DOES NOT EXIST")).SortBy(test.ID.Increasing()).Get(&test),
	).Test(t)
}
