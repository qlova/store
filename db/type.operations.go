package db

//Operator type in a condition.
type Operator int

//Operators
const (
	OpTrue Operator = iota
	OpFalse
	OpEquals
	OpContains
	OpNotEquals
	OpHasPrefix
	OpLessThan
	OpDivisibleBy
)

//LessThan returns a condition that is true if i is less then val.
func (i Int64) LessThan(val int64) Condition {
	return Condition{
		Table:  i.table,
		View:   i.view,
		driver: i.driver,

		Column:   i.column,
		Operator: OpLessThan,
		Value:    val,
	}
}

//DivisibleBy returns a condition that is true if i is divisible by val.
func (i Int64) DivisibleBy(val int64) Condition {
	return Condition{
		Table:  i.table,
		View:   i.view,
		driver: i.driver,

		Column:   i.column,
		Operator: OpDivisibleBy,
		Value:    val,
	}
}

//Contains returns a condition that is true if s contains val.
func (s String) Contains(val string) Condition {
	return Condition{
		Table:  s.table,
		View:   s.view,
		driver: s.driver,

		Column:   s.column,
		Operator: OpContains,
		Value:    val,
	}
}

//HasPrefix returns a condition that is true if s starts with val.
func (s String) HasPrefix(val string) Condition {
	return Condition{
		Table:  s.table,
		View:   s.view,
		driver: s.driver,

		Column:   s.column,
		Operator: OpHasPrefix,
		Value:    val,
	}
}
