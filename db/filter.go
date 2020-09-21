package db

import (
	"reflect"
)

//Linker links two tables together so that they can be searched on.
type Linker struct {
	From, To Viewable

	View Table
}

//Sorter defines how to sort a given column.
type Sorter struct {
	Column     string
	Decreasing bool
}

//Slicer is returned by Filter.Slice
type Slicer Filter

//Into moves the filter's selection into the given viewer.
func (s Slicer) Into(v Viewer) (int, error) {
	var columns = make([]Variable, v.Columns())
	var rviewer = reflect.ValueOf(v).Elem()

	for i := 0; i < v.Columns(); i++ {
		column := v.Column(i)

		columns[i] = rviewer.Field(column.Field()).Addr().Interface().(Variable)
	}

	s.Columns = columns

	return v.Database().Search(Filter(s)).Get(columns[0], columns[1:]...)
}

//Read reads into the given variables.
func (s Slicer) Read() (int, error) {
	return Filter(s).driver.Search(Filter(s)).Get(s.Columns[0], s.Columns[1:]...)
}

//MarshalJSON implements json.Marshaler
func (s Slicer) MarshalJSON() ([]byte, error) {
	return Filter(s).MarshalJSON()
}

//Filter describes which rows to select in a database.
type Filter struct {
	driver Driver

	Table string
	View  Table

	Condition  Condition
	Conditions []Condition

	Sort  Sorter
	Sorts []Sorter

	Link  Linker
	Links []Linker

	Offset, Length int

	Columns []Variable
}

//Link returns a filter with the given links applied.
func Link(linker Linker, linkers ...Linker) Filter {
	return Filter{
		Table:  linker.From.Table(),
		driver: linker.From.Database(),
		View:   linker.View,

		Link:  linker,
		Links: linkers,

		Length: 1,
	}
}

//If returns a filter on the database with the given conditions.
func If(condition Condition, conditions ...Condition) Filter {
	return Filter{
		Table:      condition.Table,
		View:       condition.View,
		driver:     condition.driver,
		Condition:  condition,
		Conditions: conditions,

		Length: 1,
	}
}

//If returns a filter on the database with the additional conditions.
func (f Filter) If(condition Condition, conditions ...Condition) Filter {
	if f.Table == "" {
		f.Table = condition.Table
	}

	if f.View == nil {
		f.View = condition.View
		f.driver = condition.driver
	}

	f.Condition = condition
	f.Conditions = conditions

	return f
}

//Slice selects the slice of the results that this filter should return.
func (f Filter) Slice(offset, length int, columns ...Variable) Slicer {
	f.Offset = offset
	f.Length = length
	f.Columns = columns
	return Slicer(f)
}

//Get moves the filter's selection into the given viewer.
func (f Filter) Get(v Viewer) error {
	if v.Master() {
		return ErrIllegalMaster
	}
	if v.Database() == nil {
		return ErrDisconnectedViewer
	}

	//No columns, nothing to get.
	if v.Columns() == 0 {
		return ErrNotFound
	}

	var columns = make([]Variable, v.Columns())
	var rviewer = reflect.ValueOf(v).Elem()

	for i := 0; i < v.Columns(); i++ {
		column := v.Column(i)

		columns[i] = rviewer.Field(column.Field()).Addr().Interface().(Variable)
	}

	_, err := v.Database().Search(f).Get(columns[0], columns[1:]...)

	return err
}

//Update updates the selected items with the given updates.
//Returns the number of items updated (or -1 if the statistic is unavailable).
func (f Filter) Update(update Update, updates ...Update) (int, error) {
	return update.Database().Search(f).Update(update, updates...)
}

//Delete deletes all the results from the database.
func (f Filter) Delete() (int, error) {
	return f.driver.Search(f).Delete()
}

//Read reads into the given variables.
func (f Filter) Read(v Variable, vs ...Variable) error {
	_, err := v.Database().Search(f).Get(v, vs...)
	return err
}

//SortBy sorts on the results of the filter with the given sorters.
func (f Filter) SortBy(sorter Sorter, sorters ...Sorter) Filter {
	f.Sort = sorter
	f.Sorts = sorters
	return f
}

//MarshalJSON encodes the results of the filter into JSON.
func (f Filter) MarshalJSON() ([]byte, error) {
	return f.driver.Search(f).MarshalJSON()
}

//Count counts the number of results.
func (f Filter) Count(v Viewable) (int, error) {
	return f.driver.Search(f).Count(v)
}

//Average sets the variable to the average value of all results of that column.
func (f Filter) Average(v Viewable) (float64, error) {
	return f.driver.Search(f).Average(v)
}

//Sum sets the variable to the sum of all results of that column.
func (f Filter) Sum(v Variable) error {
	return f.driver.Search(f).Sum(v)
}
