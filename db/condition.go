package db

//Condition can be used to specify a filter.
type Condition struct {
	Table, Column string
	Operator
	Value interface{}

	View Table

	driver Driver

	//If any of the cases match, this condition evaluates to true.
	Cases []Condition

	Invert bool
}

var True = Condition{}
var False = Condition{Operator: OpFalse}

//Switch returns a condition that is true if any of it's case conditions are true.
func Switch(first Condition, cases ...Condition) Condition {
	first.Cases = cases
	return first
}

//Either returns a condition that is true if either of it's arguments are true.
func Either(a, b Condition) Condition {
	return Switch(a, b)
}

//Both returns a condition that is true if all of it's arguments are true.
func Both(a, b Condition) Condition {
	c := Switch(a, b)
	c.Invert = true
	return c
}
